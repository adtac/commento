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
				cleanupOldComments()
				time.Sleep(60 * time.Second)
				fmt.Println("deleting")
			}
		}()
	}

	err = http.ListenAndServe(port, nil)
	if err != nil {
		logger.Fatalf("http.ListenAndServe: %v", err)
	}

}
