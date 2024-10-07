package models_test

import (
	models "github.com/edstef/goboyle/models"
)

var mods *models.Models

func init() {
	connString := "postgres://postgres:@localhost:5432/sss?sslmode=disable"
	logLevel := 0

	mods = models.NewModels(connString, logLevel)
	mods.DropTables()
	mods.CreateTables()
}
