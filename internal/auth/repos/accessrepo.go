package repos

import (
	"github.com/aube/gophermart/internal/auth/model"
)

type AccessRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
