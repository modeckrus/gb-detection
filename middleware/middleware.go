package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wazyiz/jwt-gin/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unathorized")
			c.Abort()
			return
		}

		c.Next()
	}
}
