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

func gameHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "game", nil)
}

func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "scoreboard", nil)
}

func main() {
	http.HandleFunc("/index", homeHandler)
	http.HandleFunc("/engines", enginesHandler)
	http.HandleFunc("/choose", chooseHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/scoreboard", scoreboardHandler)

	fileServer := http.FileServer(http.Dir("./Style"))
	http.Handle("/Style/", http.StripPrefix("/Style/", fileServer))

	var err error = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Server starting on port 8080...\n")
	}
}
