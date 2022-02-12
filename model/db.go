package model

import (
	"context"
	"fmt"
	"strconv"
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

func GetGeoRequest(division, district, upazilla, union, mouza string) (selector string, count uint, err error) {
	if division != "" {
		var divisionInt int
		divisionInt, err = strconv.Atoi(division)
		if err != nil {
			return
		}
		selector = fmt.Sprintf("%02d", divisionInt)
		count = 1
	} else {
		return
	}
	if district != "" {
		var districtInt int
		districtInt, err = strconv.Atoi(district)
		selector += "." + fmt.Sprintf("%02d", districtInt)
		count = 2
	} else {
		return
	}
	if upazilla != "" {
		var upazillaInt int
		upazillaInt, err = strconv.Atoi(upazilla)
		selector += "." + fmt.Sprintf("%02d", upazillaInt)
		count = 3
	} else {
		return
	}
	if union != "" {
		var unionInt int
		unionInt, err = strconv.Atoi(union)
		selector += "." + fmt.Sprintf("%03d", unionInt)
		count = 4
	} else {
		return
	}
	if mouza != "" {
		var mouzaInt int
		mouzaInt, err = strconv.Atoi(mouza)
		selector += "." + fmt.Sprintf("%03d", mouzaInt)
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
		conditions = "true = true GROUP BY rmo"
	default:
		return tableData, fmt.Errorf(("don't know this table name"))
	}

	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
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
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	select sum(occ) as occ,
    sum(occ2) as occ2,
    sum(occ3) as occ3,
    sum(occ4) as occ4,
    sum(occ5) as occ5,
	(sum(occ) + sum(occ2) + sum(occ3) + sum(occ4) + sum(occ5)) as total
	from agregateds
	where subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

type EducationOfTheHouseholdHead struct {
	NoEducation               uint
	Class1                    uint
	Class2                    uint
	Class3                    uint
	Class4                    uint
	Class5                    uint
	Class6                    uint
	Class7                    uint
	Class8                    uint
	Class9                    uint
	Ssc                       uint
	Hsc                       uint
	BachelorEquivalent        uint
	MastersEquivalentOrHigher uint
	Total                     uint
}

func (db *Db) GetEducationOfTheHouseholdHead(division, district, upazilla, union, mouza string) (data EducationOfTheHouseholdHead, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	select 
	sum(edu) as no_education,
    sum(edu1) as Class1,
    sum(edu2) as Class2,
    sum(edu3) as Class3,
    sum(edu4) as Class4,
    sum(edu5) as Class5,
    sum(edu6) as Class6,
    sum(edu7) as Class7,
    sum(edu8) as Class8,
    sum(edu9) as Class9, 
    sum(edu10) as Ssc,
    sum(edu12) as Hsc,
    sum(edu15) as Bachelor_equivalent,
    sum(edu18) as Masters_Equivalent_Or_Higher,
    (
        sum(edu1) + 
		sum(edu2) + 
		sum(edu3) + 
		sum(edu4) + 
		sum(edu5) + 
		sum(edu6) + 
		sum(edu7) + 
		sum(edu8) + 
		sum(edu9) + 
		sum(edu10) + 
		sum(edu12) + 
		sum(edu15) + 
		sum(edu18)
    ) as Total
from agregateds where subpath(geocode, 0,%d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

type GenderOfTheHouseholdHead struct {
	Male   uint
	Female uint
	Hijra  uint
	Total  uint
}

func (db *Db) GetGenderOfTheHouseholdHead(division, district, upazilla, union, mouza string) (data GenderOfTheHouseholdHead, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	select sum(sex) as male,
    sum(sex2) as female,
    sum(sex3) as hijra,
    (sum(sex) + sum(sex2) + sum(sex3)) as total
	from agregateds
	where subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

type FisheryHolding struct {
	NumberOfFisheryHousehold uint
	Percentage               float64
}

func (db *Db) GetFisheryHolding(division, district, upazilla, union, mouza string) (data FisheryHolding, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	SELECT sum(hh_f) as Number_Of_Fishery_Household,
    	(sum(hh_f)::NUMERIC / sum(hh_sno)::NUMERIC)::NUMERIC * 100 as Percentage
	FROM agregateds
	WHERE hh_f = 1 AND subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

type AgriculuralLaborHolding struct {
	NumberOfAgriLaborHouseHold uint
	Percentage                 float64
}

func (db *Db) GetAgriculuralLaborHolding(division, district, upazilla, union, mouza string) (data AgriculuralLaborHolding, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	SELECT sum(hh_a) as Number_Of_Agri_Labor_House_Hold,
    	(sum(hh_f)::NUMERIC / sum(hh_sno)::NUMERIC)::NUMERIC * 100 as Percentage
	FROM agregateds
	WHERE hh_a = 1 AND subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

type HouseholdHeadInformation struct {
	NoEducation                   uint
	Class_I_V                     uint
	Class_VI_IX                   uint
	SccPassed                     uint
	HscPassed                     uint
	DegreePassed                  uint
	MasterPassed                  uint
	TotalEducation                uint
	Agriculture                   uint
	Industry                      uint
	Service                       uint
	Business                      uint
	Other                         uint
	TotalOccupation               uint
	FisheryHolding                uint
	FisheryHoldingPercentage      float64
	AgriculturalHolding           uint
	AgriculturalHoldingPercentage float64
	HouseholdMemberMale           uint
	HouseholdMemberFemale         uint
	HouseholdMemberHijra          uint
	HouseholdMemberTotal          uint
	HouseholdWorkerMale           uint
	HouseholdWorkerFemale         uint
	HouseholdWorkerHijra          uint
	HouseholdWorkerTotal          uint
	HouseholdWorker_10_14Male     uint
	HouseholdWorker_10_14Female   uint
	HouseholdWorker_10_14Hijra    uint
	HouseholdWorker_10_14Total    uint
	HouseholdWorker_15PlusMale    uint
	HouseholdWorker_15PlusFemale  uint
	HouseholdWorker_15PlusHijra   uint
	HouseholdWorker_15PlusTotal   uint
}

func (db *Db) GetHouseholdHeadInformation(division, district, upazilla, union, mouza string) (data HouseholdHeadInformation, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	SELECT sum(edu) as no_education,
    	(sum(edu1) + sum(edu2) + sum(edu3) + sum(edu4) + sum(edu5)) as class_I_V,
    	(sum(edu6) + sum(edu7) + sum(edu8) + sum(edu9)) as class_VI_IX,
    	sum(edu10) as Scc_Passed,
    	sum(edu12) as Hsc_Passed,
    	sum(edu15) as Degree_Passed,
    	sum(edu18) as Master_Passed,
    	(
    	    sum(edu) + sum(edu1) + sum(edu2) + sum(edu3) + sum(edu4) + sum(edu5) + sum(edu6) + sum(edu7) + sum(edu8) + sum(edu9) + sum(edu10) + sum(edu12) + sum(edu15) + sum(edu18)
    	) as Total_Education,
    	sum(occ) as Agriculture,
    	sum(occ2) as Industry,
    	sum(occ3) as Service,
    	sum(occ4) as Business,
    	sum(occ5) as Other,
    	(sum(occ) + sum(occ2) + sum(occ3) + sum(occ4) + sum(occ5)) as Total_Occupation,
    	sum(c01m) as Household_Member_Male,
    	sum(c01f) as Household_Member_Female,
    	sum(c01h) as Household_Member_Hijra,
    	(sum(c01m) + sum(c01f) + sum(c01h)) as Household_Member_Total,
    	(sum(c02m) + sum(c03m)) as Household_Worker_Male,
    	(sum(c02f) + sum(c03f)) as Household_Worker_Female,
    	(sum(c02h) + sum(c03h)) as Household_Worker_Hijra,
    	(sum(c02m) + sum(c03m) + sum(c02f) + sum(c03f) + sum(c02h) + sum(c03h)) as Household_Worker_Total,
    	sum(c02m) as Household_Worker_10_14_Male,
    	sum(c02f) as Household_Worker_10_14_Female,
    	sum(c02h) as Household_Worker_10_14_Hijra,
    	(sum(c02m) + sum(c02f) + sum(c02h)) as Household_Worker_10_14_Total,
    	sum(c03m) as Household_Worker_15_Plus_Male,
    	sum(c03f) as Household_Worker_15_Plus_Female,
    	sum(c03h) as Household_Worker_15_Plus_Hijra,
    	(sum(c03m) + sum(c03f) + sum(c03h)) as Household_Worker_15_Plus_Total
	FROM agregateds
	WHERE hh_a = 1 AND subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

type OccupationOfHouseholdHead struct {
	Agriculture uint
	Industry    uint
	Service     uint
	Business    uint
	Others      uint
	Total       uint
}

func (db *Db) GetOccupationOfHouseholdHead(division, district, upazilla, union, mouza string) (data OccupationOfHouseholdHead, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	SELECT sum(occ) as Agriculture,
    	sum(occ2) as Industry,
    	sum(occ3) as Service,
    	sum(occ4) as Business,
    	sum(occ5) as Others,
    	(
    	    sum(occ) + sum(occ2) + sum(occ3) + sum(occ4) + sum(occ5)
    	) as Total
	FROM agregateds
	WHERE hh_a = 1 AND subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

type TotalNumberOfHouseholdMembers struct {
	Male   uint
	Female uint
	Hijra  uint
	Total  uint
}

func (db *Db) GetTotalNumberOfHouseholdMembers(division, district, upazilla, union, mouza string) (data TotalNumberOfHouseholdMembers, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	SELECT sum(c01m) as Male,
    	sum(c01f) as Female,
    	sum(c01h) as Hijra,
    	(sum(c01m) + sum(c01f) + sum(c01h)) as Total
	FROM agregateds
	WHERE hh_a = 1 AND subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (db *Db) GetTotalNumberOfHouseholdWorkers(division, district, upazilla, union, mouza string) (data TotalNumberOfHouseholdMembers, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	SELECT sum(c02m + c03m) as Male,
    	sum(c02f + c03f) as Female,
    	sum(c02h + c03h) as Hijra,
    	(sum(c02m + c03m) + sum(c02f + c03f) + sum(c02h + c03h)) as Total
	FROM agregateds
	WHERE hh_a = 1 AND subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (db *Db) GetTotalNumberOfHouseholdWorkers1014(division, district, upazilla, union, mouza string) (data TotalNumberOfHouseholdMembers, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	SELECT sum(c02m) as Male,
    	sum(c02f) as Female,
    	sum(c02h) as Hijra,
    	(sum(c02m) + sum(c02f) + sum(c02h)) as Total
	FROM agregateds
	WHERE hh_a = 1 AND subpath(geocode, 0, %d) = ?;`,
		count)
	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (db *Db) GetTotalNumberOfHouseholdWorkers15plus(division, district, upazilla, union, mouza string) (data TotalNumberOfHouseholdMembers, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`
	SELECT sum(c03m) as Male,
    	sum(c03f) as Female,
    	sum(c03h) as Hijra,
    	(sum(c03m) + sum(c03f) + sum(c03h)) as Total
	FROM agregateds
	WHERE hh_a = 1 AND subpath(geocode, 0, %d) = ?;`,
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
