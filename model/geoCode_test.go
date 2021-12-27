package model

import (
	"encoding/json"
	"testing"

	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/conf"
)

var db Db

func init() {
	config := conf.Config{
		Version: "test",
	}
	config.Load()

	db = Db{}
	err := db.Init(config)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGeocodesJson(t *testing.T) {
	g := GeoCodes{
		GeocodeID: "72.04.715.151.00.1",
		Division:  72,
		District:  04,
		Upazilla:  4,
		Union:     151,
		Mouza:     0,
		Villages: &Villages{Villages: map[int]Village{
			001: {NameEn: "test1", NameBengali: "benga1"},
			002: {NameEn: "test2", NameBengali: "benga2"},
		},
		},
		CA:               0,
		Rmo:              1,
		NameDivision:     "",
		NameDistrict:     "",
		NameUpazilla:     "",
		NameUnion:        "",
		NameMouza:        "",
		NameCountingArea: "",
		NameRMO:          "",
	}
	gjson, err := json.MarshalIndent(g, " ", "  ")
	if err != nil {
		t.Error(err)
	}
	t.Logf(string(gjson))
	fjson, err := json.Marshal(g)
	if err != nil {
		t.Error(err)
	}
	_, err = db.Conn.Model(&g).Insert()
	if err != nil {
		t.Error(err)
	}

	t.Log(string(fjson))
}

func TestGetGeoCode(t *testing.T) {
	g := GeoCodes{
		GeocodeID: "70.66.187.127.00.1",
	}
	err := db.Conn.Model(&g).WherePK().Select()
	if err != nil {
		t.Error(err)
	}
}
