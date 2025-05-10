package memory

import (
	"context"

	"github.com/aube/auth/internal/model"
)

// ActiveUserRepository ...
type ActiveUserRepository struct {
	mem map[string]*model.User
}

// Set ...
func (r *ActiveUserRepository) Set(ctx context.Context, u *model.User) error {
	r.mem[u.RandomHash] = u
	return nil
}

// Get ...
func (r *ActiveUserRepository) Get(ctx context.Context, RandomHash string) (*model.User, bool) {
	user, ok := r.mem[RandomHash]

	return user, ok
}
