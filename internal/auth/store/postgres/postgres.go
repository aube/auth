package postgres

import (
	"database/sql"

	"github.com/aube/gophermart/internal/auth/store/repos"
)

// Store ...
type SQLStore struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New ...
func New(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

// User ...
func (s *SQLStore) User() repos.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
