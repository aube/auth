// Package user provides data persistence operations for user accounts.
package user

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/domain/entities"
)

// ErrUserNotFound is returned when a requested user cannot be found.
var ErrUserNotFound = errors.New("user not found")

// UserRepository defines the interface for user account persistence operations.
// Implementations should handle database operations for user records.
//
// Methods:
//
//   - Create: Stores a new user account
//     ctx: Context for cancellation/timeout
//     user: User entity to create
//     Returns: error on failure
//
//   - FindByUsername: Retrieves user by username
//     ctx: Context for cancellation/timeout
//     username: Unique username identifier
//     Returns: (*entities.User, error)
//
//   - FindByID: Retrieves user by database ID
//     ctx: Context for cancellation/timeout
//     id: Numeric user identifier
//     Returns: (*entities.User, error)
//
//   - Exists: Checks if username is already registered
//     ctx: Context for cancellation/timeout
//     username: Username to check
//     Returns: (bool, error) - true if username exists
//
//   - Delete: Removes user account by ID
//     ctx: Context for cancellation/timeout
//     id: User identifier
//     Returns: error on failure
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	Exists(ctx context.Context, username string) (bool, error)
	Delete(ctx context.Context, id int64) error
}
