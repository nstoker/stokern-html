package main

import (
	"log"
	"net/http"
)

// Logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
