package main

import (
	"fmt"
	"os"
	"time"
)

const usage = `
usage:
	crawler <starting-url>
`

// Make sure to honor robots.txt on sites
func worker(linkq chan string, resultsq chan []string) {
	for link := range linkq {
		plinks, err := getPageLinks(link)
		if err != nil {
			fmt.Printf("ERROR fetching %s: %v", link, err)
			continue // Goes to next iteration of for loop
		}

		fmt.Printf("%s (%d links)\n", link, len(plinks.Links))
		time.Sleep(time.Millisecond * 500) // Time to wait between individual worker requests
		if len(plinks.Links) > 0 {
			// Define a temporary go function to stop the workers from blocking when attempting to write to channer
			// Define and immediately invoke function
			go func(links []string) {
				resultsq <- links
			}(plinks.Links)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	nWorkers := 20                        // More workers would go faster
	linkq := make(chan string, 1000)      // Higher number means larger memory buffer
	resultsq := make(chan []string, 1000) // Workers will get blocked less often and go faster
	for i := 0; i < nWorkers; i++ {
		go worker(linkq, resultsq)
	}

	linkq <- os.Args[1] // Add an initial link to begin crawling from

	seen := map[string]bool{} // This is safe to use in concurrency because only main ever accesses it
	for links := range resultsq {
		for _, link := range links {
			if !seen[link] {
				seen[link] = true
				linkq <- link
			}
		}
	}

}
