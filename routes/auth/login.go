package auth

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func PostLogin(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	if err := db.DB.Where("UPPER(email) like UPPER(?)", json.Email).First(&user).Error; err != nil {
		// do a hash anyway so people can check timing of response for email
		models.CheckPasswordHash("testing", "testing")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	if !models.CheckPasswordHash(json.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
