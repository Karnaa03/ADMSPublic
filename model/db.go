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

	// sex = male
	// sex2 = female
	// sex3 = hijra

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
		columns = "sum(sex)"
		conditions = "true = true GROUP BY rmo"
	case "16":
		columns = "sum(sex2)"
		conditions = "true = true GROUP BY rmo"
	case "17":
		columns = "count(c19)"
		conditions = "c19 > 0 GROUP BY rmo"
	case "18":
		columns = "sum(c19)"
		conditions = "true = true GROUP BY rmo"
	case "19":
		columns = "count(sf+lf+mf)"
		conditions = "c19 is not null GROUP BY rmo"
	case "20":
		columns = "sum(c19)"
		conditions = "c18 >= 0.05 GROUP BY rmo"
	case "21":
		columns = "c33h+c33f"
	case "22":
		columns = "c34h+c34f"
	case "23":
		columns = "sum(c33h+c33f)"
		conditions = "c18 >= 0.05 GROUP BY rmo"
	case "24":
		columns = "sum(c34h+c34f)"
		conditions = "c18 >= 0.05 GROUP BY rmo"
	case "25":
		columns = "sum(c35h+c35f)"
		conditions = "true = true GROUP BY rmo"
	case "26":
		columns = "sum(c36h+c36f)"
		conditions = "true = true GROUP BY rmo"
	case "27":
		columns = "sum(c35h+c35f)"
		conditions = "c18 >= 0.05 GROUP BY rmo"
	case "28":
		columns = "sum(c36h+c36f)"
		conditions = "true = true GROUP BY rmo"
	case "29":
		columns = "sum(c28h+c28f)"
		conditions = "true = true GROUP BY rmo"
	case "30":
		columns = "sum(c29h+c29f)"
		conditions = "true = true GROUP BY rmo"
	case "31":
		columns = "sum(c28h+c28f)"
		conditions = "c18 >= 0.05 GROUP BY rmo"
	case "32":
		columns = "sum(c29h+c29f)"
		conditions = "c18 >= 0.05 GROUP BY rmo"
	case "33":
		columns = "sum(t101+t102+t103+t104+t105+t112+t113+t114+t121+t122+t123+t124+t125+t127+t128+t129+t130+t131+t132+t134+t135+t157+t158+t159+t160+t161+t167+t169+t175+t176+t177+t179+t182+t185+t106+t107+t108+t109+t110+t111+t115+t116+t117+t118+t119+t120+t126+t133+t136+t137+t138+t139+t140+t141+t142+t143+t144+t145+t146+t147+t148+t149+t150+t151+t152+t153+t154+t155+t156+t162+t163+t164+t165+t166+t168+t170+t171+t172+t173+t174+t178+t180+t181+t183+t184+t186+t187+t188+t189+t190+t191+t192+t193+t194+t195+t196+t197+t198+t199+t200+t201+t202+t203)/c13"
		conditions = "true = true GROUP BY rmo, c13"
	case "34":
		columns = "sum(c13)"

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

type OccupationHouseHoldHead struct {
	Occ   uint
	Occ2  uint
	Occ3  uint
	Occ4  uint
	Occ5  uint
	Total uint
}

func (db *Db) GetOccupationOfHouseHold(division, district, upazilla, union, mouza string) (data OccupationHouseHoldHead, err error) {
	geoCodeReq, count := GetGeoRequest(division, district, upazilla, union, mouza)
	query := fmt.Sprintf("select sum(occ) as occ,sum(occ2) as occ2,sum(occ3) as occ3,sum(occ4) as occ4,sum(occ5) as occ5,(sum(occ)+sum(occ2)+sum(occ3)+sum(occ4)+sum(occ5)) as total from agregateds where subpath(geocode, 0,%d) = ?;",
		count)
	_, err = db.Conn.QueryOne(&data, query,
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
