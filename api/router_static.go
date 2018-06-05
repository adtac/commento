package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"mime"
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
			p := sl + dir + sl + file.Name()

			contents, err := ioutil.ReadFile("." + p)
			if err != nil {
				logger.Errorf("cannot read file %s: %v", p, err)
				return err
			}

			prefix := ""
			if dir == "js" {
				prefix = "window.origin='" + os.Getenv("ORIGIN") + "';\n"
				prefix += "window.cdn='" + os.Getenv("CDN_PREFIX") + "';\n"
			}

			asset[p] = prefix + string(contents);

			router.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(r.URL.Path)))
				fmt.Fprint(w, asset[r.URL.Path])
			})
		}
	}

	pages := []string{
		"login",
		"signup",
		"dashboard",
		"logout",
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
