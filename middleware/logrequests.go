package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func logReq(r *http.Request) {
	log.Println(r.Method, r.URL.Path)
}

func logReqs(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		start := time.Now()
		handlerFunc(w, r) // Closure allows us to use this variable inside this function
		fmt.Printf("%v\n", time.Since(start))
	}
}

// func logRequests(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s", r.Method, r.URL.Path)
// 		start := time.Now()
// 		handler.ServeHTTP(w, r)
// 		fmt.Printf("%v\n", time.Since(start))
// 	})
// }

func logRequests(logger *log.Logger) Adapter {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s", r.Method, r.URL.Path)
			start := time.Now()
			handler.ServeHTTP(w, r)
			logger.Printf("%v\n", time.Since(start))
		})
	}
}
