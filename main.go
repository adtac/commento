package main

import (
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
)

var Logger = logging.MustGetLogger("commento")

func main() {
	var err error

	err = loadConfig()
	if err != nil {
		Logger.Errorf("fatal error: cannot load config: %v\n", err)
		return
	}

	connectionStr := "sqlite:file=" + os.Getenv("COMMENTO_DATABASE_FILE")

	err = LoadDatabase(connectionStr)
	if err != nil {
		Logger.Errorf("fatal error: cannot load %s: %v\n", connectionStr, err)
		return
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

	Logger.Infof("Starting the server on port %s", port)
	err = svr.ListenAndServe()
	if err != nil {
		Logger.Errorf("fatal error: cannot start the server: %v\n", err)
		return
	}
}
