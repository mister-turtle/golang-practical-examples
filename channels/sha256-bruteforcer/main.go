package main

import (
	"bufio"
	"bytes"
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

	// here we use a wait group to track how many words have been added and read from the channel
	// we also use an additional Add here, which corresponds to the wg.Done() after the wordlist has finished reading.
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// spawn the goroutines for reading the wordlist, and hashing the words
	go readWordList(wordChan, wordFile, wg)
	for i := 0; i <= *argThreads; i++ {
		go hashAndCompare(wordChan, *argHash, wg)
	}

	wg.Wait()
	log.Println("Finished cracking.")
}

func readWordList(wordChan chan string, words io.Reader, wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(words)
	for scanner.Scan() {
		wg.Add(1)
		wordChan <- scanner.Text()
	}
	close(wordChan)
	wg.Done()
}

func hashAndCompare(wordChan chan string, target string, wg *sync.WaitGroup) {

	// firstly take the hash byte string and convert it into a byteslice
	targetBytes, err := hex.DecodeString(target)
	if err != nil {
		log.Fatalf("failed to get hex bytes from target hash: %s", err.Error())
	}

	for word := range wordChan {
		h := sha256.New()
		h.Write([]byte(word))
		if bytes.Equal(h.Sum(nil), targetBytes) {
			log.Printf("FOUND! %s was %x\n", word, h.Sum(nil))
		}
		log.Printf("Received word %s\n", word)
		wg.Done()
	}
}
