package postgres

import (
	"embed"

	"github.com/aube/auth/internal/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Postgres
//
//go:embed migrations/*.sql
var embedPostgresMigrations embed.FS

func runPostgresMigrations(pool *pgxpool.Pool) error {
	log := logger.Get().With().Str("migrations", "runPostgresMigrations").Logger()

	// Run migrations
	stdDB := stdlib.OpenDBFromPool(pool)
	defer stdDB.Close()

	goose.SetBaseFS(embedPostgresMigrations)

	log.Debug().Msg("goose.SetDialect postgres")
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// log.Debug().Msg("goose.Down migrations")
	// if err := goose.DownTo(stdDB, "migrations", 0); err != nil {
	// 	return err
	// }

	log.Debug().Msg("goose.Up migrations")
	if err := goose.Up(stdDB, "migrations"); err != nil {
		return err
	}
	return nil
}
