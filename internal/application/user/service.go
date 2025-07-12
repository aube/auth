package user

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/domain/valueobjects"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, dto CreateUserDTO) (*UserResponseDTO, error) {
	// Проверяем, существует ли пользователь
	exists, err := s.repo.Exists(ctx, dto.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Создаем value object пароля
	password, err := valueobjects.NewPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	// Хэшируем пароль
	if err := password.Hash(); err != nil {
		return nil, err
	}

	// Создаем сущность пользователя
	user, err := entities.NewUser(0, dto.Username, password)
	if err != nil {
		return nil, err
	}

	// Сохраняем в репозитории
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &UserResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserService) Login(ctx context.Context, dto LoginDTO) (*entities.User, error) {
	user, err := s.repo.FindByUsername(ctx, dto.Username)
	if err != nil {
		return nil, err
	}

	if !user.PasswordMatches(dto.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*UserResponseDTO, error) {
	// Реализация будет зависеть от вашего репозитория
	// Это примерный интерфейс
	return nil, errors.New("not implemented")
}
