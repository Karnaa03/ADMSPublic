package model

import (
	"fmt"
	"time"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/looplab/fsm"
	log "github.com/sirupsen/logrus"
)

// Capacity
const maxBooklets = 10

type Crate struct {
	Number      string    `pg:",pk,notnull"`
	Booklets    []Booklet `pg:"rel:has-many"`
	ShelfNumber string
	Shelf       *Shelf    `pg:"rel:has-one"`
	Status      string    `pg:"default:'registered',notnull"`
	ScannedOn   time.Time `binding:"-"`
	ScannedBy   string
}

var crateFsm = fsm.Events{
	{Name: addInBatch, Src: []string{registered}, Dst: registered},
	{Name: removeFromBatch, Src: []string{registered}, Dst: registered},
	{Name: moveToCuttingStation, Src: []string{registered, inCuttingStation}, Dst: inCuttingStation},
	{Name: moveToPreScanning, Src: []string{inCuttingStation, inPreScanning}, Dst: inPreScanning},
	{Name: moveToScanStation, Src: []string{inPreScanning, inScanningStation}, Dst: inScanningStation},
	{Name: moveToArchiveStation, Src: []string{inScanningStation, inArchiveStation}, Dst: inArchiveStation},
	{Name: archive, Src: []string{inArchiveStation}, Dst: registered},
}

func (c *Crate) Scann(shelfNumber string, id gin_oidc.Identity, db *Db) (err error) {
	for _, b := range c.Booklets {
		log.Debugf("propagate event down %s for booklet %s", moveToScanStation, b.Number)
		booklet, err := db.GetBooklet(b.Number)
		if err != nil {
			return err
		}
		err = booklet.transition(moveToScanStation, db, c.Number, shelfNumber, "", "", id)
		if err != nil {
			switch err.(type) {
			case fsm.NoTransitionError:
			default:
				return err
			}
		}
	}
	return
}

func (c *Crate) transition(event string, src string, booklet *Booklet, db *Db, id gin_oidc.Identity) (err error) {
	newFSM := fsm.NewFSM(
		c.Status,
		crateFsm,
		fsm.Callbacks{
			"before_" + addInBatch:           c.addingInBatch(db, booklet, id),
			"before_" + removeFromBatch:      c.removeFromBatch(booklet, id, db),
			"before_" + moveToCuttingStation: c.propagateEvent(db, src, id),
			"before_" + moveToPreScanning:    c.propagateEvent(db, src, id),
			"before_" + moveToScanStation:    c.propagateEvent(db, src, id),
			"before_" + moveToArchiveStation: c.archiveAndPropagate(booklet, db, src, id),
			"after_event":                    c.updateState(id, db),
		},
	)
	err = newFSM.Event(event)
	if err != nil {
		return
	}
	c.Status = newFSM.Current()
	return
}

func (c *Crate) addingInBatch(db *Db, booklet *Booklet, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if c.isFull() {
			event.Cancel(fmt.Errorf("this crate is full, please use an empty crate"))
			return
		}
		if event.FSM.Current() != registered {
			event.Cancel(fmt.Errorf("this crate is already moved to another production step, it is not possible to add booklet"))
			return
		}
		err := c.Shelf.transition(addInBatch, "", c, db, id)
		if err != nil {
			switch err.(type) {
			case fsm.NoTransitionError:
			default:
				event.Cancel(fmt.Errorf("shelf : %w", err))
				return
			}
		}
		c.Booklets = append(c.Booklets, *booklet)
	}
}

func (c *Crate) archiveAndPropagate(booklet *Booklet, db *Db, src string, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		log.Debugf("remove archived booklet %s from crate %s", booklet.Number, c.Number)
		c.removeBooklet(booklet.Number)
		shelf := c.Shelf
		p := c.propagate(db, src, shelf, id)
		p(event)
	}
}

func (c *Crate) propagateEvent(db *Db, src string, id gin_oidc.Identity) func(event *fsm.Event) {
	return c.propagate(db, src, c.Shelf, id)
}

func (c *Crate) propagate(db *Db, src string, shelf *Shelf, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if event.Err == nil {
			log.Debugf("propagate event %s to shelf", event.Event)
			if shelf != nil {
				if src == "" {
					src = event.Src
				}
				err := shelf.transition(event.Event, src, c, db, id)
				if err != nil {
					switch err.(type) {
					case fsm.NoTransitionError:
					default:
						log.Errorf("error on shelf transition %s : %s", event.Event, err)
						event.Cancel(err)
						return
					}
				}
			}
		}
	}
}

func (c *Crate) updateState(id gin_oidc.Identity, db *Db) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if event.Err == nil {
			switch event.Dst {
			case inScanningStation:
			case inArchiveStation:
				log.Debugf("crate state updated to %s", event.FSM.Current())
				c.Status = event.FSM.Current()
				if c.isEmpty() {
					log.Debugf("remove empty crate %s from shelf %s and move them to status %s", c.Number, c.ShelfNumber, registered)
					c.Status = registered
					c.ShelfNumber = ""
					c.Shelf = nil
				}
				err := db.UpdateEntities(nil, c, nil, nil)
				if err != nil {
					log.Errorf("error when updating shelf : %s", err)
				}
			default:
				log.Debugf("crate state updated to %s", event.FSM.Current())
				c.Status = event.FSM.Current()
			}
		}
		err := db.RecordAction(id, event.Src, event.Event, c.Status, c)
		if err != nil {
			log.Errorf("unable to record event on history : %s", err)
		}
	}
}

func (c *Crate) removeFromBatch(booklet *Booklet, id gin_oidc.Identity, db *Db) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		err := c.Shelf.transition(removeFromBatch, "", c, db, id)
		if err != nil {
			switch err.(type) {
			case fsm.NoTransitionError:
			default:
				event.Cancel(err)
				return
			}
		}
		c.removeBooklet(booklet.Number)
		if c.isEmpty() {
			c.ShelfNumber = "" // the crate is empty then we can remove from shelf
		}
	}
}

func (c *Crate) isFull() bool {
	return len(c.Booklets) >= maxBooklets
}

func (c *Crate) isEmpty() bool {
	return len(c.Booklets) == 0
}

func (c *Crate) recycle() {
	c.ShelfNumber = ""
	c.Shelf = nil
	c.Status = registered
	c.ScannedBy = ""
	c.ScannedOn = time.Time{}
}

func (c *Crate) GetType() string {
	return "Crate"
}

func (c *Crate) GetId() string {
	return c.Number
}

func (c *Crate) getFsmGraph() string {
	f := fsm.NewFSM(registered, crateFsm, nil)
	return fsm.Visualize(f)
}

func (c *Crate) removeBooklet(bookletNumber string) {
	for i, booklet := range c.Booklets {
		if booklet.Number == bookletNumber {
			c.Booklets[i] = c.Booklets[len(c.Booklets)-1]
			c.Booklets[len(c.Booklets)-1] = Booklet{}
			c.Booklets = c.Booklets[:len(c.Booklets)-1]
		}
	}
}
