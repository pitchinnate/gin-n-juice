package users

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"gin-n-juice/utils/input"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Update(c *gin.Context) {
	var json_user models.User
	if err := c.ShouldBindJSON(&json_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := input.GetPathInt(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Model(&user).Select("*").Omit("id", "created_at", "password").Updates(json_user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
