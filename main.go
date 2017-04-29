package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
)

const (
	portDefault = "8080"
)

var logger = logging.MustGetLogger("commento")
var db *sql.DB

func main() {
	if err := loadDatabase("sqlite3.db"); err != nil {
		die(err)
	}

	// Parse command line options
	var port = flag.String("port", portDefault, "port for commento service")
	flag.Parse()

        // Check that we get a valid port value
	if _, err := strconv.ParseInt(*port, 10, 32); err != nil {
		die(err)
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

	router := NewRouter()
	svr := &http.Server{
		Addr:          ":" + *port,
		Handler:       router,
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  10 * time.Second,
	}

	logger.Infof("Running on port %s", *port)
	if err := svr.ListenAndServe(); err != nil {
		logger.Fatalf("http.ListenAndServe: %v", err)
	}

}
