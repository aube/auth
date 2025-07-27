package handlers_user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aube/auth/internal/application/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, user dto.RegisterRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, credentials dto.LoginRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, credentials)
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserHandler_Register_Success(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, "test-secret")

	expectedUser := &dto.UserResponse{ID: 1, Username: "testuser"}
	mockService.On("Register", mock.Anything, mock.Anything).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/register", handler.Register)

	// Test
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(`{"username":"testuser","password":"pass12345678","email":"test@test.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_Register_Fail_Short_Password(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, "test-secret")

	// Здесь НЕ настраиваем ожидание вызова Register, так как он не должен быть вызван

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/register", handler.Register)

	// Test - пароль слишком короткий
	reqBody := `{"username":"testuser","password":"pass","email":"test@test.com"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Password")
	assert.Contains(t, w.Body.String(), "min") // Проверяем что ошибка о минимальной длине

	// Убеждаемся, что Register не был вызван
	mockService.AssertNotCalled(t, "Register")
}

func TestUserHandler_Login_Success(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, "test-secret")

	expectedUser := &dto.UserResponse{ID: 1, Username: "testuser"}
	mockService.On("Login", mock.Anything, mock.Anything).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/login", handler.Login)

	// Test
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(`{"username":"testuser","password":"pass"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
	mockService.AssertExpectations(t)
}
