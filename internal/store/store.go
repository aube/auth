package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/aube/auth/internal/model"
	"github.com/aube/auth/internal/store/memory"
	"github.com/aube/auth/internal/store/postgres"
)

// ActiveUserProvider ...
type ActiveUserProvider interface {
	Set(context.Context, *model.User) error
	Get(context.Context, string) (*model.User, bool)
}

// UserProvider ...
type UserProvider interface {
	Register(context.Context, *model.User) error
	Login(context.Context, *model.User) (*model.User, error)
	Balance(context.Context, *model.User) error
}

type Store struct {
	User       UserProvider
	ActiveUser ActiveUserProvider
}

func NewStore(config string) (Store, error) {
	db, err := NewDB(config)
	if err != nil {
		return Store{}, err
	}

	store := Store{
		User:       postgres.New(db).User(),
		ActiveUser: memory.New().ActiveUser(),
	}

	return store, nil
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
