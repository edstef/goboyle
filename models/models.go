package models

import (
	"context"
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Models struct {
	db *bun.DB
}

func (m *Models) CreateTables() error {
	log.Print("Creating tables ...")
	_, err := m.db.NewCreateTable().Model(&Profile{}).IfNotExists().Exec(context.Background())
	if err != nil {
		return err
	}

	log.Print("Finisehd creating tables ...")
	return nil
}

func (m *Models) DropTables() error {
	log.Print("Dropping tables ...")
	_, err := m.db.NewDropTable().Model(&Profile{}).IfExists().Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func NewModels(dsn string, logLevel int) *Models {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	var models *Models

	if logLevel > 0 {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	models = &Models{
		db: db,
	}

	return models
}
