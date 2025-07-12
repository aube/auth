package postgres

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Postgres
//
//go:embed migrations/*.sql
var embedPostgresMigrations embed.FS

func runPostgresMigrations(pool *pgxpool.Pool) error {
	// Run migrations
	stdDB := stdlib.OpenDBFromPool(pool)
	defer stdDB.Close()

	goose.SetBaseFS(embedPostgresMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Down(stdDB, "migrations"); err != nil {
		return err
	}

	if err := goose.Up(stdDB, "migrations"); err != nil {
		return err
	}
	return nil
}
