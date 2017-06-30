package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	err := LoadDatabase("sqlite:file=sqlite3.db")
	if err != nil {
		Die(err)
	}

	// Load configuration from the environment.
	// Values in earlier files will take precedence over later values
	//    Ex. A COMMENTO_PORT value in .env.development.local will be used
	//        even if COMMENTO_PORT exists in a .env.development file
	for _, envFile := range []string{".env.development.local", ".env.test.local", ".env.production.local", ".env.local", ".env.development", ".env.test", ".env.production", ".env"} {
		godotenv.Load(envFile)
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
	Logger.Infof("Running on port %s", port)
	err = svr.ListenAndServe()
	if err != nil {
		Logger.Fatalf("http.ListenAndServe: %v", err)
	}
}
