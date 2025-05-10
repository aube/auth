package memory

import (
	"context"

	"github.com/aube/auth/internal/model"
)

// ActiveUserProvider ...
type ActiveUserProvider interface {
	Set(context.Context, *model.User) error
	Get(context.Context, string) (*model.User, bool)
}

// Store ...
type MemoryStore struct {
	mem map[string]*model.User
	oq  []int
}

// ActiveUser ...
func (s *MemoryStore) ActiveUser() ActiveUserProvider {
	return &ActiveUserRepository{
		mem: s.mem,
	}
}

// New ...
func New() *MemoryStore {
	return &MemoryStore{
		mem: make(map[string]*model.User),
		oq:  []int{},
	}
}
