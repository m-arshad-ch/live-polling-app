package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"live-polling-app/utils"
)

func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "login required"})
			c.Abort()
			return
		}

		userID, err := utils.ParseToken(parts[1], secret)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
