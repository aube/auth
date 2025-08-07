package rest

import (
	"github.com/aube/auth/internal/api/rest/handlers_page"
	"github.com/aube/auth/internal/api/rest/middlewares"
	appPage "github.com/aube/auth/internal/application/page"

	"github.com/gin-gonic/gin"
)

func SetupPageRouter(
	api *gin.RouterGroup,
	pageService *appPage.PageService,
	jwtSecret string,
) {
	pageHandler := handlers_page.NewPageHandler(pageService, jwtSecret)

	// Защищённые маршруты
	authApi := api.Group("/")
	authApi.GET("/page", pageHandler.GetByParam)
	authApi.Use(middlewares.PaginationMiddleware())
	{
		authApi.GET("/pages", pageHandler.ListPages)
	}

	authApi.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		authApi.POST("/page", pageHandler.Create)
		authApi.PUT("/page", pageHandler.Update)
		authApi.DELETE("/page", pageHandler.Delete)
	}
}
