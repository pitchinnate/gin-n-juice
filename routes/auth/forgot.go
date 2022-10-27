package auth

import (
	"fmt"
	"gin-n-juice/db"
	"gin-n-juice/models"
	"gin-n-juice/utils/mail"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type Forgot struct {
	Email     string `json:"email" binding:"required,email"`
	ReturnUrl string `json:"return_url" binding:"required,url"`
}

func PostForgot(c *gin.Context) {
	var json Forgot
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	if err := db.DB.Where("UPPER(email) like UPPER(?)", json.Email).First(&user).Error; err != nil {
		// report that an email was sent even if it wasn't
		log.Printf("forgot password for: %s but no email found", json.Email)
		c.JSON(http.StatusOK, "email sent")
		return
	}

	if err := sendGoodEmail(user, json.ReturnUrl); err != nil {
		log.Printf("Error sending email: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "email sent")
}

func sendGoodEmail(user models.User, returnUrl string) error {
	var message = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Forgot Password</title>
	</head>
	<body>
		A forgot password request was made for %s. Click the following link or copy and paste it.<br><br>
		<a href='%s'>%s</a>
<br><br>
If you did not request a forgot password you can ignore this email.
	</body>
</html>`

	b64token := uuid.NewString()
	url := fmt.Sprintf("%s?token=%s&email=%s", returnUrl, b64token, user.Email)

	db.DB.Delete(models.PasswordReset{}, "UPPER(email) like UPPER(?)", user.Email)
	forgot := models.PasswordReset{
		Email:     user.Email,
		Token:     b64token,
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&forgot).Error; err != nil {
		return err
	}

	newMessage := mail.EmailMessage{
		user.Email,
		nil,
		"Forgot Password Request",
		fmt.Sprintf(message, user.Email, url, url),
	}
	return mail.SendEmail(newMessage)
}

func sendBadEmail(email string) error {
	var message = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Forgot Password</title>
	</head>
	<body>
		A forgot password request was made for %s however that email address does not have account with us. If you
		requested the forgot password please try a different email address. If you did not request a forgot password
		you can ignore this email.
	</body>
</html>`

	newMessage := mail.EmailMessage{
		email,
		nil,
		"Forgot Password Request",
		fmt.Sprintf(message, email),
	}
	return mail.SendEmail(newMessage)
}
