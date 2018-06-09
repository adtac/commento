package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func serveRoutes() error {
	router := mux.NewRouter()

	if err := initAPIRouter(router); err != nil {
		return err
	}

	if err := initStaticRouter(router); err != nil {
		return err
	}

	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})

	addrPort := os.Getenv("BIND_ADDRESS") + ":" + os.Getenv("PORT")

	logger.Infof("starting server on %s\n", addrPort)
	if err := http.ListenAndServe(addrPort, handlers.CORS(origins, headers, methods)(router)); err != nil {
		logger.Errorf("cannot start server: %v", err)
		return err
	}

	return nil
}
