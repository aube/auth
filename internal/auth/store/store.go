package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/aube/gophermart/internal/auth/store/postgres"
	"github.com/aube/gophermart/internal/auth/store/repos"
)

type Store interface {
	User() repos.UserRepository
}

func NewStore(config string) (Store, error) {
	db, err := NewDB(config)

	pgstore := postgres.New(db)

	return pgstore, err
}

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("PostgreSQL database connection established")

	runPostgresMigrations(db)

	return db, nil
}
