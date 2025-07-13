package rest

import (
	"net/http"
	"time"

	appFile "github.com/aube/auth/internal/application/file"
	appUser "github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/interfaces/rest/handlers"
	"github.com/aube/auth/internal/interfaces/rest/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(apiPath string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return r, r.Group(apiPath)
}

func SetupUserRouter(api *gin.RouterGroup, userService *appUser.UserService, jwtSecret string) {
	apiHandler := handlers.NewUserHandler(userService, jwtSecret)

	// API маршруты
	api.POST("/register", apiHandler.Register)
	api.POST("/login", apiHandler.Login)

	// Защищённые маршруты
	authApi := api.Group("/")
	authApi.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		authApi.GET("/profile", apiHandler.GetProfile)
		authApi.POST("/logout", apiHandler.Logout)
	}
}

func SetupFilesRouter(api *gin.RouterGroup, fileService *appFile.FileService, jwtSecret string) {
	fileHandler := handlers.NewFileHandler(fileService)

	// Защищённые маршруты
	authApi := api.Group("/")
	authApi.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		authApi.POST("/upload", fileHandler.UploadFile)
		authApi.GET("/download", fileHandler.DownloadFile)
		authApi.GET("/files", fileHandler.ListFiles)
		authApi.DELETE("/delete", fileHandler.DeleteFile)
	}
}

func SetupStaticRouter(r *gin.Engine) *gin.Engine {
	// Загрузка шаблонов (только index.html)
	r.LoadHTMLGlob("internal/interfaces/rest/templates/*")

	// Статические файлы
	r.Static("/static", "internal/interfaces/rest/static")

	webHandler := handlers.NewWebHandler()

	// Все GET запросы (кроме API) возвращают SPA
	r.GET("/", webHandler.ServeSPA)
	r.GET("/login", webHandler.ServeSPA)
	r.GET("/register", webHandler.ServeSPA)
	r.GET("/profile", webHandler.ServeSPA)

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
