package user

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/domain/entities"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	Exists(ctx context.Context, username string) (bool, error)
}
