package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"strings"

	"encoding/json"

	"golang.org/x/net/html"

	"time"

	"github.com/go-redis/redis"
)

const defaultPort = "80"
const headerContentType = "Content-Type"
const contentTypeHTML = "text/html"
const contentTypeJSON = "application/json; charset=utf-8"

//PageSummary contains summary information about a web page
type PageSummary struct {
	Title string   `json:"title"`
	Links []string `json:"links"`
}

//getPageSummary fetches PageSummary info for a given URL
func getPageSummary(URL string) (*PageSummary, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error resonse status code: %d", resp.StatusCode)
	}
	if !strings.HasPrefix(resp.Header.Get(headerContentType), contentTypeHTML) {
		return nil, fmt.Errorf("the URL did not return an HTML page")
	}

	psum := &PageSummary{}
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		ttype := tokenizer.Next()
		if ttype == html.ErrorToken {
			return psum, tokenizer.Err()
		}

		//if this is a start tag token
		if ttype == html.StartTagToken {
			token := tokenizer.Token()
			//if this is the page title
			if token.Data == "title" {
				tokenizer.Next()
				psum.Title = tokenizer.Token().Data
			}

			//if this is a hyperlink
			if token.Data == "a" {
				//get the href attribute
				for _, attr := range token.Attr {
					//ignore bookmark links
					if attr.Key == "href" && !strings.HasPrefix(attr.Val, "#") {
						psum.Links = append(psum.Links, attr.Val)
					}
				} //for all attributes
			} //if <a>
		} //if start tag
	} //for each token
} //getPageSummary()

// Setting this as a reciever for all functions allows all functions access to redis
type HandlerContext struct {
	redisClient *redis.Client
}

//SummaryHandler handles the /v1/summary resource
func (context *HandlerContext) SummaryHandler(w http.ResponseWriter, r *http.Request) {
	URL := r.FormValue("url")
	if len(URL) == 0 {
		http.Error(w, "please supply a `url` query string parameter", http.StatusBadRequest)
		return
	}

	jsonBuffer, err := context.redisClient.Get(URL).Bytes()
	if err != nil && err != redis.Nil {
		http.Error(w, "error getting from cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Cacheing data if it didn't already exist in redis
	if err == redis.Nil {
		//TODO: call getPageSummary() passing URL
		//marshal struct into JSON, and write it
		//to the response
		pageSummary, err := getPageSummary(URL)
		if err != nil && err != io.EOF {
			http.Error(w, "error getting page summary: "+err.Error(), http.StatusInternalServerError)
			return
		}

		jsonBuffer, err = json.Marshal(pageSummary)
		if err != nil {
			http.Error(w, "error marshaling json: "+err.Error(), http.StatusInternalServerError)
			return
		}

		context.redisClient.Set(URL, jsonBuffer, time.Second*60)
	}

	w.Header().Add(headerContentType, contentTypeJSON)
	w.Write(jsonBuffer)
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}
	addr := host + ":" + port

	redisOptions := redis.Options{
		Addr: "localhost:6379", // Go requires that you have a comma at the end
	}
	redisClient := redis.NewClient(&redisOptions)
	handlerContext := &HandlerContext{
		redisClient: redisClient,
	}

	http.HandleFunc("/v1/summary", handlerContext.SummaryHandler) // Summary handler will have access to redis client

	fmt.Printf("listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
