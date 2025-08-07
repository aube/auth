// Package user provides business logic for user account operations.
package user

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/domain/valueobjects"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"
)

// UserService implements core user management functionality.
// Handles registration, authentication, and account management.
// Fields:
//   - repo: Underlying user repository
//   - log: Structured logger instance
type UserService struct {
	repo UserRepository
	log  zerolog.Logger
}

// NewUserService creates a new UserService instance.
// repo: User repository implementation
// Returns: Configured *UserService

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
		log:  logger.Get().With().Str("user", "service").Logger(),
	}
}

// Register handles new user registration:
// 1. Validates username availability
// 2. Hashes password securely
// 3. Creates user entity
// 4. Persists to repository
// 5. Returns sanitized user response
//
// ctx: Context for cancellation/timeout
// userDTO: Registration data
// Returns: (*dto.UserResponse, error)
func (s *UserService) Register(ctx context.Context, userDTO dto.RegisterRequest) (*dto.UserResponse, error) {

	// Проверяем, существует ли пользователь
	exists, err := s.repo.Exists(ctx, userDTO.Username)
	if err != nil {
		s.log.Debug().Err(err).Msg("Register1")
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}
	// Создаем value object пароля
	password, err := valueobjects.NewPassword(userDTO.Password)
	if err != nil {
		s.log.Debug().Err(err).Msg("Register2")
		return nil, err
	}

	// Хэшируем пароль
	if err := password.Hash(); err != nil {
		s.log.Debug().Err(err).Msg("Register3")
		return nil, err
	}

	// Создаем сущность пользователя
	user, err := entities.NewUser(0, userDTO.Username, userDTO.Email, password)
	if err != nil {
		s.log.Debug().Err(err).Msg("Register4")
		return nil, err
	}

	// Сохраняем в репозитории
	if err := s.repo.Create(ctx, user); err != nil {
		s.log.Debug().Err(err).Msg("Register5")
		return nil, err
	}

	return dto.NewUserResponse(user), nil
}

// Login authenticates existing users:
// 1. Verifies username exists
// 2. Validates password against stored hash
// 3. Returns user profile on success
//
// ctx: Context for cancellation/timeout
// userDTO: Login credentials
// Returns: (*dto.UserResponse, error)
func (s *UserService) Login(ctx context.Context, userDTO dto.LoginRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByUsername(ctx, userDTO.Username)
	if err != nil {
		s.log.Debug().Err(err).Msg("Login")
		return nil, err
	}

	if !user.PasswordMatches(userDTO.Password) {
		s.log.Debug().Msg("123")
		s.log.Debug().Msg(user.GetHashedPassword())
		return nil, errors.New("invalid credentials")
	}

	return dto.NewUserResponse(user), nil
}

// GetUserByID retrieves user profile:
// 1. Verifies user exists
// 2. Returns sanitized profile data
//
// ctx: Context for cancellation/timeout
// id: User identifier
// Returns: (*dto.UserResponse, error)
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("GetUserByID")
		return nil, err
	}

	return dto.NewUserResponse(user), nil
}

// Delete removes user account:
// 1. Verifies user exists
// 2. Deletes from repository
//
// ctx: Context for cancellation/timeout
// id: User identifier
// Returns: error on failure
func (s *UserService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}

	return nil
}
