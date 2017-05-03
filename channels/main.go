package main

import (
	"fmt"
	"math/rand"
	"time"
)

//someLongFunc is a function that might
//take a while to complete, so we want
//to run it on its own go routine
func someLongFunc(ch chan int) {
	r := rand.Intn(2000)
	d := time.Duration(r)
	time.Sleep(time.Millisecond * d)
	ch <- r // Writing from the channel
}

func main() {
	//TODO:
	//create a channel and call
	//someLongFunc() on a go routine
	//passing the channel so that
	//someLongFunc() can communicate
	//its results
	rand.Seed(time.Now().UnixNano())
	fmt.Println("starting long-running func...")
	n := 10
	ch := make(chan int, n)
	start := time.Now()
	for i := 0; i < n; i++ {
		go someLongFunc(ch) // go means run concurrently
	}
	for i := 0; i < n; i++ {
		result := <-ch // Reading from the channel, this code yeilds when until channel is written to
		fmt.Printf("result was %d\n", result)
	}
	fmt.Printf("took %v\n", time.Since(start))
}
