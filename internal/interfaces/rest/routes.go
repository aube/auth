package rest

import (
	"net/http"

	"github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/interfaces/rest/handlers"
	"github.com/aube/auth/internal/interfaces/rest/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userService *user.UserService, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// Загрузка шаблонов (только index.html)
	r.LoadHTMLGlob("interfaces/rest/templates/*")

	// Статические файлы
	r.Static("/static", "./interfaces/rest/static")

	webHandler := handlers.NewWebHandler()
	apiHandler := handlers.NewUserHandler(userService, jwtSecret)

	// Все GET запросы (кроме API) возвращают SPA
	r.GET("/", webHandler.ServeSPA)
	r.GET("/login", webHandler.ServeSPA)
	r.GET("/register", webHandler.ServeSPA)
	r.GET("/profile", webHandler.ServeSPA)

	// API маршруты
	api := r.Group("/api/v1")
	{
		api.POST("/register", apiHandler.Register)
		api.POST("/login", apiHandler.Login)
		api.POST("/logout", apiHandler.Logout)

		// Защищённые маршруты
		authApi := api.Group("/")
		authApi.Use(middlewares.AuthMiddleware(jwtSecret))
		{
			authApi.GET("/profile", apiHandler.GetProfile)
		}
	}

	// Обработка 404 для API
	r.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		// Для всех остальных запросов возвращаем SPA
		webHandler.ServeSPA(c)
	})

	return r
}
