package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		db.createIndex()
		db.registerViewAndProcedure()
		db.loadInitialData()
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
		(*GeoCodes)(nil),
		(*ArchiveBox)(nil),
		(*Shelf)(nil),
		(*Crate)(nil),
		(*Booklet)(nil),
		(*ginOidc.Identity)(nil),
		(*Event)(nil),
		(*Warehouse)(nil),
		(*Tree)(nil),
		(*Crop)(nil),
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

func (db *Db) registerViewAndProcedure() {
	for name, proc := range procedures {
		log.Infof("register %s", name)
		_, err := db.Conn.Exec(proc)
		if err != nil {
			log.Error(err)
		}
	}
}

func (db *Db) createIndex() {
	for name, index := range index {
		log.Infof("create index : %s", name)
		_, err := db.Conn.Model((*Booklet)(nil)).Exec(index)
		if err != nil {
			log.Error(err)
		}
	}
}

func (db *Db) loadInitialData() {
	db.loadShelves()
	db.loadCrates()
	db.loadTree()
	db.loadCrop()
}

func (db *Db) loadShelves() {
	var count int
	_, err := db.Conn.Model((*Shelf)(nil)).QueryOne(pg.Scan(&count), `
	SELECT count(*) FROM ?TableName`)
	if err != nil {
		log.Error(err)
		return
	}
	if count == 0 {
		log.Info("loading shelves, this can take a while...")
		for s := 0; s < 1_000; s++ {
			shelf := Shelf{
				Number: fmt.Sprintf("S.%03d", s),
			}
			_, err := db.Conn.Model(&shelf).Insert()
			if err != nil {
				log.Error(err)
				continue
			}
		}
	} else {
		log.Info("table Shelves not empty, no need to fill them")
	}
}

func (db *Db) loadCrates() {
	var count int
	_, err := db.Conn.Model((*Crate)(nil)).QueryOne(pg.Scan(&count), `
	SELECT count(*) FROM ?TableName`)
	if err != nil {
		log.Error(err)
		return
	}
	if count == 0 {
		log.Info("loading crates, this can take a while...")
		for c := 0; c < 10_000; c++ {
			crate := Crate{
				Number: fmt.Sprintf("C.%04d", c),
			}
			_, err := db.Conn.Model(&crate).Insert()
			if err != nil {
				log.Error(err)
				continue
			}
		}
	} else {
		log.Info("table Crates not empty, no need to fill them")
	}
}

func (db *Db) loadTree() {
	if trees, ok := referencesData["tree"]; ok {
		for _, tree := range trees {
			t := Tree{
				Id:   tree.Id,
				Name: strings.ToUpper(tree.Name),
			}
			_, err := db.Conn.Model(&t).Insert()
			if pgerr, ok := err.(pg.Error); ok {
				if !pgerr.IntegrityViolation() {
					log.Errorf("error when trying to insert tree reference data : %s", pgerr.Error())
				}
			}
		}
	}
}

func (db *Db) loadCrop() {
	if crops, ok := referencesData["crop"]; ok {
		for _, crop := range crops {
			c := Crop{
				Id:   crop.Id,
				Name: strings.ToUpper(crop.Name),
			}
			_, err := db.Conn.Model(&c).Insert()
			if pgerr, ok := err.(pg.Error); ok {
				if !pgerr.IntegrityViolation() {
					log.Errorf("error when trying to insert tree reference data : %s", pgerr.Error())
				}
			}
		}
	}
}

func (db *Db) PutBooklet(booklet Booklet) (err error) {
	_, err = db.Conn.Model(&booklet).Insert()
	pgErr, ok := err.(pg.Error)
	if ok && pgErr.IntegrityViolation() {
		switch {
		case strings.Contains(pgErr.Error(), "duplicate key value violates unique constraint \"booklets_number_key\""):
			return fmt.Errorf("this booklet number has already been registered")
		case strings.Contains(pgErr.Error(), "insert or update on table \"booklets\" violates foreign key constraint \"booklets_geo_code_id_fkey\""):
			return fmt.Errorf("this Geo code is unknown")
		default:
			return pgErr
		}
	} else {
		return err
	}
}

func (db *Db) UpdateEntities(booklet *Booklet, crate *Crate, shelf *Shelf, box *ArchiveBox) (err error) {
	return db.Conn.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		if booklet != nil {
			_, err := tx.Model(booklet).WherePK().Update()
			if err != nil {
				return err
			}
		}
		if crate != nil {
			_, err := tx.Model(crate).WherePK().Update()
			if err != nil {
				return err
			}
		}
		if shelf != nil {
			_, err := tx.Model(shelf).WherePK().Update()
			if err != nil {
				return err
			}
		}
		if box != nil {
			_, err := tx.Model(box).WherePK().Update()
			if err != nil {
				return err
			}
		}

		return nil
	})
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
