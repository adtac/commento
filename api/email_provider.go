package main

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailProvider is the interface for sending transactional emails.
// Implementations can use SMTP, SendGrid, SES, Postmark, etc.
type EmailProvider interface {
	Send(to string, subject string, htmlBody string) error
}

// SmtpEmailProvider sends emails via SMTP (the default).
type SmtpEmailProvider struct{}

func (s *SmtpEmailProvider) Send(to string, subject string, htmlBody string) error {
	from := os.Getenv("SMTP_FROM_ADDRESS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		from, to, subject, htmlBody)

	addr := host + ":" + port
	return smtp.SendMail(addr, smtpAuth, from, []string{to}, []byte(msg))
}

var emailProvider EmailProvider

func emailProviderConfigure() error {
	provider := os.Getenv("EMAIL_PROVIDER")
	if provider == "" {
		provider = "smtp"
	}

	switch provider {
	case "smtp":
		if smtpConfigured {
			emailProvider = &SmtpEmailProvider{}
			logger.Infof("email provider: smtp")
		}
	default:
		logger.Warningf("unknown email provider %s, falling back to smtp", provider)
		if smtpConfigured {
			emailProvider = &SmtpEmailProvider{}
		}
	}

	return nil
}
