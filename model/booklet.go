package model

import (
	"fmt"
	"strconv"
	"time"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/looplab/fsm"
	"github.com/prometheus/common/log"
)

type Booklet struct {
	/*
	  Number : 6 digit that must :
	  	* Checking this number within the list given by BBS which contains the entire list of booklets recorded by Devnet
	  	* Checking this number within the tracking solution. Duplicate booklet number is not allowed.
	  An error message needs to be highlighted on the screen and an investigation needs to be launched to determine the mistake and to correct it.
	  To not waste time, the investigation will be done by another operator and the typist will continue to enter the data from another booklet.
	*/
	Number string `form:"BookletNumber" pg:",pk,notnull,unique"`
	/*
		GeoCode : The GEO Code number; group of:
		* District: 2 digits
		* Upazilla: 2 digits
		* Union: 3 digits
		* Mouza: 3 digits
		* Counting Area: 3 digits
		* RMO Code: 1 digit
	*/
	GeoCodeID string    `form:"GeoCode" pg:",notnull,type:ltree"`
	GeoCode   *GeoCodes `binding:"-"  pg:"rel:has-one,fk:geo_code_id"`
	/*
		Size : size of booklet, allowed values :
			* 25
			* 50
			* 100
	*/
	Size             uint   `form:"BookletSize" pg:",notnull"`
	Status           string `pg:"default:'notRegistered',notnull"`
	CrateNumber      string
	Crate            *Crate `pg:"rel:has-one"`
	ArchiveBoxNumber string
	ArchiveBox       *ArchiveBox `pg:"rel:has-one" `
	RegisteredOn     time.Time   `pg:"default:now(),notnull" binding:"-"`
	RegisteredBy     string      `pb:",notnull" binding:"-"`
	AddedInBatchOn   time.Time   `binding:"-"`
	AddedInBatchBy   string
	CutOn            time.Time `binding:"-"`
	CutBy            string
	PreparedOn       time.Time `binding:"-"`
	PreparedBy       string
	ScannedOn        time.Time `binding:"-"`
	ScannedBy        string
	ArchivedOn       time.Time `binding:"-"`
	ArchivedBy       string
}

var bookletFSM = fsm.Events{
	{Name: register, Src: []string{notRegistered}, Dst: registered},
	{Name: deregister, Src: []string{registered}, Dst: notRegistered},
	{Name: addInBatch, Src: []string{registered}, Dst: inBatch},
	{Name: removeFromBatch, Src: []string{inBatch}, Dst: registered},
	{Name: moveToCuttingStation, Src: []string{inBatch}, Dst: inCuttingStation},
	{Name: moveToPreScanning, Src: []string{inCuttingStation}, Dst: inPreScanning},
	{Name: freeze, Src: []string{inCuttingStation}, Dst: iceBox},
	{Name: moveToScanStation, Src: []string{inPreScanning}, Dst: inScanningStation},
	{Name: moveToArchiveStation, Src: []string{inScanningStation, iceBox}, Dst: archived},
	{Name: moveBookletFromBox, Src: []string{archived}, Dst: archived},
}

func RegisterNewBooklet(number, geoCode, size string, id gin_oidc.Identity, db *Db) (err error) {
	bSize, err := strconv.ParseUint(size, 10, 8)
	if err != nil {
		return
	}

	booklet, err := db.GetBookletOnly(number)
	if err != nil {
		booklet = Booklet{
			Number:       number,
			GeoCodeID:    geoCode,
			Size:         uint(bSize),
			Status:       notRegistered,
			RegisteredBy: id.FullName,
			RegisteredOn: time.Now(),
		}
		return booklet.transition(register, db, "", "", "", "", id)
	} else {
		return booklet.transition(register, db, "", "", "", "", id)
	}
}

func (b *Booklet) RegisterInBatch(crateNumber, shelfNumber string, id gin_oidc.Identity, db *Db) (err error) {
	return b.transition(addInBatch, db, crateNumber, shelfNumber, "", "", id)
}

func (b *Booklet) RemoveFromBatch(db *Db, id gin_oidc.Identity) (err error) {
	if b.notInCrate() {
		return fmt.Errorf("not in a crate")
	}
	err = b.transition(removeFromBatch, db, "", "", "", "", id)
	return
}

