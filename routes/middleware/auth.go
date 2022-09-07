package middleware

import (
	"fmt"
	"gin-n-juice/config"
	"gin-n-juice/db"
	"gin-n-juice/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("token")
		if tokenString == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "token header required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			hmacSampleSecret := []byte(config.ENCRYPT_KEY)
			return hmacSampleSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusForbidden, gin.H{"error": "problem with token"})
			c.Abort()
			return
		}

		userId, _ := strconv.Atoi(fmt.Sprint(claims["user_id"]))
		var user models.User
		err = db.DB.First(&user, userId).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
