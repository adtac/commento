package main

import (
	"bytes"
	"net/smtp"
	"os"
)

type domainExportErrorPlugs struct {
	Origin string
	Domain string
}

func smtpDomainExportError(to string, toName string, domain string) error {
	var header bytes.Buffer
	headerTemplate.Execute(&header, &headerPlugs{FromAddress: os.Getenv("SMTP_FROM_ADDRESS"), ToAddress: to, ToName: toName, Subject: "Commento Data Export"})

	var body bytes.Buffer
	templates["data-export-error"].Execute(&body, &domainExportPlugs{Origin: os.Getenv("ORIGIN")})

	err := smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), smtpAuth, os.Getenv("SMTP_FROM_ADDRESS"), []string{to}, concat(header, body))
	if err != nil {
		logger.Errorf("cannot send data export error email: %v", err)
		return errorCannotSendEmail
	}

	return nil
}
