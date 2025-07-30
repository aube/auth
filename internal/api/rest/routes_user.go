package rest

import (
	"github.com/aube/auth/internal/api/rest/handlers_user"
	"github.com/aube/auth/internal/api/rest/middlewares"
	appUser "github.com/aube/auth/internal/application/user"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(api *gin.RouterGroup, userService *appUser.UserService, jwtSecret string) {
	userHandler := handlers_user.NewUserHandler(userService, jwtSecret)

	// API маршруты
	api.POST("/register", userHandler.Register)
	api.POST("/login", userHandler.Login)

	// Защищённые маршруты
	authApi := api.Group("/")
	authApi.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		authApi.GET("/profile", userHandler.GetProfile)
		authApi.POST("/logout", userHandler.Logout)
		authApi.DELETE("/profile", userHandler.Delete)
	}
}
