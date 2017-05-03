package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

const usage = `
usage:
	concur <data-dir-path>
`

func processFile(filePath string, ch chan int) {
	//TODO: open the file, scan each line,
	//do something with the word, and write
	//the results to the channel
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	n := 0
	for scanner.Scan() {
		n++

		for i := 0; i < 100; i++ {
			h := sha256.New()
			h.Write(scanner.Bytes())
			_ = h.Sum(nil)
		}
	}
	f.Close()
	ch <- n
}

func processDir(dirPath string) {
	//TODO: iterate over the files in the directory
	//and process each, first in a serial manner,
	//and then in a concurrent manner
	fileinfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan int, len(fileinfos))
	for _, fi := range fileinfos {
		go processFile(path.Join(dirPath, fi.Name()), ch) // Compare concurrent to serial
	}
	nWords := 0
	for i := 0; i < len(fileinfos); i++ {
		nWords += <-ch
	}
	fmt.Printf("processed %d words\n", nWords)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(usage)
		os.Exit(1)
	}

	dir := os.Args[1]

	fmt.Printf("processing directory %s...\n", dir)
	start := time.Now()
	processDir(dir)
	fmt.Printf("completed in %v\n", time.Since(start))
}
