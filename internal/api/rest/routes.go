package rest

import (
	"net/http"
	"time"

	"github.com/aube/auth/internal/api/rest/handlers_common"
	"github.com/aube/auth/internal/api/rest/handlers_upload"
	"github.com/aube/auth/internal/api/rest/handlers_user"
	"github.com/aube/auth/internal/api/rest/middlewares"
	appFile "github.com/aube/auth/internal/application/file"
	appUpload "github.com/aube/auth/internal/application/upload"
	appUser "github.com/aube/auth/internal/application/user"

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

func SetupUploadsRouter(api *gin.RouterGroup, fileService *appFile.FileService, uploadService *appUpload.UploadService, jwtSecret string) {
	uploadHandler := handlers_upload.NewUploadHandler(fileService, uploadService)

	// Защищённые маршруты
	authApi := api.Group("/")
	authApi.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		authApi.GET("/file", uploadHandler.DownloadFile)
		authApi.GET("/uploads", uploadHandler.ListFiles)
		authApi.POST("/upload", uploadHandler.UploadFile)
		authApi.DELETE("/file", uploadHandler.DeleteFile)
	}
}

func SetupStaticRouter(r *gin.Engine, apiPath string) *gin.Engine {
	// Загрузка шаблонов (только index.html)
	r.LoadHTMLGlob("internal/api/rest/templates/*")

	// Статические файлы
	r.Static("/static", "internal/api/rest/static")

	webHandler := handlers_common.NewWebHandler()

	// Состояние апи
	r.GET(apiPath+"/state", webHandler.AppState200)

	// Все остальные GET запросы (кроме API) возвращают SPA
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
