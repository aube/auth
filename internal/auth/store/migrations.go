package store

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed postgres/migrations/*.sql
var embedPostgresMigrations embed.FS

func runPostgresMigrations(db *sql.DB) {

	goose.SetBaseFS(embedPostgresMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "postgres/migrations"); err != nil {
		panic(err)
	}

}
