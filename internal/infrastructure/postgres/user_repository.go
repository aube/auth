package postgres

import (
	"context"
	"errors"
	"fmt"

	appUser "github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/domain/valueobjects"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryUserInsert       string = "INSERT INTO users (username, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id"
	queryUserSelectByName string = "SELECT id, username, email, encrypted_password as password FROM users WHERE username = $1 and deleted = false"
	queryUserSelectByID   string = "SELECT id, username, email FROM users WHERE id = $1 and deleted = false"
	queryUserCheckExists  string = "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 and deleted = false)"
	queryUserDelete       string = "UPDATE users SET deleted=true WHERE id = $1"
)

type UserRepository struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:  db,
		log: logger.Get().With().Str("postgres", "user_repository").Logger(),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	err := r.db.QueryRow(ctx, queryUserInsert, user.Username, user.Email, user.GetHashedPassword()).Scan(&user.ID)
	if err != nil {
		r.log.Debug().Err(err).Msg(user.Username)
		r.log.Debug().Err(err).Msg(user.Email)
		r.log.Debug().Err(err).Msg("Create")
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var (
		id       int64
		dbUser   string
		password string
		email    string
	)

	err := r.db.QueryRow(ctx, queryUserSelectByName, username).Scan(&id, &dbUser, &email, &password)
	if err != nil {
		r.log.Debug().Err(err).Msg("FindByUsername")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appUser.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	pwd, err := valueobjects.NewPassword(password)
	if err != nil {
		return nil, fmt.Errorf("invalid password format in DB: %w", err)
	}

	return entities.NewUser(id, dbUser, email, pwd)
}

func (r *UserRepository) FindByID(ctx context.Context, userID int64) (*entities.User, error) {
	var (
		id     int64
		dbUser string
		email  string
	)

	err := r.db.QueryRow(ctx, queryUserSelectByID, userID).Scan(&id, &dbUser, &email)
	if err != nil {
		r.log.Debug().Err(err).Msg("FindByID")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appUser.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &entities.User{
		ID:       id,
		Username: dbUser,
		Email:    email,
	}, nil
}

func (r *UserRepository) Exists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, queryUserCheckExists, username).Scan(&exists)
	if err != nil {
		r.log.Debug().Err(err).Msg("Exists")
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *UserRepository) Delete(ctx context.Context, userID int64) error {

	_, err := r.db.Query(ctx, queryUserDelete, userID)

	if err != nil {
		r.log.Debug().Err(err).Msg("Delete")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