func (b *Booklet) CutBooklet(crateNumber string, shelfNumber string, id gin_oidc.Identity, db *Db) (err error) {
	err = b.transition(moveToCuttingStation, db, crateNumber, shelfNumber, "", "", id)
	return
}

func (b *Booklet) FreezeBooklet(crateNumber string, shelfNumber string, id gin_oidc.Identity, db *Db) (err error) {
	err = b.transition(freeze, db, crateNumber, shelfNumber, "", "", id)
	return
}

func (b *Booklet) PreScann(id gin_oidc.Identity, crateNumber, shelfNumber string, db *Db) (err error) {
	err = b.transition(moveToPreScanning, db, crateNumber, shelfNumber, "", "", id)
	return
}

func (b *Booklet) Archive(id gin_oidc.Identity, crateNumber, shelfNumbers, archiveBoxNumber string, db *Db) (err error) {
	archiveBox, err := db.GetOrCreateArchiveBox(archiveBoxNumber)
	if err != nil {
		return err
	}
	if len(archiveBox.Booklets) >= MaxBookletMerArchiveBoc {
		return fmt.Errorf("maximum allowed booklets per archive box : 10, current : %d", len(archiveBox.Booklets))
	}
	b.ArchiveBox = &archiveBox
	b.ArchiveBoxNumber = archiveBox.Number
	err = b.transition(moveToArchiveStation, db, crateNumber, shelfNumbers, "", "", id)
	return
}

func (b *Booklet) MoveInBox(id gin_oidc.Identity, sourceBox, destBox string, db *Db) (err error) {
	err = b.transition(moveBookletFromBox, db, "", "", sourceBox, destBox, id)
	return
}

func (b *Booklet) transition(event string, db *Db, crateNumber, shelfNumber, sourceBox, destBox string, id gin_oidc.Identity) (err error) {
	newFSM := fsm.NewFSM(
		b.Status,
		bookletFSM,
		fsm.Callbacks{
			"enter_" + registered:            b.register(db),
			"before_" + addInBatch:           b.checkForAddingInBatch(crateNumber, shelfNumber, db, id),
			"before_" + removeFromBatch:      b.propagateEvent("", "", db, id),
			"after_" + removeFromBatch:       b.removeBookletFromBatch(db),
			"before_" + moveToCuttingStation: b.propagateEvent(crateNumber, shelfNumber, db, id),
			"before_" + moveToPreScanning:    b.propagateEvent(crateNumber, shelfNumber, db, id),
			"before_" + moveToScanStation:    b.propagateEvent(crateNumber, shelfNumber, db, id),
			"before_" + moveToArchiveStation: b.archiveAndPropagate(crateNumber, shelfNumber, db, id),
			"before_" + freeze:               b.freeze(crateNumber, shelfNumber, db, id),
			"before_" + moveBookletFromBox:   b.moveInBox(sourceBox, destBox, db, id),
			"after_event":                    b.updateState(id, db), // not applied for register & removeFromBatch
		},
	)
	err = newFSM.Event(event)
	if err != nil {
		switch err.(type) {
		case fsm.NoTransitionError:
			err = nil
		default:
			return
		}
	}
	b.Status = newFSM.Current()
	return

}

func (b *Booklet) freeze(crateNumber, shelfNumber string, db *Db, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		check := b.checkCrateAndShelfLinks(crateNumber, shelfNumber, db)
		check(event)

		crate := b.Crate
		shelf := b.Crate.Shelf

		// Release Crate if empty after Booklet removal
		crate.removeBooklet(b.Number)
		if crate.isEmpty() {
			crate.recycle()
			// Release Shelf if empty after current crate removal
			shelf.removeCrate(crate.Number)
			if shelf.isEmpty() {
				shelf.recycle()
			}
		}
		err := db.UpdateEntities(nil, crate, shelf, nil)
		if err != nil {
			log.Errorf("error when trying to update crate and shelf on frozen booklet : %s", err)
			event.Cancel(err)
		}
		b.Crate = nil
		b.CrateNumber = ""
	}
}

