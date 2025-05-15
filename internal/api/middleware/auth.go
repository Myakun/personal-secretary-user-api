package middleware

import (
	"errors"
	"net/http"
	"strings"

	jwtPkg "github.com/Myakun/personal-secretary-user-api/internal/common/jwt"

	"github.com/gin-gonic/gin"
)

// UserContext keys
const (
	UserIDKey = "user_id"
	EmailKey  = "email"
)

// AuthMiddleware is a middleware that validates the Bearer token in the Authorization header
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Validate the token
		claims, err := jwtPkg.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			var statusCode int
			var errorMessage string

			switch {
			case errors.Is(err, jwtPkg.ErrExpiredToken):
				statusCode = http.StatusUnauthorized
				errorMessage = "token has expired"
			default:
				statusCode = http.StatusUnauthorized
				errorMessage = "invalid token"
			}

			c.AbortWithStatusJSON(statusCode, gin.H{"error": errorMessage})
			return
		}

		// Set user information in the context
		c.Set(UserIDKey, claims.UserId)
		c.Set(EmailKey, claims.Email)

		c.Next()
	}
}
