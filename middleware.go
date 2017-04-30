package main

import (
	"fmt"
	"net/http"
)

func CORSHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Cookies")
		for _, cookie := range r.Cookies() {
			fmt.Fprint(w, cookie.Name)
		}

		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)		

		next.ServeHTTP(w, r)
	})
}
