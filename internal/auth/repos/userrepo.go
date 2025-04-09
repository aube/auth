package repos

import (
	"context"

	"github.com/aube/gophermart/internal/auth/model"
)

type IUserRepository interface {
	Create(context.Context, *model.User) error
	Find(context.Context, int) (*model.User, error)
	FindByEmail(context.Context, string) (*model.User, error)
}
