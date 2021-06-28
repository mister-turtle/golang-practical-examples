package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"os"
	"sync"
)

func main() {

	// Remove the time and date from logging output and use stdout instead of stderr
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	log.Println("Simple Hash Brute Forcer")

	// define the command line arguments for the utility.
	argWordlist := flag.String("w", "", "Newline delimited wordlist")
	argHash := flag.String("h", "", "SHA256 hash to brute force")
	argThreads := flag.Int("t", 4, "Number of goroutines")

	// parse the given arguments and check for mandatory parameters
	flag.Parse()

	if *argWordlist == "" || *argHash == "" {
		// flag.Usage will print out generated cli usage
		flag.Usage()
		os.Exit(1)
	}

	// verify we can open the specified word file.
	wordFile, err := os.Open(*argWordlist)
	if err != nil {
		log.Fatalf("failed to open %s: %s\n", *argWordlist, err.Error())
	}

	// create a channel of type string that we can use to send words read from the file
	wordChan := make(chan string)

	// create a wait group that we can use to control the flow of execution
	wg := &sync.WaitGroup{}

	// finally, create a context that we can use to stop goroutines running.
	// this will be used if we find the matching word for the hash in the hashlist, we can gracefully terminate the goroutines.
	ctx, cancel := context.WithCancel(context.Background())

	// firstly we spawn a gortouine that will read the word file line by line.
	// we add one to the waitgroup for each goroutine that we spawn.
	// we remove one from the waitgroup when that go routine finishes.
	// the go routine takes the context, the word channel, and the word file.
	wg.Add(1)
	go func() {
		readWordList(ctx, wordChan, wordFile)
		wg.Done()
	}()

	// create a channel that will be used by the hashing go routines if they find the correct word.
	found := make(chan string)

	// create multiple go routines (determined by the -t flag) to do the hashing
	// again we add one to the waitgroup for each go routine, and remove one when the go routine exits.
	// the hashing go routine takes the context, the worc channel, the hash to compare to, and the found channel
	for i := 0; i <= *argThreads; i++ {
		wg.Add(1)
		go func() {
			hashAndCompare(ctx, wordChan, *argHash, found)
			wg.Done()
		}()
	}

	// we spawn another goroutine that receives items from the found channel
	// when an item is received, it will print to stdout that the hash has been cracked, and print the result.
	// it then calls cancel(), which is from the context above to alert the other go routines that they can stop.
	go func() {
		result := <-found
		log.Printf("Cracked: %s\n", result)
		cancel()
	}()

	// meanwhile, in the main thread we simply wait for the waitgroup to be empty.
	// the word list reader will close the word channel when its done, meaning the go routines will also close when there are no more items on the word channel.
	// each of the wg.Done() calls will be made when the go routines exit, and the wait here will unblock.
	// if a result is found, the call to cancel() will stop the go routines, unblocking us here as well.
	wg.Wait()
	log.Println("Finished cracking.")
}

func readWordList(ctx context.Context, wordChan chan string, words io.Reader) {

	// make sure we close the channel when this go routine exits.
	defer func() {
		close(wordChan)
	}()

	// create a new scanner from the file passed in.
	// an *os.File implements the io.Reader interface so we can pass it in directly.
	scanner := bufio.NewScanner(words)

	// this is a loop tag, used so we can break out of a nested loop.
scanloop:
	for {
		if !scanner.Scan() {

			// the scanner will return false if there is an error, including an EOF.
			// however, if it's an EOF error, the scanner.Err() will not be set.
			// this is defined in the documentation for a bufio.Scanner
			if err := scanner.Err(); err != nil {
				log.Fatalf("Failed to read word file: %s\n", err.Error())
			}
			// this is likely EOF (end of file)
			break scanloop
		}

		// select will choose between reading from scanner.Text() and placing on the word channel, or if available the context.Done.
		// context.Done will be usable if we call the cancel() function associated with it.
		// this way if a result is found, and cancel() called from main, we will stop reading words from the word file.
		select {
		case wordChan <- scanner.Text():
		case <-ctx.Done():
			break scanloop
		}
	}
}

func hashAndCompare(ctx context.Context, wordChan chan string, target string, found chan string) {

	// firstly take the hash byte string and convert it into a byteslice
	targetBytes, err := hex.DecodeString(target)
	if err != nil {
		log.Fatalf("failed to get hex bytes from target hash: %s", err.Error())
	}

	// initalise a new hash container and utilise Reset() later to avoid a new instance each time.
	h := sha256.New()

	// iterate over the channel until it is empty. If the channel is closed, the range will finish when there are no more items available.
	// this is the same method as the readWordList function
forchan:
	for {
		select {

		// reading two items off the channel will yield a value, and a bool related to the status of the channel (open:true, closed:false)
		case word, ok := <-wordChan:

			// if !ok, eg invert the ok boolean, then the channel was closed
			if !ok {
				break forchan
			}

			// reset the hash container and write the bytes to it.
			h.Reset()
			h.Write([]byte(word))

			// compare the hash bytes of our current word to the target bytes.
			// if found send to the found channel
			if bytes.Equal(h.Sum(nil), targetBytes) {
				found <- word
			}

		case <-ctx.Done():
			break forchan
		}
	}
}
