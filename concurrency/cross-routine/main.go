package main

import (
	"log"
	"math/rand"
	"time"
)

// IntHandler takes two channels of type integer. It will receive integers from ichan, and if they are larger than 500, pass them on to aboveChan
func IntHandler(ichan chan int, aboveChan chan int) {
	log.Println("Int Handler Started")

	// range over the incoming integer channel to receive values from it
	for i := range ichan {
		log.Printf("Int Handler received: %d\n", i)

		// if the value received is over 500, pass it on to the aboveChan
		if i > 500 {
			aboveChan <- i
		}
	}
}

// Above500Handler takes only a single channel, as it only expects to receive values over 500 from the IntHandler.
func Above500Handler(aboveChan chan int) {
	log.Println("Above 500 Handler Started")
	for i := range aboveChan {

		// although we expect IntHandler to have filtered these values already, its sensible to do our own validation.
		// continue here will jump to the next iteration of the for loop.
		if i <= 500 {
			log.Printf("Error: Above 500 received %d (expected > 500)\n", i)
			continue
		}
		log.Printf("Above 500 received: %d\n", i)
	}
}

func main() {

	// Remove the time and date stamp from the Log output
	log.SetFlags(0)
	log.Println("Go Channels Demo")

	// Create two channels of type int. These are unbuffered channels.
	intChan := make(chan int)
	aboveChan := make(chan int)

	// Start our two go routines running the handlers.
	go IntHandler(intChan, aboveChan)
	go Above500Handler(aboveChan)

	// In this example, we will just populate the intChan channel with data from main()
	// as the channels are unbuffered, they will only take one value at a time.
	// this means the channel needs to be empty or sending to the channel will block.
	for i := 0; i < 10; i++ {
		log.Println("\n -----Sending new values-----")

		// we sleep here just so it doesn't complete instantly
		time.Sleep(1 * time.Second)

		intChan <- rand.Intn(600-1) + 1
	}

	// Close our chanels now we're done with them
	close(aboveChan)
	close(intChan)

	log.Println("Exiting gracefully.")
}
