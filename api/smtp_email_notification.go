package main

import (
	"bytes"
	"fmt"
	ht "html/template"
	"net/smtp"
	"os"
	tt "text/template"
)

type emailNotificationText struct {
	emailNotification
	Html ht.HTML
}

type emailNotificationPlugs struct {
	Origin               string
	Kind                 string
	Subject              string
	UnsubscribeSecretHex string
	Notifications        []emailNotificationText
}

func smtpEmailNotification(to string, toName string, unsubscribeSecretHex string, notifications []emailNotificationText, kind string) error {
	var subject string
	if kind == "reply" {
		var verb string
		if len(notifications) > 1 {
			verb = "replies"
		} else {
			verb = "reply"
		}
		subject = fmt.Sprintf("%d new comment %s", len(notifications), verb)
	} else {
		var verb string
		if len(notifications) > 1 {
			verb = "comments"
		} else {
			verb = "comment"
		}
		if kind == "pending-moderation" {
			subject = fmt.Sprintf("%d new %s pending moderation", len(notifications), verb)
		} else {
			subject = fmt.Sprintf("%d new %s on your website", len(notifications), verb)
		}
	}

	h, err := tt.New("header").Parse(`MIME-Version: 1.0
From: Commento <{{.FromAddress}}>
To: {{.ToName}} <{{.ToAddress}}>
Content-Type: text/html; charset=UTF-8
Subject: {{.Subject}}

`)

	var header bytes.Buffer
	h.Execute(&header, &headerPlugs{FromAddress: os.Getenv("SMTP_FROM_ADDRESS"), ToAddress: to, ToName: toName, Subject: "[Commento] " + subject})

	t, err := ht.ParseFiles(fmt.Sprintf("%s/templates/email-notification.txt", os.Getenv("STATIC")))
	if err != nil {
		logger.Errorf("cannot parse %s/templates/email-notification.txt: %v", os.Getenv("STATIC"), err)
		return errorMalformedTemplate
	}

	var body bytes.Buffer
	err = t.Execute(&body, &emailNotificationPlugs{
		Origin:               os.Getenv("ORIGIN"),
		Kind:                 kind,
		Subject:              subject,
		UnsubscribeSecretHex: unsubscribeSecretHex,
		Notifications:        notifications,
	})
	if err != nil {
		logger.Errorf("error generating templated HTML for email notification: %v", err)
		return err
	}

	err = smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), smtpAuth, os.Getenv("SMTP_FROM_ADDRESS"), []string{to}, concat(header, body))
	if err != nil {
		logger.Errorf("cannot send email notification: %v", err)
		return errorCannotSendEmail
	}

	return nil
}
