package models

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

const DefaultSender = "support@toucan.com"

// SMTPConfig load field from env
type SMTPConfig struct {
	HOST     string
	PORT     int
	Username string
	Password string
}

type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}

type EmailService struct {
	dialer        *gomail.Dialer
	DefaultSender string
}

func NewEmailService() *EmailService {
	var config SMTPConfig
	config.HOST = os.Getenv("SMTP_HOST")
	config.PORT, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	config.Username = os.Getenv("SMTP_USERNAME")
	config.Password = os.Getenv("SMTP_PASSWORD")
	return &EmailService{dialer: gomail.NewDialer(config.HOST, config.PORT, config.Username, config.Password)}
}

func (es *EmailService) Send(email Email) error {
	message := gomail.NewMessage()
	message.SetHeader("To", email.To)
	es.setFrom(message, email)
	message.SetHeader("Subject", email.Subject)
	switch {
	case email.PlainText != "" && email.HTML != "":
		message.SetBody("text/plain", email.PlainText)
		message.SetBody("text/html", email.HTML)
	case email.PlainText != "":
		message.SetBody("text/plain", email.PlainText)
	case email.HTML != "":
		message.SetBody("text/html", email.HTML)
	}
	err := es.dialer.DialAndSend(message)
	if err != nil {
		return fmt.Errorf("send email error %w", err)
	}
	return nil
}

func (es *EmailService) setFrom(message *gomail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	message.SetHeader("From", from)
}

func (es EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		From:      "support@toucan",
		To:        to,
		Subject:   "Reset Your Password",
		PlainText: fmt.Sprintf("To set your password, please visit the following link: %s", resetURL),
		HTML:      fmt.Sprintf(`To set your password, please visit the following link: <a href="%s">%s</a>`, resetURL, resetURL),
	}
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password reset %w", err)
	}
	return nil
}
