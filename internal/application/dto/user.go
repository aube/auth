package dto

import "github.com/aube/auth/internal/domain/entities"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
