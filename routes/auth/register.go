package auth

import (
	"encoding/base64"
	"fmt"
	"gin-n-juice/db"
	"gin-n-juice/models"
	"gin-n-juice/utils/mail"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Register struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	VerifyUrl       string `json:"verify_url" binding:"required,url"`
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

	if err := sendVerifyEmail(user, json.VerifyUrl); err != nil {
		log.Print("Error sending email: ", err)
	}

	c.JSON(http.StatusCreated, "Check your email to verify your email address")
}

func sendVerifyEmail(user models.User, returnUrl string) error {
	var message = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Verify Email Address</title>
	</head>
	<body>
		Click the following link or copy and paste it to verify your email address.<br><br>
		<a href='%s'>%s</a>
<br><br>
If you did not create an account, no further action is required.
	</body>
</html>`

	token := models.GenerateJwt(user)
	b64token := base64.URLEncoding.EncodeToString([]byte(token))
	url := fmt.Sprintf("%s?token=%s", returnUrl, b64token)

	newMessage := mail.EmailMessage{
		user.Email,
		nil,
		"Verify Email Address",
		fmt.Sprintf(message, url, url),
	}
	return mail.SendEmail(newMessage)
}
