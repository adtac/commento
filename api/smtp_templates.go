package main

import (
	"fmt"
	"os"
	"text/template"
)

var headerTemplate *template.Template

type headerPlugs struct {
	FromAddress string
	ToName      string
	ToAddress   string
	Subject     string
}

var templates map[string]*template.Template

func smtpTemplatesLoad() error {
	var err error
	headerTemplate, err = template.New("header").Parse(`MIME-Version: 1.0
From: Commento <{{.FromAddress}}>
To: {{.ToName}} <{{.ToAddress}}>
Content-Type: text/plain; charset=UTF-8
Subject: {{.Subject}}

`)
	if err != nil {
		logger.Errorf("cannot parse header template: %v", err)
		return errorMalformedTemplate
	}

	names := []string{
		"confirm-hex",
		"reset-hex",
		"domain-export",
		"domain-export-error",
	}

	templates = make(map[string]*template.Template)

	logger.Infof("loading templates: %v", names)
	for _, name := range names {
		var err error
		templates[name] = template.New(name)
		templates[name], err = template.ParseFiles(fmt.Sprintf("%s/templates/%s.txt", os.Getenv("STATIC"), name))
		if err != nil {
			logger.Errorf("cannot parse %s/templates/%s.txt: %v", os.Getenv("STATIC"), name, err)
			return errorMalformedTemplate
		}
	}

	return nil
}
