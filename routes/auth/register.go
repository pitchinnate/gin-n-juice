package auth

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Register struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

func PostRegister(c *gin.Context) {
	var json Register
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	result := db.DB.Where("UPPER(email) like UPPER(?)", json.Email).First(&user)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email address already used, use forgot password"})
		return
	}

	user.Password = json.Password
	user.Email = json.Email
	user.Admin = false

	result = db.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	token := models.GenerateJwt(user)

	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}
