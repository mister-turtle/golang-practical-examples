package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	argPort := flag.Int("p", 8000, "TCP port to listen on")
	argAddress := flag.String("l", "127.0.0.1", "IP address to listen on")

	flag.Parse()

	// set the log output to stdout
	log.SetOutput(os.Stdout)

	// this example uses the net/http router as it doesnt require any advanced routing
	// HandleFunc takes a string for the route, and a handler
	// a handler is defined as a function with the signature func (http.ResponseWriter, *http.Request)
	http.HandleFunc("/", Index)
	http.HandleFunc("/redirect", Redirect)
	http.HandleFunc("/cookie", Cookie)
	http.HandleFunc("/headers", Headers)

	// following convention and deining the variable for address close to where its used
	address := fmt.Sprintf("%s:%d", *argAddress, *argPort)
	log.Printf("Listening on %s\n", address)

	// start the server listening on our address, but don't pass in an additional router.
	log.Fatal(http.ListenAndServe(address, nil))

}

func Index(w http.ResponseWriter, r *http.Request) {

	// we can use the response writer directly and pass it []bytes to write.
	// it will write the headers at this point, and then our data as the body
	w.Write([]byte("Hello, world!"))
}

func Redirect(w http.ResponseWriter, r *http.Request) {

	// we can also pass the writer and request on to other functions
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func Cookie(w http.ResponseWriter, r *http.Request) {

	// here is another example where we set cookies using the stdlib http package
	// note we have to do this before we've written data to the http.ResponseWriter
	newCookie := http.Cookie{
		Name:  "MyCookie",
		Value: "123456789",
	}

	http.SetCookie(w, &newCookie)
	w.Write([]byte("Check the cookie in the response!"))
}

func Headers(w http.ResponseWriter, r *http.Request) {

	// adding headers is pretty straightforward too
	w.Header().Add("custom-header", "super useful value")
	w.Write([]byte("Check the headers!"))
}
