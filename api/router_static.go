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

type staticAssetPlugs struct {
	Origin string
}

type staticHtmlPlugs struct {
	CdnPrefix string
}

func initStaticRouter(router *mux.Router) error {
	asset := make(map[string]string)

	for _, dir := range []string{"js", "css", "images"} {
		files, err := ioutil.ReadDir("./" + dir)
		if err != nil {
			logger.Errorf("cannot read directory ./%s: %v", dir, err)
			return err
		}

		for _, file := range files {
			sl := string(os.PathSeparator)
			path := sl + dir + sl + file.Name()

			contents, err := ioutil.ReadFile("." + path)
			if err != nil {
				logger.Errorf("cannot read file %s: %v", path, err)
				return err
			}

			if dir == "js" {
				asset[path] = "window.API='" + os.Getenv("ORIGIN") + "/api';" + string(contents)
			} else {
				asset[path] = string(contents)
			}

			logger.Debugf("routing %s", path)
			router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, asset[r.URL.Path])
			})
		}
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
			logger.Errorf("cannot parse /%s template: %v", page, err)
			return err
		}

		var buf bytes.Buffer
		t.Execute(&buf, &staticHtmlPlugs{CdnPrefix: os.Getenv("CDN_PREFIX")})

		html["/" + page] = buf.String()
	}

	for _, page := range pages {
		router.HandleFunc("/" + page, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, html[r.URL.Path])
		})
	}

	router.HandleFunc("/", redirectLogin)

	return nil
}
