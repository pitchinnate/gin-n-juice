package users

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetList(c *gin.Context) {
	var users []models.User

	if err := db.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
