// Package dto contains data transfer objects for user operations.
package dto

import "github.com/aube/auth/internal/domain/entities"

// RegisterRequest represents user registration input.
// Fields:
//   - Username: Required, 3-50 characters.
//   - Email: Required, valid email format.
//   - Password: Required, minimum 8 characters.
//
// Validation tags are used for Gin request binding.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents user authentication input.
// Fields:
//   - Username: Required, 3-50 characters (alternative to Email).
//   - Password: Required.
//   - Email: Optional alternative to Username.
//
// At least one of Username or Email must be provided.
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

// UserResponse represents user profile data in API responses.
// Fields:
//   - ID: User's unique identifier.
//   - Username: User's display name.
//   - Email: User's contact email.
type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// NewUserResponse creates a UserResponse from an entities.User.
// user: Source user entity.
// Returns: Populated UserResponse DTO.
func NewUserResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
