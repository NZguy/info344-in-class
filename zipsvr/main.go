package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type zip struct {
	Zip   string `json:"zip"`
	City  string `json:"city"`
	State string `json:"state"`
}

type zipSlice []*zip
type zipIndex map[string]zipSlice

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	w.Header().Add("Content-Type", "text/plain")

	w.Write([]byte("Hello " + name))
}

// http.ResponseWriter is an interface and is always passed by reference
// http.Request needs a * to tell it to pass the pointer
// zi is equivalent to this in java, you need to pass in the parent struct
func (zi zipIndex) zipsForCityHandler(w http.ResponseWriter, r *http.Request) {
	_, city := path.Split(r.URL.Path)
	lcity := strings.ToLower(city)

	// Header for content type
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	// Any origin can call our function
	w.Header().Add("Access-Control-Allow-Origin", "*")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(zi[lcity]); err != nil {
		http.Error(w, "error encoding json: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// var addr string = is the same thing
	// := automatically figures out the variables type
	// No loose typing, addr will always be a string
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		log.Fatal("please set ADDR environment variable")
	}

	f, err := os.Open("../data/zips.json")
	if err != nil {
		log.Fatal("error opening zips file " + err.Error())
	}

	zips := make(zipSlice, 0, 43000)
	decoder := json.NewDecoder(f)
	// decoder.Decode may return an error, if the error != nil then there was an error
	if err := decoder.Decode(&zips); err != nil {
		log.Fatal("error decoding zips json " + err.Error())
	}
	fmt.Printf("loaded %d zips\n", len(zips))

	zi := make(zipIndex) // map joining city string to zip object pointer

	// range is like foreach
	// _ means ignore, we are ignoring the index variable
	for _, z := range zips {
		lower := strings.ToLower(z.City)
		zi[lower] = append(zi[lower], z)
	}

	fmt.Printf("there are %d zips in Seattle\n", len(zi["seattle"]))

	http.HandleFunc("/hello", helloHandler)

	http.HandleFunc("/zips/city/", zi.zipsForCityHandler)

	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
