package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob("Templates/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index", nil)
}

func enginesHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "engines", nil)
}

func chooseHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "choose", nil)
}

func main() {
	http.HandleFunc("/index", homeHandler)
	http.HandleFunc("/engines", enginesHandler)
	http.HandleFunc("/choose", chooseHandler)

	fileServer := http.FileServer(http.Dir("./Style"))
	http.Handle("/Style/", http.StripPrefix("/Style/", fileServer))

	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
