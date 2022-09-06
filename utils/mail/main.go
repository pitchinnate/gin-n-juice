package mail

import (
	"gin-n-juice/config"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type EmailMessage struct {
	To      string
	From    string
	Subject string
	Body    string
}

func setupMailer() mailgun.Mailgun {
	return mailgun.NewMailgun(config.MAILGUN_DOMAIN,
		config.MAILGUN_PRIVATE_KEY,
		config.MAILGUN_VALIDATION_KEY)
}

func SendEmail(message EmailMessage) error {
	if config.IS_TESTING {
		return nil
	}
	mg := setupMailer()
	newMessage := mg.NewMessage(message.From, message.Subject, "Please view in HTML", message.To)
	newMessage.SetHtml(message.Body)
	_, _, err := mg.Send(newMessage)
	return err
}
