package model

import (
	"fmt"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/looplab/fsm"
	log "github.com/sirupsen/logrus"
)

const MaxSize = 21

type Shelf struct {
	Number string  `pg:",pk"`
	Crates []Crate `pg:"rel:has-many"`
	Status string  `pg:"default:'registered',notnull"`
}

var shelfFsm = fsm.Events{
	{Name: addInBatch, Src: []string{registered}, Dst: registered},
	{Name: removeFromBatch, Src: []string{registered}, Dst: registered},
	{Name: moveToCuttingStation, Src: []string{registered, inCuttingStation}, Dst: inCuttingStation},
	{Name: moveToPreScanning, Src: []string{inCuttingStation, inPreScanning}, Dst: inPreScanning},
	{Name: moveToScanStation, Src: []string{inPreScanning, inScanningStation}, Dst: inScanningStation},
	{Name: moveToArchiveStation, Src: []string{inScanningStation, inArchiveStation}, Dst: inArchiveStation},
	{Name: archive, Src: []string{inArchiveStation}, Dst: registered},
}

func (s *Shelf) isFull(crate Crate) bool {
	full := len(s.Crates) >= MaxSize
	alreadyIn := s.containCrate(crate)
	fullCrate := crate.isFull()
	switch {
	case full && alreadyIn:
		return fullCrate
	case full && !alreadyIn:
		return true
	case !full:
		return false
	default:
		return true
	}
}

func (s *Shelf) transition(event string, src string, currentCrate *Crate, db *Db, id gin_oidc.Identity) (err error) {
	newFSM := fsm.NewFSM(
		s.Status,
		shelfFsm,
		fsm.Callbacks{
			"before_" + addInBatch:           s.checkConstraintForAddingCrate(currentCrate),
			"before_" + moveToPreScanning:    s.checkBookletsStatus(db, src),
			"before_" + moveToScanStation:    s.checkBookletsStatus(db, src),
			"before_" + moveToArchiveStation: s.archiveAndPropagate(db, currentCrate, src),
			"after_event":                    s.updateState(db, id),
		},
	)
	err = newFSM.Event(event)
	if err != nil {
		return
	}
	s.Status = newFSM.Current()
	return
}

func (s *Shelf) removeCrate(crateNumber string) {
	for i, crate := range s.Crates {
		if crate.Number == crateNumber {
			s.Crates[i] = s.Crates[len(s.Crates)-1]
			s.Crates[len(s.Crates)-1] = Crate{}
			s.Crates = s.Crates[:len(s.Crates)-1]
		}
	}
}

func (s *Shelf) isEmpty() bool {
	return len(s.Crates) == 0
}

func (s *Shelf) recycle() {
	s.Status = registered
	s.Crates = nil
}

func (s *Shelf) archiveAndPropagate(db *Db, crate *Crate, src string) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		check := s.checkAllBooklets(db, src)
		check(event)
		if crate.isEmpty() {
			log.Debugf("in shelf : remove crate %s from shelf %s", crate.Number, s.Number)
			s.removeCrate(crate.Number)
		}
	}
}

func (s *Shelf) checkBookletsStatus(db *Db, src string) func(event *fsm.Event) {
	return s.checkAllBooklets(db, src)
}

func (s *Shelf) checkAllBooklets(db *Db, src string) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		status, err := db.GetBookletsStatus(s.Number)
		if err != nil {
			event.Cancel(fmt.Errorf("error when checking booklets status : %w", err))
			return
		}
		if !in(status, []string{src, event.FSM.Current()}) {
			event.Cancel(fmt.Errorf("all booklets in shelf must be at least in %s", src))
			return
		}
	}
}

func in(proposal, possible []string) bool {
	for _, s := range proposal {
		var found bool
		for _, s2 := range possible {
			if s2 == s {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (s *Shelf) updateState(db *Db, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if event.Err == nil {
			switch event.Dst {
			case inArchiveStation:
				log.Debugf("shelf state updated to %s", event.FSM.Current())
				s.Status = event.FSM.Current()
				if s.isEmpty() {
					log.Debugf("move empty shelf %s to status %s", s.Number, registered)
					s.Status = registered
				}
				err := db.UpdateEntities(nil, nil, s, nil)
				if err != nil {
					log.Errorf("error when updating shelf : %s", err)
				}
			default:
				log.Debugf("shelf state updated to %s", event.FSM.Current())
				s.Status = event.FSM.Current()
			}
		}
		err := db.RecordAction(id, event.Src, event.Event, s.Status, s)
		if err != nil {
			log.Errorf("unable to record event on history : %s", err)
		}
	}
}

func (s *Shelf) checkConstraintForAddingCrate(currentCrate *Crate) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if currentCrate.ShelfNumber != "" && currentCrate.ShelfNumber != s.Number {
			event.Cancel(fmt.Errorf("this crate has already been registered in another shelf"))
			return
		}
		if s.isFull(*currentCrate) {
			event.Cancel(fmt.Errorf("this shelf is full, please use an empty shelf"))
			return
		}
		if event.FSM.Current() != registered {
			event.Cancel(fmt.Errorf("this shelf is already moved to another production step, it is not possible to add crate"))
			return
		}
	}
}

func (s *Shelf) GetType() string {
	return "Shelf"
}

func (s *Shelf) GetId() string {
	return s.Number
}

func (s *Shelf) GetFsmGraph() string {
	f := fsm.NewFSM(registered, shelfFsm, nil)
	return fsm.Visualize(f)
}

func (s *Shelf) containCrate(crate Crate) bool {
	for _, c := range s.Crates {
		if c.Number == crate.Number {
			return true
		}
	}
	return false
}
