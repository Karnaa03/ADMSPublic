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

	_, err = db.GetAgregate("20", "46", "65", "", "", "1")

	assert.Nil(t, err)
}
