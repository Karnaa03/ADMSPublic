package main

import (
	"log"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/conf"
	db2 "git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/routes"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/s3"
)

var Version string

func main() {
	config := conf.Config{
		Version: Version,
	}
	config.Load()

	db := db2.Db{}
	err := db.Init(config)
	if err != nil {
		log.Fatal(err)
	}

	s3client, err := s3.New(config)
	if err != nil {
		log.Fatal(err)
	}

	srv := routes.Server{
		Config: config,
		Db:     db,
		S3:     s3client,
	}
	err = srv.InitRouter()
	if err != nil {
		log.Fatal(err)
	}
	srv.Run()
}
