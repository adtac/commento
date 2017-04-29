package main

import (
	"context"
	"net/http"

	httprouter "github.com/julienschmidt/httprouter"
)

type Router struct {
	httprouter.Router
}

func NewRouter() *Router {

	baseRouter := httprouter.New()

	router := Router{
		*baseRouter,
	}
	router.initializeRoutes()

	return &router
}

func (router *Router) Get(path string, handler http.Handler) {
	router.GET(path, wrapHandler(handler))
}

func (router *Router) Post(path string, handler http.Handler) {
	router.POST(path, wrapHandler(handler))
}

func (router *Router) Put(path string, handler http.Handler) {
	router.PUT(path, wrapHandler(handler))
}

func (router *Router) Delete(path string, handler http.Handler) {
	router.DELETE(path, wrapHandler(handler))
}

func (router *Router) Options(path string, handler http.Handler) {
	router.OPTIONS(path, wrapHandler(handler))
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if runEnv == "dev" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		ctx := context.WithValue(r.Context(), "params", ps)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
