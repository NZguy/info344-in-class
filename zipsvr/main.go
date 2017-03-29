package main

import "os"
import "log"
import "net/http"
import "fmt"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	w.Header().Add("Content-Type", "text/plain")

	w.Write([]byte("Hello " + name))
}

func main() {
	// var addr string = is the same thing
	// := automatically figures out the variables type
	// No loose typing, addr will always be a string
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		log.Fatal("please set ADDR environment variable")
	}

	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
