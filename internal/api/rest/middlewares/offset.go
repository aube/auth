package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func OffsetMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := 1
		size := 5
		maxSize := 10 // Maximum allowed page size

		if pStr := c.Query("page"); pStr != "" {
			if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
				page = p
			}
		}

		if sStr := c.Query("size"); sStr != "" {
			if s, err := strconv.Atoi(sStr); err == nil && s > 0 {
				if s > maxSize {
					size = maxSize
				} else {
					size = s
				}
			}
		}

		offset := (page - 1) * size
		c.Set("offset", offset)
		c.Set("limit", size)
		c.Next()
	}
}
