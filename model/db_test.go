package model

import (
	"testing"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/conf"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDb_GetAgregate(t *testing.T) {
	config := conf.Config{
		Version: "test",
		DbLog:   true,
	}
	config.Load()

	db := Db{}
	err := db.Init(config)
	if err != nil {
		log.Fatal(err)
	}

	data, err := db.GetAgregate("46", "65", "", "", "", "1")
	t.Logf("returned data : %v", data)
	assert.Nil(t, err)
}

func TestGetGeoRequest(t *testing.T) {
	type args struct {
		division string
		district string
		upazilla string
		union    string
		mouza    string
	}
	tests := []struct {
		name         string
		args         args
		wantSelector string
	}{
		{
			name: "just division",
			args: args{
				division: "20",
			},
			wantSelector: "20",
		},
		{
			name: "division with district",
			args: args{
				division: "20",
				district: "65",
			},
			wantSelector: "20.65",
		},
		{
			name: "division with district with upazilla",
			args: args{
				division: "20",
				district: "65",
				upazilla: "34",
			},
			wantSelector: "20.65.34",
		},
		{
			name: "with all",
			args: args{
				division: "20",
				district: "65",
				upazilla: "34",
				union:    "45",
				mouza:    "23",
			},
			wantSelector: "20.65.34.45.23",
		},
		{
			name: "with hole",
			args: args{
				division: "20",
				district: "65",
				union:    "23",
			},
			wantSelector: "20.65",
		},
		{
			name:         "with nothing",
			args:         args{},
			wantSelector: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSelector := GetGeoRequest(tt.args.division, tt.args.district, tt.args.upazilla, tt.args.union, tt.args.mouza); gotSelector != tt.wantSelector {
				t.Errorf("GetGeoRequest() = %v, want %v", gotSelector, tt.wantSelector)
			}
		})
	}
}