func (b *Booklet) archiveAndPropagate(crateNumber, shelfNumber string, db *Db, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		// Frozen booklet hasn't crate & shelf numbers
		if crateNumber != "" && shelfNumber != "" {
			check := b.checkCrateAndShelfLinks(crateNumber, shelfNumber, db)
			check(event)

			crate := b.Crate
			b.Crate = nil
			b.CrateNumber = ""
			propagate := b.propagate(db, crate, id)
			propagate(event)
		} else {
			propagate := b.propagate(db, nil, id)
			propagate(event)
		}
	}
}

func (b *Booklet) moveInBox(srcBox, dstBox string, db *Db, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if b.ArchiveBoxNumber != srcBox {
			msg := fmt.Errorf("the given source box (%s) didn't correspond to actually recorded : %s", srcBox, b.ArchiveBoxNumber)
			log.Error(msg)
			event.Cancel(msg)
		}
		if b.ArchiveBox == nil || b.ArchiveBox.Status != archived {
			msg := fmt.Errorf("the given source box (%s) is not in archived", srcBox)
			log.Error(msg)
			event.Cancel(msg)
		}
		dest, err := db.GetArchiveBox(dstBox)
		if err != nil {
			msg := fmt.Errorf("unable to find destination archive box : %s", err)
			log.Error(msg)
			event.Cancel(msg)
		}
		b.ArchiveBoxNumber = dstBox
		b.ArchiveBox = &dest
		propagate := b.propagate(db, nil, id)
		propagate(event)
	}
}

func (b *Booklet) register(db *Db) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if event.Src != inBatch {
			log.Info("record booklet")
			b.Status = event.FSM.Current()
			err := db.PutBooklet(*b)
			if err != nil {
				event.Cancel(fmt.Errorf("Unable to save booklet in database for registration : %w", err))
				return
			}
		}
	}
}

func (b *Booklet) checkCrateAndShelfLinks(crateNumber, shelfNumber string, db *Db) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if crateNumber == "" || shelfNumber == "" {
			event.Cancel(fmt.Errorf("Crate and Shelf number is mandatory"))
			return
		}
		if b.CrateNumber != "" && b.CrateNumber != crateNumber {
			event.Cancel(fmt.Errorf("This booklet has been registered in another crate (%s)", b.CrateNumber))
			return
		}
		crate, err := db.GetCrate(crateNumber)
		if err != nil {
			event.Cancel(fmt.Errorf("Can't find crate with number %s", crateNumber))
			return
		}
		b.Crate = &crate
		b.CrateNumber = crate.Number

		shelf, err := db.GetShelf(shelfNumber)
		if err != nil {
			event.Cancel(fmt.Errorf("Can't find shelf with number %s", shelfNumber))
			return
		}

		if crate.ShelfNumber != "" && crate.ShelfNumber != shelf.Number {
			event.Cancel(fmt.Errorf("This crate has already been registered in another shelf (%s)", crate.ShelfNumber))
			return
		}
		b.Crate.Shelf = &shelf
		b.Crate.ShelfNumber = shelf.Number
	}
}

func (b *Booklet) checkForAddingInBatch(crateNumber, shelfNumber string, db *Db, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if !b.notInCrate() {
			event.Cancel(fmt.Errorf("Already in a crate"))
			return
		}
		check := b.checkCrateAndShelfLinks(crateNumber, shelfNumber, db)
		check(event)

		err := b.Crate.transition(addInBatch, event.Src, b, db, id)
		if err != nil {
			switch err.(type) {
			case fsm.NoTransitionError:
			default:
				event.Cancel(fmt.Errorf("Crate : %w", err))
				return
			}
		}
	}
}

func (b *Booklet) removeBookletFromBatch(db *Db) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		log.Infof("remove booklet %s from crate %s and shelf %s", b.Number, b.CrateNumber, b.Crate.ShelfNumber)
		crate := b.Crate
		shelf := crate.Shelf
		b.Status = event.FSM.Current()
		b.CrateNumber = ""
		b.AddedInBatchBy = ""
		b.AddedInBatchOn = time.Time{}
		err := db.UpdateEntities(b, crate, shelf, nil)
		if err != nil {
			event.Cancel(fmt.Errorf("Unable to update booklet in database for batch removal : %w", err))
			return
		}
	}
}

