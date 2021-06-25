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

	// validate the word file being used
	wordFile, err := os.Open(*argWordlist)
	if err != nil {
		log.Fatalf("failed to open %s: %s\n", *argWordlist, err.Error())
	}

	// create a channel that we can send words from the wordlist to.
	wordChan := make(chan string, *argThreads*2)

	// we use an Add here, which corresponds to the wg.Done() after the wordlist has finished reading.
	wg := &sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	// spawn the goroutines for reading the wordlist, and hashing the words
	// add one to the WaitGroup per goroutine running.
	// this is removed from the WaitGroup when the goroutine exits.
	go readWordList(wordChan, wordFile, wg, ctx)
	for i := 0; i <= *argThreads; i++ {
		wg.Add(1)
		go hashAndCompare(ctx, cancel, wordChan, *argHash, wg)
	}

	wg.Wait()
	log.Println("Finished cracking.")
}

func readWordList(wordChan chan string, words io.Reader, wg *sync.WaitGroup, ctx context.Context) {

	defer func() {
		wg.Done()
		close(wordChan)
	}()

	scanner := bufio.NewScanner(words)

scanloop:
	for {
		scanner.Scan()
		wordChan <- scanner.Text()
		select {
		case <-ctx.Done():
			break scanloop
		}
	}
}

func hashAndCompare(ctx context.Context, cancel context.CancelFunc, wordChan chan string, target string, wg *sync.WaitGroup) {

	// remove one from the WaitGroup when this goroutine finishes.
	defer wg.Done()

	// firstly take the hash byte string and convert it into a byteslice
	targetBytes, err := hex.DecodeString(target)
	if err != nil {
		log.Fatalf("failed to get hex bytes from target hash: %s", err.Error())
	}

	// iterate over the channel until it is empty. If the channel is closed, the range will finish when there are no more items available.
forchan:
	for {
		select {
		case word, ok := <-wordChan:
			if !ok {
				break forchan
			}
			h := sha256.New()
			h.Write([]byte(word))
			if bytes.Equal(h.Sum(nil), targetBytes) {
				cancel()
				log.Printf("FOUND! %s was %x\n", word, h.Sum(nil))
			}
		case <-ctx.Done():
			break forchan
		}
	}
}
