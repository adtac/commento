package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (router *Router) initializeRoutes() {

	// Middleware is added on here. Currently there is no middleware
	middleware := alice.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.Get("/", middleware.ThenFunc(indexHandler))
	router.Post("/", middleware.ThenFunc(indexHandler))

	router.Post("/get", middleware.ThenFunc(getCommentsHandler))
	router.Post("/create", middleware.ThenFunc(createCommentHandler))
}
