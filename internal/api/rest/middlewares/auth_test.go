package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	jwtSecret := "test-secret"
	r.Use(AuthMiddleware(jwtSecret))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "authorization header required")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	jwtSecret := "test-secret"
	r.Use(AuthMiddleware(jwtSecret))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}
