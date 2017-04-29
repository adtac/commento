package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
)

const (
	portDefault = "8080"
)

var logger = logging.MustGetLogger("commento")
var db *sql.DB

// The run environment for the application. Either dev, staging, or prod
var runEnv string

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

	// Check for the run environment, default to dev
	runEnv = strings.ToLower(os.Getenv("RUN_ENV"))
	switch runEnv {
	case "": // When the run environment is not set it is assumed to be dev mode
		runEnv = "dev"
	case "dev":
	case "staging":
	case "prod":
	default:
		logger.Errorf("Unrecognized run environment: %s", runEnv)
		panic("Bad run environment")
	}

	logger.Infof("Run environment: %s", runEnv)

	if demoEnv := os.Getenv("DEMO"); demoEnv == "true" {
		go func() {
			for {
				err := cleanupOldComments()
				if err != nil {
					logger.Errorf("Error cleaning up old comments %v", err)
				}
				time.Sleep(t)
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
