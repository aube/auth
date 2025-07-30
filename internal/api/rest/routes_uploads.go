package rest

import (
	"github.com/aube/auth/internal/api/rest/handlers_upload"
	"github.com/aube/auth/internal/api/rest/middlewares"
	appFile "github.com/aube/auth/internal/application/file"
	appUpload "github.com/aube/auth/internal/application/upload"

	"github.com/gin-gonic/gin"
)

func SetupUploadsRouter(api *gin.RouterGroup, fileService *appFile.FileService, uploadService *appUpload.UploadService, jwtSecret string) {
	uploadHandler := handlers_upload.NewUploadHandler(fileService, uploadService)

	// Защищённые маршруты
	authApi := api.Group("/")
	authApi.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		authApi.GET("/upload", uploadHandler.DownloadFile)
		authApi.POST("/upload", uploadHandler.UploadFile)
		authApi.DELETE("/upload", uploadHandler.DeleteFile)
	}
	authApi.Use(middlewares.PaginationMiddleware())
	{
		authApi.GET("/uploads", uploadHandler.ListFiles)
	}
}
