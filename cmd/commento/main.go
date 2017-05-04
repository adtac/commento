package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adtac/commento"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dir  = flag.String("dir", "../..", "commento working directory")
	port = flag.Int("port", 8080, "commento http port")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := os.Chdir(*dir); err != nil {
		log.Fatalf("unable to change to commento directory: %v", err)
	}

	if err := commento.LoadDatabase("sqlite3.db"); err != nil {
		log.Fatalf("failed to load database: %v", err)
	}

	if os.Getenv("DEMO") == "true" {
		go func() {
			log.Print("Demo Env: cleaning up old comments")
			for {
				if err := commento.CleanupOldComments(); err != nil {
					log.Printf("Error cleaning up old comments: %v", err)
				}
				time.Sleep(time.Minute)
			}
		}()
	}

	listenAddr := fmt.Sprintf(":%d", *port)
	if err := commento.Serve(listenAddr); err != nil {
		log.Fatal(err)
	}

}
