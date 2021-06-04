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

	argPortStart := flag.Int("s", 1, "Starting port")
	argPortEnd := flag.Int("e", 1000, "Ending port")
	argTarget := flag.String("i", "", "Target IP address")
	argThreads := flag.Int("t", 5, "Number of threads to run")

	flag.Parse()

	if *argTarget == "" {
		flag.Usage()
		os.Exit(1)
	}

	portChan := make(chan Port)
	go func() {
		for i := *argPortStart; i <= *argPortEnd; i++ {
			portChan <- Port{
				Target: *argTarget,
				PortId: i,
			}
		}
		close(portChan)
	}()

	var wg sync.WaitGroup
	resChan := make(chan Port)
	log.Printf("Starting %d threads\n", *argThreads)
	for i := 1; i <= *argThreads; i++ {
		wg.Add(1)
		go ScannerThread(i, &wg, portChan, resChan)
	}

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for res := range resChan {
		log.Printf("%-16s : %-6d - %s\n", res.Target, res.PortId, res.Status)
	}

	log.Println("Finished.")
}

func ScannerThread(num int, wg *sync.WaitGroup, ports chan Port, res chan Port) {

	defer wg.Done()

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

