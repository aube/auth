package handlers

import (
	"net/http"
	"time"

	"github.com/aube/auth/internal/api/rest/dto"
	appUser "github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	userService *appUser.UserService
	jwtSecret   []byte
	log         zerolog.Logger
}

func NewUserHandler(userService *appUser.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtSecret:   []byte(jwtSecret),
		log:         logger.Get().With().Str("handlers", "user_handler").Logger(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug().Err(err).Msg("Register1")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	userDTO := appUser.CreateUserDTO{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	createdUser, err := h.userService.Register(ctx, userDTO)
	if err != nil {
		h.log.Debug().Err(err).Msg("Register2")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *UserHandler) Logout(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})

}

func (h *UserHandler) Login(c *gin.Context) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug().Err(err).Msg("Login1")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	loginDTO := appUser.LoginDTO{
		Username: req.Username,
		Password: req.Password,
	}

	userEntity, err := h.userService.Login(ctx, loginDTO)
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

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userService.GetUserByID(ctx, userID.(int64))
	if err != nil {
		h.log.Debug().Err(err).Msg("GetProfile")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
