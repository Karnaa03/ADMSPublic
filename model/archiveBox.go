package model

import (
	"fmt"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/looplab/fsm"
	log "github.com/sirupsen/logrus"
)

const MaxBookletMerArchiveBoc = 10

type ArchiveBox struct {
	Number                    string    `pg:",pk" form:"ArchiveBoxNumber"`
	Booklets                  []Booklet `pg:"rel:has-many"`
	Status                    string    `pg:"default:'registered',notnull"`
	WarehouseRowNumber        int       `form:"WarehouseRowNumber" binding:"gt=0,lt=99"`
	WarehouseShelfNumber      int       `form:"WarehouseShelfNumber" binding:"gt=0,lt=999"`
	WarehouseShelfLevelNumber int       `form:"WarehouseShelfLevelNumber" binding:"gt=0,lt=9"`
	CheckedOutBy              string
	CheckedInBy               string
}

var archiveBoxFsm = fsm.Events{
	{Name: moveToArchiveStation, Src: []string{registered, inUse}, Dst: inUse},
	{Name: archive, Src: []string{inUse}, Dst: archived},
	{Name: checkOut, Src: []string{archived}, Dst: checkedOut},
	{Name: checkIn, Src: []string{checkedOut}, Dst: archived},
	{Name: moveBookletFromBox, Src: []string{archived}, Dst: archived},
}

func (a *ArchiveBox) transition(event string, db *Db, id gin_oidc.Identity) (err error) {
	newFSM := fsm.NewFSM(
		a.Status,
		archiveBoxFsm,
		fsm.Callbacks{
			"enter_state": a.updateState(db, id),
		},
	)
	err = newFSM.Event(event)
	if err != nil {
		return
	}
	a.Status = newFSM.Current()
	return
}

func (a *ArchiveBox) updateState(db *Db, id gin_oidc.Identity) func(event *fsm.Event) {
	return func(event *fsm.Event) {
		if event.Err == nil {
			switch {
			case event.Dst == archived:
				switch event.Src {
				case inUse:
					a.Status = event.FSM.Current()
					err := db.UpdateEntities(nil, nil, nil, a)
					if err != nil {
						log.Errorf("unable to update archive box in database : %s", err)
					}
				case checkedOut:
					a.Status = event.FSM.Current()
					a.CheckedInBy = id.FullName
					err := db.UpdateEntities(nil, nil, nil, a)
					if err != nil {
						log.Errorf("unable to update archive box in database : %s", err)
					}
				}
			case event.Dst == checkedOut:
				a.Status = event.FSM.Current()
				a.CheckedOutBy = id.FullName
				err := db.UpdateEntities(nil, nil, nil, a)
				if err != nil {
					log.Errorf("unable to update archive box in database : %s", err)
				}
			default:
				log.Debugf("archive box state updated to %s", event.FSM.Current())
				a.Status = event.FSM.Current()
			}
		}
		err := db.RecordAction(id, event.Src, event.Event, event.FSM.Current(), a)
		if err != nil {
			log.Errorf("unable to record event on history : %s", err)
		}
	}
}

func (a *ArchiveBox) GetType() string {
	return "ArchiveBox"
}

func (a *ArchiveBox) GetId() string {
	return a.Number
}

func (a *ArchiveBox) GetFsmGraph() string {
	f := fsm.NewFSM(registered, archiveBoxFsm, nil)
	return fsm.Visualize(f)
}

func (a *ArchiveBox) RegisterInWarehouse(id gin_oidc.Identity, db *Db, registrationRequest ArchiveBox) (err error) {
	found, err := db.CheckWarehousePosition(registrationRequest.WarehouseRowNumber, registrationRequest.WarehouseShelfNumber, registrationRequest.WarehouseShelfLevelNumber)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("no warehouse position found with row : %d, shelf : %d and shelf level : %d",
			registrationRequest.WarehouseRowNumber,
			registrationRequest.WarehouseShelfNumber,
			registrationRequest.WarehouseShelfLevelNumber)
	}
	a.WarehouseRowNumber = registrationRequest.WarehouseRowNumber
	a.WarehouseShelfNumber = registrationRequest.WarehouseShelfNumber
	a.WarehouseShelfLevelNumber = registrationRequest.WarehouseShelfLevelNumber
	err = a.transition(archive, db, id)
	return
}

func (a *ArchiveBox) CheckOut(id gin_oidc.Identity, db *Db) (err error) {
	err = a.transition(checkOut, db, id)
	return
}

func (a *ArchiveBox) CheckIn(id gin_oidc.Identity, db *Db) (err error) {
	err = a.transition(checkIn, db, id)
	return
}

func (a ArchiveBox) GetActionLink() string {
	switch a.Status {
	case "archived":
		return fmt.Sprintf(`			
				<a href="#" data-href="/warehouse/searchInArchives.html?checkout=%s" data-toggle="modal" data-target="#confirm-checkout">
					<center>
						<i class="fa fa-sign-out"></i>
					</center>
				</a>			
		`, a.Number)
	case "checkedOut":
		return fmt.Sprintf(`
				<a href="#" data-href="/warehouse/searchInArchives.html?checkin=%s" data-toggle="modal" data-target="#confirm-checkin">
					<center>
						<i class="fa fa-sign-in"></i>
					</center>
				</a>
		`, a.Number)
	default:
		return ""
	}
}
