// Package handlers_common provides handlers for static content and SPA serving.
package handlers_common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebHandler handles static content and single-page application (SPA) routes.
type WebHandler struct {
	// Для SPA нам не нужны зависимости сервиса в этом обработчике
	// Аутентификация будет через API
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

// ServeSPA: Serves the main SPA (index.html) for all non-API GET requests.
func (h *WebHandler) ServeSPA(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Auth App",
	})
}

// AppState200: A simple health check endpoint returning "ololo".
func (h *WebHandler) AppState200(c *gin.Context) {
	c.String(http.StatusOK, "ololo")
}
