package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// func optionsHandler(w http.ResponseWriter, r *http.Request) {}

func (router *Router) initializeRoutes() {

	// Middleware is added on here. Currently there is no middleware
	middleware := alice.New(CORSHandler)

	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	// router.Options("*", middleware.ThenFunc(optionsHandler))

	router.Get("/", middleware.ThenFunc(indexHandler))
	router.Post("/", middleware.ThenFunc(indexHandler))
	// router.Options("/", middleware.ThenFunc(CORSHandler))

	router.Get("/get", middleware.ThenFunc(getCommentsHandler))
	router.Post("/create", middleware.ThenFunc(createCommentHandler))
	// router.Get("/comments", getCommentsHandler)
	// router.Post("/comments", createCommentHandler)
}
