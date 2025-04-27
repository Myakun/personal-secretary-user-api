package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"personal-secretary-user-ap/internal/entity/accesstoken"
	"strings"
)

// AuthMiddleware is a middleware that validates the Bearer token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if "" == authHeader {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if "" == token {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessTokenService := accesstoken.GetAccessTokenService()
		accessToken, err := accessTokenService.FindOneByToken(token)
		if nil != err {
			// TODO: Log CRITICAL
			// TODO: Idea about monitoring: send stat to grafana? And grafana send alert its critical part
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if nil == accessToken {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("accessToken", accessToken)

		c.Next()
	}
}
