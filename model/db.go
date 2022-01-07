package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	ginOidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/conf"
)

type Db struct {
	Conn *pg.DB
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		log.Error(err)
	}
	log.Info(string(query))
	return nil
}

func (db *Db) Init(conf conf.Config) (err error) {
	db.Conn = pg.Connect(&pg.Options{
		Addr:     conf.DbHost,
		Dialer:   nil,
		User:     conf.DbUser,
		Password: conf.DbPassword,
		Database: conf.DbDatabase,
	})
	if conf.DbLog {
		logger := dbLogger{}
		db.Conn.AddQueryHook(logger)
	}
	if conf.DbInit {
		db.createExtension()
		err = db.createSchema()
		if err != nil {
			return
		}

	}
	return
}

func (db *Db) Close() {
	err := db.Conn.Close()
	if err != nil {
		log.Error(err)
	}
}

func (db *Db) createSchema() (err error) {
	for _, model := range []interface{}{
		(*Agregated)(nil),
	} {
		err := db.Conn.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return
}

func (db *Db) createExtension() {
	_, err := db.Conn.Exec("CREATE EXTENSION if not exists ltree;")
	if err != nil {
		log.Error(err)
	}
}

func (db *Db) GetAgregate(division, district, upazilla, union, mouza uint, tableName string) {
	agregate := Agregated{
		Geocode: fmt.Sprintf("%s.%s.%s.%d", division, district, upazilla, union, mouza),
	}
	err = db.Conn.Model(&agregate).Select()

}

func (db *Db) GetBooklet(bookletNumber string) (booklet Booklet, err error) {
	booklet = Booklet{
		Number: bookletNumber,
	}
	err = db.Conn.Model(&booklet).
		Relation("ArchiveBox").
		Relation("Crate").
		Relation("Crate.Booklets").
		Relation("Crate.Shelf").
		Relation("Crate.Shelf.Crates").
		// Relation("Crate.Shelf.Crates.Booklets").
		Where("Booklet.number = ?", bookletNumber).
		Select()
	return
}

func (db *Db) GetGeoCode(geoCodeNumber string) (geoCode GeoCodes, err error) {
	geoCode = GeoCodes{
		GeocodeID: geoCodeNumber,
	}
	err = db.Conn.Model(&geoCode).WherePK().Select()
	return
}

func (db *Db) GetArchiveBox(archiveBoxNumber string) (archiveBox ArchiveBox, err error) {
	archiveBox = ArchiveBox{
		Number: archiveBoxNumber,
	}
	err = db.Conn.Model(&archiveBox).WherePK().Select()
	return
}

func (db *Db) GetBookletsStatus(shelfNumber string) (status []string, err error) {
	_, err = db.Conn.Query(&status, `
	select distinct(b.status) as status
from booklets b,
     crates c,
     shelves s
where b.crate_number = c.number
  and c.shelf_number = s.number
  and s.number = ?;`, shelfNumber)
	return
}

func (db *Db) GetBookletOnly(bookletNumber string) (booklet Booklet, err error) {
	booklet = Booklet{Number: bookletNumber}
	err = db.Conn.Model(&booklet).WherePK().Select()
	return
}

func (db *Db) GetCrate(crateNumber string) (crate Crate, err error) {
	crate = Crate{
		Number: crateNumber,
	}
	err = db.Conn.Model(&crate).
		Relation("Booklets").
		Relation("Shelf").
		Where("Crate.number = ?", crateNumber).
		Select()
	return
}

func (db *Db) GetShelf(shelfNumber string) (shelf Shelf, err error) {
	shelf = Shelf{
		Number: shelfNumber,
	}
	err = db.Conn.Model(&shelf).
		Relation("Crates").
		Where("Shelf.number = ?", shelfNumber).
		Select()
	return
}

func (db *Db) GetOrCreateArchiveBox(archiveBoxNumber string) (archiveBox ArchiveBox, err error) {
	archiveBox = ArchiveBox{
		Number: archiveBoxNumber,
	}
	err = db.Conn.Model(&archiveBox).
		Where("number = ?", archiveBoxNumber).
		Relation("Booklets").
		Select()
	if errors.Is(err, pg.ErrNoRows) {
		_, err = db.Conn.Model(&archiveBox).Insert()
		if err != nil {
			return ArchiveBox{}, err
		}
	}
	return
}

func (db *Db) RecordAction(id ginOidc.Identity, sourceState, eventType, finalState string, subject EventSubject) (err error) {

	e := Event{
		SourceState: sourceState,
		EventType:   eventType,
		FinalState:  finalState,
		TimeStamp:   time.Now(),
		IdentityId:  id.Id,
		Identity:    &id,
		Subject:     subject,
	}

	_, err = db.Conn.Model(e.Identity).Where("id = ?id").SelectOrInsert()

	if err == nil {
		switch e.Subject.(type) {
		case *Booklet:
			b := e.Subject.(*Booklet)
			e.BookletNumber = b.Number
			e.CrateNumber = b.CrateNumber
			e.ShelfNumber = b.GetShelfNumber()
			e.ArchiveBoxNumber = b.ArchiveBoxNumber
		case *Crate:
			c := e.Subject.(*Crate)
			e.CrateNumber = c.Number
			e.ShelfNumber = c.ShelfNumber
		case *Shelf:
			e.ShelfNumber = e.Subject.GetId()
		case *ArchiveBox:
			e.ArchiveBoxNumber = e.Subject.GetId()
		}
		_, err = db.Conn.Model(&e).Insert()
	}
	return
}

func (db *Db) CheckWarehousePosition(row, shelf, level int) (found bool, err error) {
	var count int
	_, err = db.Conn.Model((*Warehouse)(nil)).QueryOne(pg.Scan(&count), `
	select count(*) from ?TableName where row = ? and shelf = ? and ? = any(shelf_level);`, row, shelf, level)
	if err != nil {
		return false, err
	} else {
		return count == 1, err
	}
}
