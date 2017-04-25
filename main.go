package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
	"net/http"
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

	http.ListenAndServe(":8080", nil)
}
