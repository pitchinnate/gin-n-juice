package auth

import (
	user2 "gin-n-juice/utils/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Resend struct {
	VerifyUrl string `json:"verify_url" binding:"required,url"`
}

func PostResend(c *gin.Context) {
	var json Resend
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := user2.GetUserFromContext(c)

	if err := sendVerifyEmail(*user, json.VerifyUrl); err != nil {
		log.Print("Error sending email: ", err)
	}

	c.JSON(http.StatusCreated, "Check your email to verify your email address")
}
