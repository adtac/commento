package main

import (
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error

	err = loadConfig()
	if err != nil {
		Die(err)
	}

	fp := os.Getenv("COMMENTO_DATABASE_FILE")
	err = LoadDatabase("sqlite:file=" + fp)
	if err != nil {
		Die(err)
	}

	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/create", CreateCommentHandler)
	http.HandleFunc("/get", GetCommentsHandler)

	port := os.Getenv("COMMENTO_PORT")
	svr := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	Logger.Infof("Running on port %s", port)
	err = svr.ListenAndServe()
	if err != nil {
		Die(err)
	}
}
