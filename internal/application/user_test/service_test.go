package user_test

import (
	"context"
	"testing"

	"github.com/aube/auth/internal/application/dto"
	appUser "github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_Register_Success(t *testing.T) {
	// Setup
	mockRepo := new(UserRepository)
	service := appUser.NewUserService(mockRepo)

	// Test data
	registerReq := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockRepo.On("Exists", mock.Anything, "testuser").Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.User")).
		Run(func(args mock.Arguments) {
			userArg := args.Get(1).(*entities.User)
			assert.Equal(t, "testuser", userArg.Username)
			assert.Equal(t, "test@example.com", userArg.Email)
			assert.NotEmpty(t, userArg.GetHashedPassword())
		}).
		Return(nil)

	// Execute
	response, err := service.Register(context.Background(), registerReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "test@example.com", response.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_UserExists(t *testing.T) {
	// Setup
	mockRepo := new(UserRepository)
	service := appUser.NewUserService(mockRepo)

	// Mock expectations
	mockRepo.On("Exists", mock.Anything, "existinguser").Return(true, nil)

	// Execute
	_, err := service.Register(context.Background(), dto.RegisterRequest{
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "password123",
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Login_Success(t *testing.T) {
	// Setup
	mockRepo := new(UserRepository)
	service := appUser.NewUserService(mockRepo)

	// Test data
	hashedPassword, _ := valueobjects.NewPassword("password123")
	_ = hashedPassword.Hash()
	testUser := &entities.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	// Mock expectations
	mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(testUser, nil)

	// Execute
	response, err := service.Login(context.Background(), dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "testuser", response.Username)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	// Setup
	mockRepo := new(UserRepository)
	service := appUser.NewUserService(mockRepo)

	// Test data
	hashedPassword, _ := valueobjects.NewPassword("correctpassword")
	_ = hashedPassword.Hash()
	testUser := &entities.User{
		Password: hashedPassword,
	}

	// Mock expectations
	mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(testUser, nil)

	// Execute
	_, err := service.Login(context.Background(), dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	// Setup
	mockRepo := new(UserRepository)
	service := appUser.NewUserService(mockRepo)

	// Test data
	testUser := &entities.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Mock expectations
	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(testUser, nil)

	// Execute
	response, err := service.GetUserByID(context.Background(), 1)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "testuser", response.Username)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_Success(t *testing.T) {
	// Setup
	mockRepo := new(UserRepository)
	service := appUser.NewUserService(mockRepo)

	// Mock expectations
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	// Execute
	err := service.Delete(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
