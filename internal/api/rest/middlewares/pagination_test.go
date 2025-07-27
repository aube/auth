package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPaginationMiddleware_DefaultValues(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.Use(PaginationMiddleware())
	r.GET("/test", func(c *gin.Context) {
		offset := c.GetInt("offset")
		limit := c.GetInt("limit")
		c.JSON(http.StatusOK, gin.H{"offset": offset, "limit": limit})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"offset":0`)
	assert.Contains(t, w.Body.String(), `"limit":10`)
}

func TestPaginationMiddleware_CustomValues(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.Use(PaginationMiddleware())
	r.GET("/test", func(c *gin.Context) {
		offset := c.GetInt("offset")
		limit := c.GetInt("limit")
		c.JSON(http.StatusOK, gin.H{"offset": offset, "limit": limit})
	})

	req, _ := http.NewRequest("GET", "/test?page=2&size=20", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"offset":20`)
	assert.Contains(t, w.Body.String(), `"limit":20`)
}
