package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func redirectLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", 301)
}

type staticHtmlPlugs struct {
	CdnPrefix string
}

func initStaticRouter(router *mux.Router) error {
	for _, path := range []string{"js", "css", "images"} {
		router.PathPrefix("/" + path + "/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f, err := os.Stat("." + r.URL.Path)
			if err != nil || f.IsDir() {
				http.NotFound(w, r)
			}

			http.ServeFile(w, r, "."+r.URL.Path)
		})
	}

	pages := []string{
		"login",
		"signup",
		"dashboard",
		"account",
	}

	html := make(map[string]string)
	for _, page := range pages {
		contents, err := ioutil.ReadFile(page + ".html")
		if err != nil {
			logger.Errorf("cannot read file %s.html: %v", page, err)
			return err
		}

		t, err := template.New(page).Delims("<<<", ">>>").Parse(string(contents))
		if err != nil {
			logger.Errorf("cannot parse %s.html template: %v", page, err)
			return err
		}

		var buf bytes.Buffer
		t.Execute(&buf, &staticHtmlPlugs{CdnPrefix: os.Getenv("CDN_PREFIX")})

		html[page] = buf.String()
	}

	for _, page := range pages {
		router.HandleFunc("/"+page, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, html[page])
		})
	}

	router.HandleFunc("/", redirectLogin).Methods("GET")

	return nil
}
