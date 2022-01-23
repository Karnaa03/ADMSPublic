package model

import (
	"context"
	"fmt"
	"strings"

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

	}
	return
}

func (db *Db) Close() {
	err := db.Conn.Close()
	if err != nil {
		log.Error(err)
	}
}

func (db *Db) createIndex() {
	for name, index := range index {
		log.Infof("create index : %s", name)
		_, err := db.Conn.Model((*Agregated)(nil)).Exec(index)
		if err != nil {
			log.Error(err)
		}
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

type RawTableData struct {
	Data float64
	Rmo  uint
}

func GetGeoRequest(division, district, upazilla, union, mouza string) (selector string, count uint) {
	if division != "" {
		selector = division
		count = 1
	} else {
		return
	}
	if district != "" {
		selector += "." + district
		count = 2
	} else {
		return
	}
	if upazilla != "" {
		selector += "." + upazilla
		count = 3
	} else {
		return
	}
	if union != "" {
		selector += "." + union
		count = 4
	} else {
		return
	}
	if mouza != "" {
		selector += "." + mouza
		count = 5
	} else {
		return
	}
	return
}

func (db *Db) GetAgregate(division, district, upazilla, union, mouza, tableName string) (tableData []RawTableData, err error) {
	columns := ""
	conditions := ""
	switch tableName {
	case "1":
		columns = "hh_sno"
	case "2":
		columns = "sf+mf+lf"
	case "3":
		columns = "sf"
	case "4":
		columns = "mf"
	case "5":
		columns = "lf"
	case "6":
		columns = "hh_sno"
		conditions = "c04 = 0"
	case "7":
		columns = "sum(hh_a)"
		conditions = "hh_a = 1 GROUP BY rmo"
	case "8":
		columns = "sum(hh_f)"
		conditions = "hh_f = 1 GROUP BY rmo"
	case "9":
		columns = "c02m+c02f+c02h+c03m+c03f+c03h"
	case "10":
		columns = "c02m+c02f+c02h+c03m+c03f+c03h "
		conditions = "c18 >= 0.05"
	case "11":
		columns = "sum(c07)"
		conditions = "true = true GROUP BY rmo"
	case "12":
		columns = "sum(c07)"
		conditions = "c18 >= 0.05 GROUP BY rmo"
	case "13":
		columns = "sum(c08)"
		conditions = "true = true GROUP BY rmo"
	case "14":
		columns = "sum(c08)"
		conditions = "c18 >= 0.05 GROUP BY rmo"
	case "15":

	default:
		return tableData, fmt.Errorf(("don't know this table name"))
	}

	geoCodeReq, count := GetGeoRequest(division, district, upazilla, union, mouza)
	query := fmt.Sprintf(`SELECT %s as data, rmo FROM agregateds where subpath(geocode, 0,%d) = ?;`, columns, count)
	if conditions != "" {
		query = strings.Replace(query, ";", fmt.Sprintf(" AND %s;", conditions), 1)
	}
	_, err = db.Conn.Query(&tableData, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (db *Db) GetGeoCode(geoCodeNumber string) (geoCode GeoCodes, err error) {
	geoCode = GeoCodes{
		GeocodeID: geoCodeNumber,
	}
	err = db.Conn.Model(&geoCode).WherePK().Select()
	return
}
