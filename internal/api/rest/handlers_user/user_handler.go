package handlers_user

import (
	"context"
	"net/http"
	"time"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	Delete(ctx context.Context, id int64) error
	GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	Login(ctx context.Context, userDTO dto.LoginRequest) (*dto.UserResponse, error)
	Register(ctx context.Context, userDTO dto.RegisterRequest) (*dto.UserResponse, error)
}

type UserHandler interface {
	Delete(c *gin.Context)
	GetProfile(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Register(c *gin.Context)
}

type Handler struct {
	userService UserService
	jwtSecret   []byte
	log         zerolog.Logger
}

func NewUserHandler(userService UserService, jwtSecret string) UserHandler {
	return &Handler{
		userService: userService,
		jwtSecret:   []byte(jwtSecret),
		log:         logger.Get().With().Str("handlers", "user_handler").Logger(),
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug().Err(err).Msg("Register1")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	userDTO := dto.RegisterRequest(req)

	createdUser, err := h.userService.Register(ctx, userDTO)
	if err != nil {
		h.log.Debug().Err(err).Msg("Register2")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) Logout(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})

}

func (h *Handler) Login(c *gin.Context) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug().Err(err).Msg("Login1")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	LoginRequest := dto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	userEntity, err := h.userService.Login(ctx, LoginRequest)
	if err != nil {
		h.log.Debug().Err(err).Msg("Login2")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userEntity.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		h.log.Debug().Err(err).Msg("Login3")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func (h *Handler) GetProfile(c *gin.Context) {
	uID, exists := c.Get("userID")
	if !exists {
		h.log.Debug().Msg("GetProfile not exists")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := uID.(int64)
	if !ok {
		h.log.Debug().Msg("GetProfile not ok")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized2"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		h.log.Debug().Err(err).Msg("GetProfile")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	h.log.Debug().Msg(user.Username)

	c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c *gin.Context) {
	uID, exists := c.Get("userID")
	if !exists {
		h.log.Debug().Msg("GetProfile not exists")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := uID.(int64)
	if !ok {
		h.log.Debug().Msg("GetProfile not ok")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized2"})
		return
	}

	ctx := c.Request.Context()
	err := h.userService.Delete(ctx, userID)
	if err != nil {
		h.log.Debug().Err(err).Msg("Delete")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
