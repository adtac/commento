package main

import (
	"database/sql"
	"net/http"
	"os"

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

	err = http.ListenAndServe(port, nil)
	if err != nil {
		logger.Fatalf("http.ListenAndServe: %v", err)
	}
}
