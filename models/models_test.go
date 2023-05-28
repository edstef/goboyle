package models_test

import (
	models "github.com/edstef/goboyle/models"
	"log"
)

var mods *models.Models

func init() {
	connString := "postgres://postgres:@localhost:5432/sss?sslmode=disable"
	logLevel := 0

	mods = models.NewModels(connString, logLevel)
	err := mods.DropTables()
	if err != nil {
		log.Fatal(err)
	}

	err = mods.CreateTables()
	if err != nil {
		log.Fatal(err)
	}
}
