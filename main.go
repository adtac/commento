package main

import (
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := LoadDatabase("sqlite3.db")
	if err != nil {
		Die(err)
	}
	
	fs := http.FileServer(http.Dir("assets"))
	
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/create", CreateCommentHandler)
	http.HandleFunc("/get", GetCommentsHandler)
	
	var port string
	
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = ":" + fromEnv
	} else {
		port = ":8080"
	}
	
	if demoEnv := os.Getenv("DEMO"); demoEnv == "true" {
		t := time.Second * 60
		Logger.Infof("Demo Env: Cleaning old comments every %s", t)
		go func() {
			for true {
				err := CleanupOldComments()
				if err != nil {
					Logger.Errorf("Error cleaning up old comments %s", err)
				}
				time.Sleep(t)
			}
		}()
	}
	
	svr := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	Logger.Infof("Running on port %s", port)
	err = svr.ListenAndServe()
	if err != nil {
		Logger.Fatalf("http.ListenAndServe: %v", err)
	}
}
