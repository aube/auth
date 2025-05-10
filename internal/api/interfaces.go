package api

import (
	"context"

	"github.com/aube/auth/internal/model"
)

type UserProvider interface {
	Register(context.Context, *model.User) error
	Login(context.Context, *model.User) (*model.User, error)
	Balance(context.Context, *model.User) error
}

type ActiveUserProvider interface {
	Set(context.Context, *model.User) error
	Get(context.Context, string) (*model.User, bool)
}
