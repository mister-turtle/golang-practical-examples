package main

import (
	"embed"
	"html/template"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
)

// Embed all template files from the templates directory.
//
//go:embed templates/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

// Map to store compiled templates by page name.
var templatesMap map[string]*template.Template

// Function to load and compile templates.
func loadTemplates() error {
	templatesMap = make(map[string]*template.Template)

	// Get all template files in the templates directory.
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return err
	}

	// Iterate over each template file except the base template and styles.css.
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name == "base.html" || filepath.Ext(name) != ".html" {
			continue
		}

		newTemplate := template.Must(template.New(name).ParseFS(templatesFS, "templates/base.html", "templates/"+name))

		// Store the compiled template in the map.
		templatesMap[name] = newTemplate
	}

	return nil
}

// Function to render a template by page name.
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templatesMap[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	// Execute the "layout" template.
	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Load and compile all templates at startup.
	err := loadTemplates()
	if err != nil {
		panic(err)
	}

	// Serve static files (CSS) from the templates directory under /static/.
	staticFiles, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	_ = mime.AddExtensionType(".css", "text/css")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

	// Set up HTTP handlers for different routes.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "home.html", nil)
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "about.html", nil)
	})
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "services.html", nil)
	})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "contact.html", nil)
	})

	// Start the HTTP server.
	http.ListenAndServe(":8000", nil)
}
