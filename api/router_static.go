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
	http.Redirect(w, r, os.Getenv("ORIGIN")+"/login", 301)
}

type staticAssetPlugs struct {
	Origin string
}

type staticHtmlPlugs struct {
	Origin    string
	CdnPrefix string
	Footer    template.HTML
}

func staticRouterInit(router *mux.Router) error {
	asset := make(map[string][]byte)
	gzippedAsset := make(map[string][]byte)

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
				prefix = "window.commentoOrigin='" + os.Getenv("ORIGIN") + "';\n"
				prefix += "window.commentoCdn='" + os.Getenv("CDN_PREFIX") + "';\n"
			}

			gzip := (os.Getenv("GZIP_STATIC") == "true")

			subdir := pathStrip(os.Getenv("ORIGIN"))

			asset[subdir+p] = []byte(prefix + string(contents))
			if gzip {
				gzippedAsset[subdir+p], err = gzipStatic(asset[subdir+p])
				if err != nil {
					logger.Errorf("error gzipping %s: %v", p, err)
					return err
				}
			}

			// faster than checking inside the handler
			if !gzip {
				router.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(r.URL.Path)))
					w.Write(asset[r.URL.Path])
				})
			} else {
				router.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(r.URL.Path)))
					w.Header().Set("Content-Encoding", "gzip")
					w.Write(gzippedAsset[r.URL.Path])
				})
			}
		}
	}

	footer, err := ioutil.ReadFile(os.Getenv("STATIC") + string(os.PathSeparator) + "footer.html")
	if err != nil {
		logger.Errorf("cannot read file footer.html: %v", err)
		return err
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

		t, err := template.New(page).Delims("[[[", "]]]").Parse(string(contents))
		if err != nil {
			logger.Errorf("cannot parse %s%s template: %v", os.Getenv("STATIC"), file, err)
			return err
		}

		var buf bytes.Buffer
		t.Execute(&buf, &staticHtmlPlugs{
			Origin:    os.Getenv("ORIGIN"),
			CdnPrefix: os.Getenv("CDN_PREFIX"),
			Footer:    template.HTML(string(footer)),
		})

		subdir := pathStrip(os.Getenv("ORIGIN"))

		html[subdir+page] = buf.String()
	}

	for _, page := range pages {
		router.HandleFunc("/"+page, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, html[r.URL.Path])
		})
	}

	router.HandleFunc("/", redirectLogin)

	return nil
}
