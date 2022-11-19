package main

import (
	"crypto/tls"
	"errors"
	"net/smtp"
	"os"
	"strings"
)

var testHookStartTLS func(*tls.Config) // nil, except for tests

func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func SendSMTPMail(to string, msg []byte) error {

	from := os.Getenv("SMTP_FROM_ADDRESS")
	host := os.Getenv("SMTP_HOST")
	addr := host + ":" + os.Getenv("SMTP_PORT")

	if err := validateLine(from); err != nil {
		return err
	}
	if err := validateLine(to); err != nil {
		return err
	}

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	localName := os.Getenv("SMTP_HELO_FQDN")
	if localName == "" {
		localName, err = os.Hostname()
		if err != nil {
			localName = "localhost"
		}
	}
	if err := c.Hello(localName); err != nil {
		return err
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host}
		if testHookStartTLS != nil {
			testHookStartTLS(config)
		}
		if err := c.StartTLS(config); err != nil {
			return err
		}
	}

	if smtpAuth != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}
		if err := c.Auth(smtpAuth); err != nil {
			return err
		}
	}

	if err := c.Mail(from); err != nil {
		return err
	}

	if err := c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