func (b *Booklet) propagateEvent(crateNumber string, shelfNumber string, db *Db, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		check := b.checkCrateAndShelfLinks(crateNumber, shelfNumber, db)
		check(event)
		propagate := b.propagate(db, b.Crate, id)
		propagate(event)
	}
}

func (b *Booklet) propagate(db *Db, crate *Crate, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if event.Err == nil {
			if crate != nil {
				log.Debugf("propagate event %s to crate", event.Event)
				err := crate.transition(event.Event, event.Src, b, db, id)
				if err != nil {
					switch err.(type) {
					case fsm.NoTransitionError:
					default:
						log.Errorf("error on crate transition %s : %s", event.Event, err)
						event.Cancel(err)
						return
					}
				}
			}
			if b.ArchiveBox != nil {
				log.Debugf("propagate event %s to archiveBox", event.Event)
				err := b.ArchiveBox.transition(event.Event, db, id)
				if err != nil {
					switch err.(type) {
					case fsm.NoTransitionError:
						log.Debugf("archiveBox remain in same state")
					default:
						event.Cancel(err)
						return
					}
				}
			}
		}
	}
}

func (b *Booklet) updateState(id gin_oidc.Identity, db *Db) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		log.Debugf("booklet state updated to %s", event.FSM.Current())
		if event.Err == nil && event.Event != register && event.Event != removeFromBatch {
			b.Status = event.FSM.Current()
			switch event.Dst {
			case inBatch:
				b.AddedInBatchOn = time.Now()
				b.AddedInBatchBy = id.FullName
			case inCuttingStation:
				b.CutOn = time.Now()
				b.CutBy = id.FullName
			case inPreScanning:
				b.PreparedOn = time.Now()
				b.PreparedBy = id.FullName
			case inScanningStation:
				b.ScannedOn = time.Now()
				b.ScannedBy = id.FullName
			case archived:
				b.ArchivedOn = time.Now()
				b.ArchivedBy = id.FullName
			}
			var err error
			if b.Crate == nil {
				err = db.UpdateEntities(b, nil, nil, b.ArchiveBox)
			} else {
				err = db.UpdateEntities(b, b.Crate, b.Crate.Shelf, b.ArchiveBox)
			}
			if err != nil {
				event.Cancel(fmt.Errorf("Unable to update booklet status in database : %w", err))
				return
			}
		}
		err := db.RecordAction(id, event.Src, event.Event, b.Status, b)
		if err != nil {
			log.Errorf("unable to record event on history : %s", err)
		}
	}
}

func (b *Booklet) notInCrate() bool {
	return b.Crate == nil
}

type BookletPersistance interface {
	GetBooklet(bookletNumber string) (booklet Booklet, err error)
	PutBooklet(booklet Booklet) (err error)
}

func (b *Booklet) GetType() string {
	return b.Number
}

func (b *Booklet) GetId() string {
	return b.Number
}

func (b *Booklet) GetCrateNumber() string {
	if b.Crate != nil {
		return b.Crate.Number
	}
	return ""
}

func (b *Booklet) GetShelfNumber() string {
	if b.Crate != nil {
		if b.Crate.Shelf != nil {
			return b.Crate.Shelf.Number
		}
	}
	return ""
}

func (b *Booklet) GetFsmGraph() string {
	f := fsm.NewFSM(notRegistered, bookletFSM, nil)
	return fsm.Visualize(f)
}

func (b Booklet) GetArchiveBox() ArchiveBox {
	if b.ArchiveBox == nil {
		return ArchiveBox{}
	} else {
		return *b.ArchiveBox
	}
}

func (b Booklet) GetActionLink() string {
	if b.ArchiveBox != nil {
		return b.ArchiveBox.GetActionLink()
	} else {
		return ""
	}
}

func (b Booklet) GetCheckedOutBy() string {
	if b.ArchiveBox != nil {
		return b.ArchiveBox.CheckedOutBy
	} else {
		return ""
	}
}

func (b Booklet) GetCheckedInBy() string {
	if b.ArchiveBox != nil {
		return b.ArchiveBox.CheckedInBy
	} else {
		return ""
	}
}
