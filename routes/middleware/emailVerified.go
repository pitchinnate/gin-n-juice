package middleware

import (
	"gin-n-juice/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EmailVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "user required"})
			c.Abort()
			return
		}
		userObj := user.(models.User)
		if userObj.EmailVerified == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "user email has not been verified"})
			c.Abort()
			return
		}
		c.Next()
	}
}
