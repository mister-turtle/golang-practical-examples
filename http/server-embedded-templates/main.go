package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed embedded/*.html
var templateFS embed.FS

//go:embed embedded/*.css
//go:embed embedded/*.js
var staticFS embed.FS

type Server struct {
	Address   string
	router    chi.Router
	templates *template.Template
}

func (s Server) Start() error {
	return http.ListenAndServe(s.Address, s.router)
}

func (s Server) GetIndex(w http.ResponseWriter, r *http.Request) {
	err := s.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func NewServer(addr string) (Server, error) {

	newServer := Server{
		Address: addr,
	}

	templates, err := template.ParseFS(templateFS, "embedded/*.html")
	if err != nil {
		return newServer, err
	}
	newServer.templates = templates

	staticSubFS, err := fs.Sub(staticFS, "embedded")
	if err != nil {
		return newServer, err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", newServer.GetIndex)
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS))))
	newServer.router = router

	return newServer, nil
}
func main() {

	argAddress := flag.String("a", "127.0.0.1", "IP address to bind to")
	argPort := flag.String("p", "8000", "Port to listen on")
	flag.Parse()

	server, err := NewServer(fmt.Sprintf("%s:%s", *argAddress, *argPort))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())

}
