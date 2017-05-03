package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
)

var (
	logger = logging.MustGetLogger("commento")
	db     *sql.DB
	port   *int
	isDemo *bool
)

func init() {
	port = flag.Int("port", 8080, "HTTP server port.")
	isDemo = flag.Bool("demo", false, "Use commento demo server.")

	flag.Parse()
}

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

	if *isDemo {
		fmt.Println("clearing")
		go func() {
			for true {
				err := cleanupOldComments()
				if err != nil {
					logger.Errorf("Error cleaning up old comments %v", err)
				}
				time.Sleep(60 * time.Second)
				fmt.Println("deleting")
			}
		}()
	}

	svr := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Infof("Running on port %d", *port)
	err = svr.ListenAndServe()
	if err != nil {
		logger.Fatalf("http.ListenAndServe: %v", err)
	}

}
