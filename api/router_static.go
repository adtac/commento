package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
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
		sl := string(os.PathSeparator)
		dir = sl + dir

		files, err := ioutil.ReadDir(os.Getenv("STATIC") + dir)
		if err != nil {
			logger.Errorf("cannot read directory %s%s: %v", os.Getenv("STATIC"), dir, err)
			return err
		}

		for _, file := range files {
			p := dir + sl + file.Name()

			contents, err := ioutil.ReadFile(os.Getenv("STATIC") + p)
			if err != nil {
				logger.Errorf("cannot read file %s%s: %v", os.Getenv("STATIC"), p, err)
				return err
			}

			prefix := ""
			if dir == "/js" {
				prefix = "window.commento_origin='" + os.Getenv("ORIGIN") + "';\n"
				prefix += "window.commento_cdn='" + os.Getenv("CDN_PREFIX") + "';\n"
			}

			asset[p] = prefix + string(contents)

			router.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(r.URL.Path)))
				fmt.Fprint(w, asset[r.URL.Path])
			})
		}
	}

	pages := []string{
		"login",
		"forgot",
		"reset-password",
		"signup",
		"dashboard",
		"logout",
	}

	html := make(map[string]string)
	for _, page := range pages {
		sl := string(os.PathSeparator)
		page = sl + page
		file := page + ".html"

		contents, err := ioutil.ReadFile(os.Getenv("STATIC") + file)
		if err != nil {
			logger.Errorf("cannot read file %s%s: %v", os.Getenv("STATIC"), file, err)
			return err
		}

		t, err := template.New(page).Delims("<<<", ">>>").Parse(string(contents))
		if err != nil {
			logger.Errorf("cannot parse %s%s template: %v", os.Getenv("STATIC"), file, err)
			return err
		}

		var buf bytes.Buffer
		t.Execute(&buf, &staticHtmlPlugs{CdnPrefix: os.Getenv("CDN_PREFIX")})

		html[page] = buf.String()
	}

	for _, page := range pages {
		router.HandleFunc("/"+page, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, html[r.URL.Path])
		})
	}

	router.HandleFunc("/", redirectLogin)

	return nil
}
