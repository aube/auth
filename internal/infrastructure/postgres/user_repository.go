package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/domain/valueobjects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (username, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.GetHashedPassword()).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT id, username, encrypted_password as password FROM users WHERE username = $1`

	var (
		id       int64
		dbUser   string
		password string
	)

	err := r.db.QueryRow(ctx, query, username).Scan(&id, &dbUser, &password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	pwd, err := valueobjects.NewPassword(password)
	if err != nil {
		return nil, fmt.Errorf("invalid password format in DB: %w", err)
	}

	return entities.NewUser(id, dbUser, pwd)
}

func (r *UserRepository) Exists(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}
