package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

func redirectLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, os.Getenv("ORIGIN")+"/login", 301)
}

type staticPlugs struct {
	Origin    string
	CdnPrefix string
	Footer    string
}

var asset map[string][]byte = make(map[string][]byte)
var contentType map[string]string = make(map[string]string)
var footer string
var compress bool

func fileDetemplate(f string) ([]byte, error) {
	contents, err := ioutil.ReadFile(f)
	if err != nil {
		logger.Errorf("cannot read file %s: %v", f, err)
		return []byte{}, err
	}

	x := string(contents)
	x = strings.Replace(x, "[[[.Origin]]]", os.Getenv("ORIGIN"), -1)
	x = strings.Replace(x, "[[[.CdnPrefix]]]", os.Getenv("CDN_PREFIX"), -1)
	x = strings.Replace(x, "[[[.Footer]]]", footer, -1)
	x = strings.Replace(x, "[[[.Version]]]", version, -1)

	return []byte(x), nil
}

func footerInit() error {
	contents, err := fileDetemplate(os.Getenv("STATIC") + "/footer.html")
	if err != nil {
		logger.Errorf("cannot init footer: %v", err)
		return err
	}

	footer = string(contents)
	return nil
}

func fileLoad(f string) ([]byte, error) {
	b, err := fileDetemplate(f)
	if err != nil {
		logger.Errorf("cannot load file %s: %v", f, err)
		return []byte{}, err
	}

	if !compress {
		return b, nil
	}

	return gzipStatic(b)
}

func staticRouterInit(router *mux.Router) error {
	var err error

	subdir := pathStrip(os.Getenv("ORIGIN"))

	if err = footerInit(); err != nil {
		logger.Errorf("error initialising static router: %v", err)
		return err
	}

	for _, dir := range []string{"/js", "/css", "/images", "/fonts"} {
		files, err := ioutil.ReadDir(os.Getenv("STATIC") + dir)
		if err != nil {
			logger.Errorf("cannot read directory %s%s: %v", os.Getenv("STATIC"), dir, err)
			return err
		}

		for _, file := range files {
			f := dir + "/" + file.Name()
			asset[subdir+f], err = fileLoad(os.Getenv("STATIC") + f)
			if err != nil {
				logger.Errorf("cannot detemplate %s%s: %v", os.Getenv("STATIC"), f, err)
				return err
			}
		}
	}

	pages := []string{
		"/login",
		"/forgot",
		"/reset",
		"/signup",
		"/confirm-email",
		"/unsubscribe",
		"/dashboard",
		"/settings",
		"/logout",
		"/profile",
	}

	for _, page := range pages {
		f := page + ".html"
		asset[subdir+page], err = fileLoad(os.Getenv("STATIC") + f)
		if err != nil {
			logger.Errorf("cannot detemplate %s%s: %v", os.Getenv("STATIC"), f, err)
			return err
		}
	}

	for p, _ := range asset {
		if path.Ext(p) != "" {
			contentType[p] = mime.TypeByExtension(path.Ext(p))
		} else {
			contentType[p] = "text/html; charset=utf-8"
		}

		router.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType[r.URL.Path])
			if compress {
				w.Header().Set("Content-Encoding", "gzip")
			}
			w.Write(asset[r.URL.Path])
		})
	}

	router.HandleFunc("/", redirectLogin)

	return nil
}
