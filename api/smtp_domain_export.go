package main

import (
	"bytes"
	"net/smtp"
	"os"
)

type domainExportPlugs struct {
	Origin    string
	Domain    string
	ExportHex string
}

func smtpDomainExport(to string, toName string, domain string, exportHex string) error {
	var header bytes.Buffer
	headerTemplate.Execute(&header, &headerPlugs{FromAddress: os.Getenv("SMTP_FROM_ADDRESS"), ToAddress: to, ToName: toName, Subject: "Commento Data Export"})

	var body bytes.Buffer
	templates["domain-export"].Execute(&body, &domainExportPlugs{Origin: os.Getenv("ORIGIN"), ExportHex: exportHex})

	err := smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), smtpAuth, os.Getenv("SMTP_FROM_ADDRESS"), []string{to}, concat(header, body))
	if err != nil {
		logger.Errorf("cannot send data export email: %v", err)
		return errorCannotSendEmail
	}

	return nil
}
