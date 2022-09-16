package users

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Create(c *gin.Context) {
	var json_user models.User
	if err := c.ShouldBindJSON(&json_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&json_user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": json_user})
}
