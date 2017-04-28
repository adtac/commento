package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("commento")
var db *sql.DB

func main() {
	err := loadDatabase("sqlite3.db")
	if err != nil {
		die(err)
	}

	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createCommentHandler)
	http.HandleFunc("/get", getCommentsHandler)

	var port string

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = ":" + fromEnv
	} else {
		port = ":8080"
	}

	if demoEnv := os.Getenv("DEMO"); demoEnv == "true" {
		fmt.Println("clearing")
		go func() {
			for true {
				err := cleanupOldComments()
				if err != nil {
					logger.Errorf("Error cleaning up old comments %s", err)
				}
				time.Sleep(60 * time.Second)
				fmt.Println("deleting")
			}
		}()
	}

	svr := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Infof("Running on port %s", port)
	err = svr.ListenAndServe()
	if err != nil {
		logger.Fatalf("http.ListenAndServe: %v", err)
	}

}
