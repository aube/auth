package rest

import (
	"net/http"

	"github.com/aube/auth/internal/api/rest/handlers_common"

	"github.com/gin-gonic/gin"
)

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
