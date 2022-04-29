package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/oleiade/reflections"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

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
		_, err := db.Conn.Model((*Aggregates)(nil)).Exec(index)
		if err != nil {
			log.Error(err)
		}
	}
}

func (db *Db) createSchema() (err error) {
	for _, model := range []interface{}{
		(*Aggregates)(nil),
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
		columns = "hh_sno-c04gtrhh"
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
	query := fmt.Sprintf(`SELECT %s as data, rmo FROM aggregates where subpath(geocode, 0,%d) = ?;`, columns, count)
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
	from aggregates
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
from aggregates where subpath(geocode, 0,%d) = ?;`,
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
	from aggregates
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
	FROM aggregates
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
	FROM aggregates
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
	FROM aggregates
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
	FROM aggregates
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
	FROM aggregates
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
	FROM aggregates
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
	FROM aggregates
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
	FROM aggregates
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

type HouseholdLandInformation struct {
	Name                      string
	Column                    string
	NumberOfReportingHoldings uint
	NumberOfFarmHoldings      uint
	TotalAreaOfOwnLand        float64
	TotalFarmHoldingArea      float64
	AverageAreaPerFarmHolding float64
}

func (db *Db) GetHouseholdLandInformation(division, district, upazilla, union, mouza string) (data []HouseholdLandInformation, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	data = []HouseholdLandInformation{
		{
			Name:   "Own land",
			Column: "c04",
		},
		{
			Name:   "Given land",
			Column: "c05",
		},
		{
			Name:   "Taken land",
			Column: "c06",
		},
		{
			Name:   "Operated land",
			Column: "c07",
		},
		{
			Name:   "Homestead",
			Column: "c08",
		},
		{
			Name:   "Permanent Fellow land",
			Column: "c11",
		},
		{
			Name:   "Uncultivated land",
			Column: "c12",
		},
		{
			Name:   "Land under temporary crops",
			Column: "c13",
		},
		{
			Name:   "Land under permanent crops",
			Column: "c14",
		},
		{
			Name:   "Land under nursery",
			Column: "c16",
		},
		{
			Name:   "Land under current fallow",
			Column: "c17",
		},
		{
			Name:   "Total cultivated land",
			Column: "c18",
		},
		{
			Name:   "Land under irrigation",
			Column: "c19",
		},
		{
			Name:   "Land under salt cultivation",
			Column: "c20",
		},
	}

	for i, c := range data {
		query := fmt.Sprintf(`
		SELECT (
			SELECT count(hh_sno)
			FROM aggregates
			WHERE %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS number_of_reporting_holdings,
		(
			SELECT count(hh_sno)
			FROM aggregates
			WHERE c18 >= 0.05
				AND %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS number_of_farm_holdings,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS total_area_of_own_land,
		(
			SELECT sum(c18)
			FROM aggregates
			WHERE c18 >= 0.05
				AND %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS total_farm_holding_area;`,
			c.Column, count,
			c.Column, count,
			c.Column, count,
			c.Column, count)
		_, err = db.Conn.QueryOne(&c, query,
			geoCodeReq, geoCodeReq, geoCodeReq, geoCodeReq)
		if err != nil {
			log.Error(err)
			return
		}
		data[i] = c
	}
	return
}

func (db *Db) GetHouseholdFisheryLandInformation(division, district, upazilla, union, mouza string) (data []HouseholdLandInformation, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	data = []HouseholdLandInformation{
		{
			Name:   "Land under ponds/digi",
			Column: "c21",
		},
		{
			Name:   "Fishery Land other than ponds",
			Column: "(c22+c23+c24)",
		},
		{
			Name:   "Fishery Land under salt cultivation",
			Column: "c25",
		},
		{
			Name:   "Fishery Land cultivated under pan/cage",
			Column: "c26",
		},
		{
			Name:   "Fishery Land under fish cultivation by Creek",
			Column: "c27",
		},
	}

	for i, c := range data {
		query := fmt.Sprintf(`
		SELECT (
			SELECT count(hh_sno)
			FROM aggregates
			WHERE %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS number_of_reporting_holdings,
		(
			SELECT count(hh_sno)
			FROM aggregates
			WHERE c18 >= 0.05
				AND %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS number_of_farm_holdings,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS total_area_of_own_land,
		(
			SELECT sum(c18)
			FROM aggregates
			WHERE c18 >= 0.05
				AND %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS total_farm_holding_area;`,
			c.Column, count,
			c.Column, count,
			c.Column, count,
			c.Column, count)
		_, err = db.Conn.QueryOne(&c, query,
			geoCodeReq, geoCodeReq, geoCodeReq, geoCodeReq)
		if err != nil {
			log.Error(err)
			return
		}
		data[i] = c
	}
	return
}

type HouseholdPoultryInformation struct {
	Name                                     string
	Column                                   string
	NumberOfHouseholdPoultryColumn           string
	NumberOfHouseholdAttachFarmPoultryColumn string
	NumberOfReportingHoldings                uint
	TotalNumberOfPoultry                     uint
	NumberOfHouseholdPoultry                 uint
	NumberOfHouseholdAttachFarmPoultry       uint
	AverageTypeOfPoultryPerHolding           float64
}

func (db *Db) GetHouseholdPoultryInformation(division, district, upazilla, union, mouza string) (data []HouseholdPoultryInformation, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	data = []HouseholdPoultryInformation{
		{
			Name:                                     "Cock/Hen",
			Column:                                   "(c28h + c28f)",
			NumberOfHouseholdPoultryColumn:           "c28h",
			NumberOfHouseholdAttachFarmPoultryColumn: "c28f",
		},
		{
			Name:                                     "Duck",
			Column:                                   "(c29h + c29f)",
			NumberOfHouseholdPoultryColumn:           "c29h",
			NumberOfHouseholdAttachFarmPoultryColumn: "c29f",
		},
		{
			Name:                                     "Pigeon",
			Column:                                   "(c30h + c30f)",
			NumberOfHouseholdPoultryColumn:           "c30h",
			NumberOfHouseholdAttachFarmPoultryColumn: "c30f",
		},
		{
			Name:                                     "Quail",
			Column:                                   "(c31h + c31f)",
			NumberOfHouseholdPoultryColumn:           "c31h",
			NumberOfHouseholdAttachFarmPoultryColumn: "c31f",
		},
		{
			Name:                                     "Turkey",
			Column:                                   "(c32h + c32f)",
			NumberOfHouseholdPoultryColumn:           "c32h",
			NumberOfHouseholdAttachFarmPoultryColumn: "c32f",
		},
	}

	for i, c := range data {
		query := fmt.Sprintf(`
		SELECT (
			SELECT count(hh_sno)
			FROM aggregates
			WHERE %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS number_of_reporting_holdings,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS total_number_of_poultry,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS number_of_household_poultry,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS number_of_household_attach_farm_poultry;`,
			c.Column, count,
			c.Column, count,
			c.NumberOfHouseholdPoultryColumn, count,
			c.NumberOfHouseholdAttachFarmPoultryColumn, count)
		_, err = db.Conn.QueryOne(&c, query,
			geoCodeReq, geoCodeReq, geoCodeReq, geoCodeReq)
		if err != nil {
			log.Error(err)
			return
		}
		data[i] = c
	}
	return
}

type HouseholdCattleInformation struct {
	Name                                    string
	Column                                  string
	NumberOfHouseholdCattleColumn           string
	NumberOfHouseholdAttachFarmCattleColumn string
	NumberOfReportingHoldings               uint
	TotalNumberOfCattle                     uint
	NumberOfHouseholdCattle                 uint
	NumberOfHouseholdAttachFarmCattle       uint
	AverageTypeOfCattlePerHolding           float64
}

func (db *Db) GetHouseholdCattlenformation(division, district, upazilla, union, mouza string) (data []HouseholdCattleInformation, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	data = []HouseholdCattleInformation{
		{
			Name:                                    "Cow",
			Column:                                  "(c33h + c33f)",
			NumberOfHouseholdCattleColumn:           "c33h",
			NumberOfHouseholdAttachFarmCattleColumn: "c33f",
		},
		{
			Name:                                    "Buffalo",
			Column:                                  "(c34h + c34f)",
			NumberOfHouseholdCattleColumn:           "c34h",
			NumberOfHouseholdAttachFarmCattleColumn: "c34f",
		},
		{
			Name:                                    "Goat",
			Column:                                  "(c35h + c35f)",
			NumberOfHouseholdCattleColumn:           "c35h",
			NumberOfHouseholdAttachFarmCattleColumn: "c35f",
		},
		{
			Name:                                    "Sheep",
			Column:                                  "(c36h + c36f)",
			NumberOfHouseholdCattleColumn:           "c36h",
			NumberOfHouseholdAttachFarmCattleColumn: "c36f",
		},
		{
			Name:                                    "Pig",
			Column:                                  "(c37h + c37f)",
			NumberOfHouseholdCattleColumn:           "c37h",
			NumberOfHouseholdAttachFarmCattleColumn: "c37f",
		},
		{
			Name:                                    "Horse",
			Column:                                  "(c38h + c38f)",
			NumberOfHouseholdCattleColumn:           "c38h",
			NumberOfHouseholdAttachFarmCattleColumn: "c38f",
		},
	}

	for i, c := range data {
		query := fmt.Sprintf(`
		SELECT (
			SELECT count(hh_sno)
			FROM aggregates
			WHERE %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS number_of_reporting_holdings,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS total_number_of_cattle,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS number_of_household_cattle,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS number_of_household_attach_farm_cattle;`,
			c.Column, count,
			c.Column, count,
			c.NumberOfHouseholdCattleColumn, count,
			c.NumberOfHouseholdAttachFarmCattleColumn, count)
		_, err = db.Conn.QueryOne(&c, query,
			geoCodeReq, geoCodeReq, geoCodeReq, geoCodeReq)
		if err != nil {
			log.Error(err)
			return
		}
		data[i] = c
	}
	return
}

type HouseholdAgricultureEquipement struct {
	Name                              string
	NumberOfReportingHoldingsColumn   string
	NumberOfReportingHoldings         uint
	TotalNumberColumn                 string
	TotalNumber                       uint
	NumberOfNonMechanicalDeviceColumn string
	NumberOfNonMechanicalDevice       uint
	NumberOfDieselDeviceColumn        string
	NumberOfDieselDevice              uint
	NumberOfElectricalDeviceColumn    string
	NumberOfElectricalDevice          uint
}

func (db *Db) GetHouseholdAgricultureEquipement(division, district, upazilla, union, mouza string) (data []HouseholdAgricultureEquipement, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}
	data = []HouseholdAgricultureEquipement{
		{
			Name:                            "Tractor",
			NumberOfReportingHoldingsColumn: "c39",
			TotalNumberColumn:               "c39",
			NumberOfDieselDeviceColumn:      "c39",
		},
		{
			Name:                            "Power tiller",
			NumberOfReportingHoldingsColumn: "c40",
			TotalNumberColumn:               "c40",
			NumberOfDieselDeviceColumn:      "c40",
		},
		{
			Name:                            "Power pump",
			NumberOfReportingHoldingsColumn: "(c41a + c41b)",
			TotalNumberColumn:               "(c41a + c41b)",
			NumberOfDieselDeviceColumn:      "c41a",
			NumberOfElectricalDeviceColumn:  "c41b",
		},
		{
			Name:                            "Deep/Shallow tube well",
			NumberOfReportingHoldingsColumn: "(c42a + c42b)",
			TotalNumberColumn:               "(c42a + c42b)",
			NumberOfDieselDeviceColumn:      "c42a",
			NumberOfElectricalDeviceColumn:  "c42b",
		},
		{
			Name:                              "Crop planting machine",
			NumberOfReportingHoldingsColumn:   "(c43a + c43b)",
			TotalNumberColumn:                 "(c43a + c43b)",
			NumberOfNonMechanicalDeviceColumn: "c43a",
			NumberOfDieselDeviceColumn:        "c43b",
		},
		{
			Name:                              "Crop cutting machine",
			NumberOfReportingHoldingsColumn:   "(c44a + c44b)",
			TotalNumberColumn:                 "(c44a + c44b)",
			NumberOfNonMechanicalDeviceColumn: "c44a",
			NumberOfDieselDeviceColumn:        "c44b",
		},
		{
			Name:                              "Crop threshing machine",
			NumberOfReportingHoldingsColumn:   "(c45a + c45b + c45c)",
			TotalNumberColumn:                 "(c45a + c45b + c45c)",
			NumberOfNonMechanicalDeviceColumn: "c45a",
			NumberOfDieselDeviceColumn:        "c45b",
			NumberOfElectricalDeviceColumn:    "c45c",
		},
		{
			Name:                              "Fertilizer Appling machine",
			NumberOfReportingHoldingsColumn:   "(c46a + c46b)",
			TotalNumberColumn:                 "(c46a + c46b)",
			NumberOfNonMechanicalDeviceColumn: "c46a",
			NumberOfDieselDeviceColumn:        "c46b",
		},
		{
			Name:                              "Fish catching boat/trailer",
			NumberOfReportingHoldingsColumn:   "(c47a + c47b)",
			TotalNumberColumn:                 "(c47a + c47b)",
			NumberOfNonMechanicalDeviceColumn: "c47a",
			NumberOfDieselDeviceColumn:        "c47b",
		},
		{
			Name:                              "Fish catching net (business)",
			NumberOfReportingHoldingsColumn:   "c48",
			TotalNumberColumn:                 "c48",
			NumberOfNonMechanicalDeviceColumn: "c48",
		},
		{
			Name:                              "Plough",
			NumberOfReportingHoldingsColumn:   "c49",
			TotalNumberColumn:                 "c49",
			NumberOfNonMechanicalDeviceColumn: "c49",
		},
	}

	for i, c := range data {
		query := fmt.Sprintf(`
		SELECT (
			SELECT count(hh_sno)
			FROM aggregates
			WHERE %s > 0
				AND subpath(geocode, 0, %d) = ?
		) AS number_of_reporting_holdings,
		(
			SELECT sum(%s)
			FROM aggregates
			WHERE subpath(geocode, 0, %d) = ?
		) AS total_number`,
			c.NumberOfReportingHoldingsColumn, count,
			c.TotalNumberColumn, count)

		if c.NumberOfNonMechanicalDeviceColumn != "" {
			query += fmt.Sprintf(`
			,(
				SELECT sum(%s)
				FROM aggregates
				WHERE subpath(geocode, 0, %d) = ?
			) AS number_of_non_mechanical_device
			`, c.NumberOfNonMechanicalDeviceColumn, count)
		}
		if c.NumberOfDieselDeviceColumn != "" {
			query += fmt.Sprintf(`
			,(
				SELECT sum(%s)
				FROM aggregates
				WHERE subpath(geocode, 0, %d) = ?
			) AS number_of_diesel_device
			`, c.NumberOfDieselDeviceColumn, count)
		}
		if c.NumberOfElectricalDeviceColumn != "" {
			query += fmt.Sprintf(`
			,(
				SELECT sum(%s)
				FROM aggregates
				WHERE subpath(geocode, 0, %d) = ?
			) AS number_of_electrical_device
			`, c.NumberOfElectricalDeviceColumn, count)
		}
		_, err = db.Conn.QueryOne(&c, query,
			geoCodeReq, geoCodeReq, geoCodeReq, geoCodeReq, geoCodeReq)
		if err != nil {
			log.Error(err)
			return
		}
		data[i] = c
	}
	return
}

type TemporaryCrops struct {
	NumberOfFarmHoldings uint
	CropArea             float64
	T101                 float64
	T102                 float64
	T103                 float64
	T104                 float64
	T105                 float64
	T106                 float64
	T107                 float64
	T108                 float64
	T109                 float64
	T110                 float64
	T111                 float64
	T112                 float64
	T113                 float64
	T114                 float64
	T115                 float64
	T116                 float64
	T117                 float64
	T118                 float64
	T119                 float64
	T120                 float64
	T121                 float64
	T122                 float64
	T123                 float64
	T124                 float64
	T125                 float64
	T126                 float64
	T127                 float64
	T128                 float64
	T129                 float64
	T130                 float64
	T131                 float64
	T132                 float64
	T133                 float64
	T134                 float64
	T135                 float64
	T136                 float64
	T137                 float64
	T138                 float64
	T139                 float64
	T140                 float64
	T141                 float64
	T142                 float64
	T143                 float64
	T144                 float64
	T145                 float64
	T146                 float64
	T147                 float64
	T148                 float64
	T149                 float64
	T150                 float64
	T151                 float64
	T152                 float64
	T153                 float64
	T154                 float64
	T155                 float64
	T156                 float64
	T157                 float64
	T158                 float64
	T159                 float64
	T160                 float64
	T161                 float64
	T162                 float64
	T163                 float64
	T164                 float64
	T165                 float64
	T166                 float64
	T167                 float64
	T168                 float64
	T169                 float64
	T170                 float64
	T171                 float64
	T172                 float64
	T173                 float64
	T174                 float64
	T175                 float64
	T176                 float64
	T177                 float64
	T178                 float64
	T179                 float64
	T180                 float64
	T181                 float64
	T182                 float64
	T183                 float64
	T184                 float64
	T185                 float64
	T186                 float64
	T187                 float64
	T188                 float64
	T189                 float64
	T190                 float64
	T191                 float64
	T192                 float64
	T193                 float64
	T194                 float64
	T195                 float64
	T196                 float64
	T197                 float64
	T198                 float64
	T199                 float64
	T200                 float64
	T201                 float64
	T202                 float64
	T203                 float64
}

func (c TemporaryCrops) PercentageOfCropArea(cropArea string) string {
	p := message.NewPrinter(language.English)
	r, err := reflections.GetField(c, cropArea)
	if err != nil {
		log.Errorf("unable to get %s field from Crops struct : %s", cropArea, err)
		return "err"
	}
	return p.Sprintf("%f%%", (r.(float64)/
		(c.T101+
			c.T102+
			c.T103+
			c.T104+
			c.T105+
			c.T106+
			c.T107+
			c.T108+
			c.T109+
			c.T110+
			c.T111+
			c.T112+
			c.T113+
			c.T114+
			c.T115+
			c.T116+
			c.T117+
			c.T118+
			c.T119+
			c.T120+
			c.T121+
			c.T122+
			c.T123+
			c.T124+
			c.T125+
			c.T126+
			c.T127+
			c.T128+
			c.T129+
			c.T130+
			c.T131+
			c.T132+
			c.T133+
			c.T134+
			c.T135+
			c.T136+
			c.T137+
			c.T138+
			c.T139+
			c.T140+
			c.T141+
			c.T142+
			c.T143+
			c.T144+
			c.T145+
			c.T146+
			c.T147+
			c.T148+
			c.T149+
			c.T150+
			c.T151+
			c.T152+
			c.T153+
			c.T154+
			c.T155+
			c.T156+
			c.T157+
			c.T158+
			c.T159+
			c.T160+
			c.T161+
			c.T162+
			c.T163+
			c.T164+
			c.T165+
			c.T166+
			c.T167+
			c.T168+
			c.T169+
			c.T170+
			c.T171+
			c.T172+
			c.T173+
			c.T174+
			c.T175+
			c.T176+
			c.T177+
			c.T178+
			c.T179+
			c.T180+
			c.T181+
			c.T182+
			c.T183+
			c.T184+
			c.T185+
			c.T186+
			c.T187+
			c.T188+
			c.T189+
			c.T190+
			c.T191+
			c.T192+
			c.T193+
			c.T194+
			c.T195+
			c.T196+
			c.T197+
			c.T198+
			c.T199+
			c.T200+
			c.T201+
			c.T202+
			c.T203))*100)
}

func (db *Db) GetTemporaryCrops(division, district, upazilla, union, mouza string) (data TemporaryCrops, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}

	query := fmt.Sprintf(`
	SELECT 
    	sum(sf + mf + lf) as number_of_farm_holdings,
    	sum(c13) as crop_area,
		sum(t101) as t101,
		sum(t102) as t102,
		sum(t103) as t103,
		sum(t104) as t104,
		sum(t105) as t105,
		sum(t106) as t106,
		sum(t107) as t107,
		sum(t108) as t108,
		sum(t109) as t109,
		sum(t110) as t110,
		sum(t111) as t111,
		sum(t112) as t112,
		sum(t113) as t113,
		sum(t114) as t114,
		sum(t115) as t115,
		sum(t116) as t116,
		sum(t117) as t117,
		sum(t118) as t118,
		sum(t119) as t119,
		sum(t120) as t120,
		sum(t121) as t121,
		sum(t122) as t122,
		sum(t123) as t123,
		sum(t124) as t124,
		sum(t125) as t125,
		sum(t126) as t126,
		sum(t127) as t127,
		sum(t128) as t128,
		sum(t129) as t129,
		sum(t130) as t130,
		sum(t131) as t131,
		sum(t132) as t132,
		sum(t133) as t133,
		sum(t134) as t134,
		sum(t135) as t135,
		sum(t136) as t136,
		sum(t137) as t137,
		sum(t138) as t138,
		sum(t139) as t139,
		sum(t140) as t140,
		sum(t141) as t141,
		sum(t142) as t142,
		sum(t143) as t143,
		sum(t144) as t144,
		sum(t145) as t145,
		sum(t146) as t146,
		sum(t147) as t147,
		sum(t148) as t148,
		sum(t149) as t149,
		sum(t150) as t150,
		sum(t151) as t151,
		sum(t152) as t152,
		sum(t153) as t153,
		sum(t154) as t154,
		sum(t155) as t155,
		sum(t156) as t156,
		sum(t157) as t157,
		sum(t158) as t158,
		sum(t159) as t159,
		sum(t160) as t160,
		sum(t161) as t161,
		sum(t162) as t162,
		sum(t163) as t163,
		sum(t164) as t164,
		sum(t165) as t165,
		sum(t166) as t166,
		sum(t167) as t167,
		sum(t168) as t168,
		sum(t169) as t169,
		sum(t170) as t170,
		sum(t171) as t171,
		sum(t172) as t172,
		sum(t173) as t173,
		sum(t174) as t174,
		sum(t175) as t175,
		sum(t176) as t176,
		sum(t177) as t177,
		sum(t178) as t178,
		sum(t179) as t179,
		sum(t180) as t180,
		sum(t181) as t181,
		sum(t182) as t182,
		sum(t183) as t183,
		sum(t184) as t184,
		sum(t185) as t185,
		sum(t186) as t186,
		sum(t187) as t187,
		sum(t188) as t188,
		sum(t189) as t189,
		sum(t190) as t190,
		sum(t191) as t191,
		sum(t192) as t192,
		sum(t193) as t193,
		sum(t194) as t194,
		sum(t195) as t195,
		sum(t196) as t196,
		sum(t197) as t197,
		sum(t198) as t198,
		sum(t199) as t199,
		sum(t200) as t200,
		sum(t201) as t201,
		sum(t202) as t202,
		sum(t203) as t203
	FROM aggregates
	WHERE subpath(geocode, 0, %d) = ?;
	`, count)

	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return

}

type PermanantCrops struct {
	NumberOfFarmHoldings uint
	CropArea             float64
	T101                 float64
	T102                 float64
	T103                 float64
	T104                 float64
	T105                 float64
	T106                 float64
	T107                 float64
	T108                 float64
	T109                 float64
	T110                 float64
	T111                 float64
	T112                 float64
	T113                 float64
	T114                 float64
	T115                 float64
	T116                 float64
	T117                 float64
	T118                 float64
	T119                 float64
	T120                 float64
	T121                 float64
	T122                 float64
	T123                 float64
	T124                 float64
	T125                 float64
	T126                 float64
	T127                 float64
	T128                 float64
	T129                 float64
	T130                 float64
	T131                 float64
	T132                 float64
	T133                 float64
	T134                 float64
	T135                 float64
	T136                 float64
	T137                 float64
	T138                 float64
	T139                 float64
	T140                 float64
	T141                 float64
	T142                 float64
	T143                 float64
	T144                 float64
	T145                 float64
	T146                 float64
	T147                 float64
	T148                 float64
	T149                 float64
	T150                 float64
	T151                 float64
	T152                 float64
	T153                 float64
	T154                 float64
	T155                 float64
	T156                 float64
	T157                 float64
	T158                 float64
	T159                 float64
	T160                 float64
	T161                 float64
	T162                 float64
	T163                 float64
	T164                 float64
	T165                 float64
	T166                 float64
	T167                 float64
	T168                 float64
	T169                 float64
	T170                 float64
	T171                 float64
	T172                 float64
	T173                 float64
	T174                 float64
	T175                 float64
	T176                 float64
	T177                 float64
	T178                 float64
	T179                 float64
	T180                 float64
	T181                 float64
	T182                 float64
	T183                 float64
	T184                 float64
	T185                 float64
	T186                 float64
	T187                 float64
	T188                 float64
	T189                 float64
	T190                 float64
	T191                 float64
	T192                 float64
	T193                 float64
	T194                 float64
	T195                 float64
	T196                 float64
	T197                 float64
	T198                 float64
	T199                 float64
	T200                 float64
	T201                 float64
	T202                 float64
	T203                 float64
}

func (c TemporaryCrops) PercentageOfPermantCropArea(cropArea string) string {
	p := message.NewPrinter(language.English)
	r, err := reflections.GetField(c, cropArea)
	if err != nil {
		log.Errorf("unable to get %s field from Crops struct : %s", cropArea, err)
		return "err"
	}
	return p.Sprintf("%f%%", (r.(float64)/
		(c.T101+
			c.T102+
			c.T103+
			c.T104+
			c.T105+
			c.T106+
			c.T107+
			c.T108+
			c.T109+
			c.T110+
			c.T111+
			c.T112+
			c.T113+
			c.T114+
			c.T115+
			c.T116+
			c.T117+
			c.T118+
			c.T119+
			c.T120+
			c.T121+
			c.T122+
			c.T123+
			c.T124+
			c.T125+
			c.T126+
			c.T127+
			c.T128+
			c.T129+
			c.T130+
			c.T131+
			c.T132+
			c.T133+
			c.T134+
			c.T135+
			c.T136+
			c.T137+
			c.T138+
			c.T139+
			c.T140+
			c.T141+
			c.T142+
			c.T143+
			c.T144+
			c.T145+
			c.T146+
			c.T147+
			c.T148+
			c.T149+
			c.T150+
			c.T151+
			c.T152+
			c.T153+
			c.T154+
			c.T155+
			c.T156+
			c.T157+
			c.T158+
			c.T159+
			c.T160+
			c.T161+
			c.T162+
			c.T163+
			c.T164+
			c.T165+
			c.T166+
			c.T167+
			c.T168+
			c.T169+
			c.T170+
			c.T171+
			c.T172+
			c.T173+
			c.T174+
			c.T175+
			c.T176+
			c.T177+
			c.T178+
			c.T179+
			c.T180+
			c.T181+
			c.T182+
			c.T183+
			c.T184+
			c.T185+
			c.T186+
			c.T187+
			c.T188+
			c.T189+
			c.T190+
			c.T191+
			c.T192+
			c.T193+
			c.T194+
			c.T195+
			c.T196+
			c.T197+
			c.T198+
			c.T199+
			c.T200+
			c.T201+
			c.T202+
			c.T203))*100)
}

func (db *Db) GetPermanantCrops(division, district, upazilla, union, mouza string) (data PermanantCrops, err error) {
	geoCodeReq, count, err := GetGeoRequest(division, district, upazilla, union, mouza)
	if err != nil {
		return
	}

	query := fmt.Sprintf(`
	SELECT 
    	sum(sf + mf + lf) as number_of_farm_holdings,
    	sum(c13) as crop_area,
		sum(t101) as t101,
		sum(t102) as t102,
		sum(t103) as t103,
		sum(t104) as t104,
		sum(t105) as t105,
		sum(t106) as t106,
		sum(t107) as t107,
		sum(t108) as t108,
		sum(t109) as t109,
		sum(t110) as t110,
		sum(t111) as t111,
		sum(t112) as t112,
		sum(t113) as t113,
		sum(t114) as t114,
		sum(t115) as t115,
		sum(t116) as t116,
		sum(t117) as t117,
		sum(t118) as t118,
		sum(t119) as t119,
		sum(t120) as t120,
		sum(t121) as t121,
		sum(t122) as t122,
		sum(t123) as t123,
		sum(t124) as t124,
		sum(t125) as t125,
		sum(t126) as t126,
		sum(t127) as t127,
		sum(t128) as t128,
		sum(t129) as t129,
		sum(t130) as t130,
		sum(t131) as t131,
		sum(t132) as t132,
		sum(t133) as t133,
		sum(t134) as t134,
		sum(t135) as t135,
		sum(t136) as t136,
		sum(t137) as t137,
		sum(t138) as t138,
		sum(t139) as t139,
		sum(t140) as t140,
		sum(t141) as t141,
		sum(t142) as t142,
		sum(t143) as t143,
		sum(t144) as t144,
		sum(t145) as t145,
		sum(t146) as t146,
		sum(t147) as t147,
		sum(t148) as t148,
		sum(t149) as t149,
		sum(t150) as t150,
		sum(t151) as t151,
		sum(t152) as t152,
		sum(t153) as t153,
		sum(t154) as t154,
		sum(t155) as t155,
		sum(t156) as t156,
		sum(t157) as t157,
		sum(t158) as t158,
		sum(t159) as t159,
		sum(t160) as t160,
		sum(t161) as t161,
		sum(t162) as t162,
		sum(t163) as t163,
		sum(t164) as t164,
		sum(t165) as t165,
		sum(t166) as t166,
		sum(t167) as t167,
		sum(t168) as t168,
		sum(t169) as t169,
		sum(t170) as t170,
		sum(t171) as t171,
		sum(t172) as t172,
		sum(t173) as t173,
		sum(t174) as t174,
		sum(t175) as t175,
		sum(t176) as t176,
		sum(t177) as t177,
		sum(t178) as t178,
		sum(t179) as t179,
		sum(t180) as t180,
		sum(t181) as t181,
		sum(t182) as t182,
		sum(t183) as t183,
		sum(t184) as t184,
		sum(t185) as t185,
		sum(t186) as t186,
		sum(t187) as t187,
		sum(t188) as t188,
		sum(t189) as t189,
		sum(t190) as t190,
		sum(t191) as t191,
		sum(t192) as t192,
		sum(t193) as t193,
		sum(t194) as t194,
		sum(t195) as t195,
		sum(t196) as t196,
		sum(t197) as t197,
		sum(t198) as t198,
		sum(t199) as t199,
		sum(t200) as t200,
		sum(t201) as t201,
		sum(t202) as t202,
		sum(t203) as t203
	FROM aggregates
	WHERE subpath(geocode, 0, %d) = ?;
	`, count)

	_, err = db.Conn.QueryOne(&data, query,
		geoCodeReq)
	if err != nil {
		log.Error(err)
		return
	}
	return

}
