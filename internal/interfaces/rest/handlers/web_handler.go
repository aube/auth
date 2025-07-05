package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebHandler struct {
	// Для SPA нам не нужны зависимости сервиса в этом обработчике
	// Аутентификация будет через API
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

// ServeSPA обрабатывает все HTML-запросы и возвращает index.html
func (h *WebHandler) ServeSPA(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Auth App",
	})
}
