package main

import (
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
	"database/sql"
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

    http.ListenAndServe(":8080", nil)
}
