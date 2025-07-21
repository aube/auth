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

type UserService struct {
	repo UserRepository
	log  zerolog.Logger
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
		log:  logger.Get().With().Str("user", "service").Logger(),
	}
}

func (s *UserService) Register(ctx context.Context, userDTO dto.CreateUserDTO) (*dto.UserResponseDTO, error) {

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

func (s *UserService) Login(ctx context.Context, userDTO dto.LoginDTO) (*dto.UserResponseDTO, error) {
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

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*dto.UserResponseDTO, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("GetUserByID")
		return nil, err
	}

	return dto.NewUserResponse(user), nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}

	return nil
}
