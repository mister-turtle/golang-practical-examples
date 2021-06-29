package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

const indexTemplate = `
<html>
<head>
	<title>Welcome, {{ .IP }}</title>
</head>
<body>
	<h1>Hello, {{ .IP }}</h1>
	<br>
	<p>The current server time is {{ .Time }}
</body>
</html>
`

func main() {

	log.SetOutput(os.Stdout)

	// define some useful command line arguments
	argPort := flag.Int("p", 8000, "TCP port to listen on")
	argAddress := flag.String("l", "127.0.0.1", "IP address to listen on")
	flag.Parse()

	// parse our template from the const defined at the top of program
	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// unlike the simplest-server example, we're defining our function in the call to HandleFunc
	// this is so we can use a closure, allowing us to access the tmpl defined above within the handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// data in the template is accessed via a map, so lets put some data in
		data := map[string]string{
			"IP":   r.RemoteAddr,
			"Time": time.Now().String(),
		}

		// templates have an Execute method that take an io.Writer, an http.ResponseWriter satisfies the io.Writer interface
		// so we can pass in the writer directly to the execute function
		err = tmpl.Execute(w, data)
		if err != nil {

			// it is bad practice to send internal errors out to clients
			// even in test code lets follow that
			// we log the error locally but send out a generic error message to the client
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong..."))
			log.Printf("ERROR: %s\n", err.Error())

			// make sure to return here otherwise execution will continue
			// although no further logic exists in this function, we could add some in the future
			return
		}
	})

	// assemble the address from the command line args and start our server
	address := fmt.Sprintf("%s:%d", *argAddress, *argPort)
	log.Fatal(http.ListenAndServe(address, nil))
}
