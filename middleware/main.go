package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := "localhost:4000"

	mux := http.NewServeMux()
	muxLogged := http.NewServeMux()
	muxLogged.HandleFunc("/v1/hello1", HelloHandler1)
	muxLogged.HandleFunc("/v1/hello2", HelloHandler2)

	mux.HandleFunc("/v1/hello3", HelloHandler3)
	logger := log.New(os.Stdout, "", log.LstdFlags)
	mux.Handle("/v1/", Adapt(muxLogged,
		logRequests(logger),
		throttleRequests(4, time.Minute))) // Add multiple mux middlewares with an Adapt function

	//mux.Handle("/v1/", logRequests(logger)(muxLogged)) //logRequests(logger) returns a function that is immediately executed
	//mux.Handle("/v1/", logRequests(muxLogged))
	// Only handler 1 and 2 will log requests

	//http.HandleFunc("/v1/hello1", logReqs(HelloHandler1)) // Option 2: Wrap around handler
	// http.HandleFunc("/v1/hello1", HelloHandler1)
	// http.HandleFunc("/v1/hello2", HelloHandler2)
	// http.HandleFunc("/v1/hello3", HelloHandler3)

	fmt.Printf("listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))

	//log.Fatal(http.ListenAndServe(addr, logRequests(mux))) // Option 3: wrap whole mux
}
