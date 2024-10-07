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
	m.db.NewCreateTable().Model(&User{}).
		IfNotExists().Exec(context.Background())

	m.db.NewCreateTable().Model(&UserToken{}).
		ForeignKey(`("user") REFERENCES "users" ("id") ON DELETE CASCADE`).
		IfNotExists().Exec(context.Background())

	log.Print("Finisehd creating tables ...")
	return nil
}

func (m *Models) DropTables() error {
	log.Print("Dropping tables ...")
	m.db.NewDropTable().Model(&UserToken{}).IfExists().Exec(context.Background())
	m.db.NewDropTable().Model(&User{}).IfExists().Exec(context.Background())

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
