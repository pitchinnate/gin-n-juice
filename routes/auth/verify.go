package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tagdeploy/config"
	"tagdeploy/db"
	"tagdeploy/models"
	"time"
)

func GetVerify(c *gin.Context) {
	tokenQuery := c.Query("token")
	if tokenQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		return
	}

	tokenBytes, err := base64.URLEncoding.DecodeString(tokenQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		return
	}

	token, err := jwt.Parse(string(tokenBytes), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		hmacSampleSecret := []byte(config.ENCRYPT_KEY)
		return hmacSampleSecret, nil
	})

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "problem with token"})
		return
	}

	userId, _ := strconv.Atoi(fmt.Sprint(claims["user_id"]))
	var user models.User
	err = db.DB.First(&user, userId).Error
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&user).Update("EmailVerified", time.Now())

	newToken := models.GenerateJwt(user)

	c.JSON(http.StatusOK, gin.H{"user": user, "token": newToken})
}
