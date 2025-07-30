package rest

import (
	"github.com/aube/auth/internal/api/rest/handlers_image"
	"github.com/aube/auth/internal/api/rest/middlewares"
	appFile "github.com/aube/auth/internal/application/file"
	appImage "github.com/aube/auth/internal/application/image"

	"github.com/gin-gonic/gin"
)

func SetupImagesRouter(api *gin.RouterGroup, fileService *appFile.FileService, imageService *appImage.ImageService, jwtSecret string) {
	imageHandler := handlers_image.NewImageHandler(fileService, imageService)

	// Защищённые маршруты
	authApi := api.Group("/")
	authApi.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		authApi.GET("/image", imageHandler.DownloadFile)
		authApi.POST("/image", imageHandler.UploadImage)
		authApi.DELETE("/image", imageHandler.DeleteFile)
	}
	authApi.Use(middlewares.PaginationMiddleware())
	{
		authApi.GET("/images", imageHandler.ListFiles)
	}
}
