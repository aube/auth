package dto_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRequest_Fields(t *testing.T) {
	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Equal(t, "testuser", req.Username)
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestLoginRequest_Fields(t *testing.T) {
	req := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	assert.Equal(t, "testuser", req.Username)
	assert.Equal(t, "password123", req.Password)
	assert.Equal(t, "test@example.com", req.Email)
}

func TestNewUserResponse(t *testing.T) {
	user := &entities.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	resp := dto.NewUserResponse(user)

	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
}

func TestUserResponse_Empty(t *testing.T) {
	resp := dto.UserResponse{}

	assert.Equal(t, int64(0), resp.ID)
	assert.Empty(t, resp.Username)
	assert.Empty(t, resp.Email)
}

func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.RegisterRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: dto.RegisterRequest{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "password123",
			},
			valid: true,
		},
		{
			name: "short username",
			request: dto.RegisterRequest{
				Username: "ab",
				Email:    "valid@example.com",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "wrong email",
			request: dto.RegisterRequest{
				Username: "ab",
				Email:    "2`34`c235v134v5",
				Password: "password123",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый контекст Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Конвертируем запрос в JSON
			jsonData, err := json.Marshal(tt.request)
			require.NoError(t, err)

			// Создаем запрос
			c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// Пробуем привязать
			var req dto.RegisterRequest
			err = c.ShouldBindJSON(&req)

			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
