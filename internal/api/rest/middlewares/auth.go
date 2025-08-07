// Package middlewares provides gin middleware.
package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens and sets the userID in the request context.
// jwtSecret: Secret key for token validation.
// Returns: Gin middleware function.
// Behavior:
//   - Extracts the "Authorization" header.
//   - Validates the JWT token and its claims.
//   - Aborts with 401 if validation fails.
//   - Sets the userID in the context for downstream handlers.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// fmt.Println("DEBUG: Authenticated user ID:", 11)
		// c.Set("userID", int64(11))
		// c.Next()
		// return

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		tokenString := authHeader[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token gone bad"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token parce error"})
			}
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		userID := claims["sub"].(float64)
		c.Set("userID", int(userID))
		fmt.Println("Authenticated user ID:", userID)
		c.Next()
	}
}
