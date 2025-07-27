package handlers_common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWebHandler_ServeSPA(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.LoadHTMLGlob("../templates/*")

	handler := NewWebHandler()
	r.GET("/", handler.ServeSPA)

	// Test
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
}

func TestWebHandler_AppState200(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	handler := NewWebHandler()
	r.GET("/state", handler.AppState200)

	// Test
	req, _ := http.NewRequest("GET", "/state", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ololo", w.Body.String())
}
