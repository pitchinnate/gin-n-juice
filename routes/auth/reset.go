package auth

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Reset struct {
	Email           string `json:"email" binding:"required,email"`
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

func PostReset(c *gin.Context) {
	var json Reset
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	if err := db.DB.Where("UPPER(email) like UPPER(?)", json.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error updating user"})
		return
	}

	reset := models.PasswordReset{}
	err := db.DB.Where("UPPER(email) like UPPER(?)", json.Email).First(&reset).Error
	if err != nil || reset.Token != json.Token {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	password, err := models.HashPassword(json.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Password = password
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Delete(models.PasswordReset{}, "UPPER(email) like UPPER(?)", user.Email)

	c.JSON(http.StatusOK, "password updated")
}
