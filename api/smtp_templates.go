package main

import (
	"os"
	"fmt"
	"html/template"
)

var headerTemplate *template.Template

type headerPlugs struct {
	FromAddress string
	ToName      string
	ToAddress   string
	Subject     string
}

var templates map[string]*template.Template

func loadTemplates() error {
	var err error
	headerTemplate, err = template.New("header").Parse(`MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8
From: {{.FromAddress}}
To: {{.ToName}} <{{.ToAddress}}>
Subject: {{.Subject}}

`)
	if err != nil {
		logger.Errorf("cannot parse header template: %v", err)
		return errorMalformedTemplate
	}

	names := []string{"confirm-hex", "reset-hex"}

	templates = make(map[string]*template.Template)

	logger.Infof("loading templates: %v", names)
	for _, name := range names {
		var err error
		templates[name] = template.New(name)
		templates[name], err = template.ParseFiles(fmt.Sprintf("%s/templates/%s.html", os.Getenv("COMMENTO_STATIC"), name))
		if err != nil {
			logger.Errorf("cannot parse %s/templates/%s.html: %v", os.Getenv("STATIC"), name, err)
			return errorMalformedTemplate
		}
	}

	return nil
}
