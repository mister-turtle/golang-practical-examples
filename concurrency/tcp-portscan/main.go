package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Port struct {
	Target string
	PortId int
	Status string
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	log.Println("Simple Concurrent TCP Port Scanner")

	// Set up the command line arguments for this example
	argPortStart := flag.Int("s", 1, "Starting port")
	argPortEnd := flag.Int("e", 1000, "Ending port")
	argTarget := flag.String("i", "", "Target IP address")
	argThreads := flag.Int("t", 5, "Number of threads to run")

	flag.Parse()

	if *argTarget == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Create a channel which will take one port item including the target and port number
	ports := make(chan Port)

	// Start populating this channel with the port range in a go routine, so that the rest of the application can continue.
	go func() {
		for i := *argPortStart; i <= *argPortEnd; i++ {
			ports <- Port{
				Target: *argTarget,
				PortId: i,
			}
		}
		// Once all the ports have been put onto the channel, close the channel so that the scanner routines can finish their for-loop and stop.
		close(ports)
	}()

	// Create a waitgroup to use in the scanning go routines. This will allow us to wait and close down the results channel later on. This will have the effect of
	// letting the main thread which is looping over the results channel exit gracefully.
	var wg sync.WaitGroup
	results := make(chan Port)
	log.Printf("Starting %d threads\n", *argThreads)
	for i := 1; i <= *argThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ScannerThread(ports, results)
		}()
	}

	// Spawn a go routine which will wait until all of the goroutines are closed (each one calls wg.Done() when it's finished) and then close the results channel.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Loop around the results channel priting out all the of the results. This range will end when the results is closed by the go routine directly above.
	for res := range results {
		log.Printf("%-16s : %-6d - %s\n", res.Target, res.PortId, res.Status)
	}

	log.Println("Finished.")
}

func ScannerThread(ports chan Port, res chan Port) {

	for port := range ports {
		target := fmt.Sprintf("%s:%d", port.Target, port.PortId)

		conn, err := net.DialTimeout("tcp", target, time.Second*5)
		if err != nil {
			continue
		}
		conn.Close()

		res <- Port{
			Target: port.Target,
			PortId: port.PortId,
			Status: "open",
		}
	}
}
