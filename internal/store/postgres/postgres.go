package postgres

import (
	"context"
	"database/sql"

	"github.com/aube/auth/internal/model"
)

// BillingProvider ...
type BillingProvider interface {
	Withdrawals(context.Context, int) ([]model.Withdraw, error)
	BalanceWithdraw(context.Context, *model.Withdraw, *model.User) error
}

// OrderProvider ...
type OrderProvider interface {
	Orders(context.Context, int) ([]model.Order, error)
	UploadOrders(context.Context, *model.Order) error
	GetNewOrdersID() ([]int, error)
	SetStatus(int, string) error
	SetAccrual(int, int) error
}

// UserProvider ...
type UserProvider interface {
	Register(context.Context, *model.User) error
	Login(context.Context, *model.User) (*model.User, error)
	Balance(context.Context, *model.User) error
}

// Store ...
type SQLStore struct {
	db *sql.DB
}

// User ...
func (s *SQLStore) User() UserProvider {
	return &UserRepository{
		db: s.db,
	}
}

// Order ...
func (s *SQLStore) Order() OrderProvider {
	return &OrderRepository{
		db: s.db,
	}
}

// Billing ...
func (s *SQLStore) Billing() BillingProvider {
	return &BillingRepository{
		db: s.db,
	}
}

// New ...
func New(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}
