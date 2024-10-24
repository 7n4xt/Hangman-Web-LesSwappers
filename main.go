package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
)

func main() {
	temp, tempErr := template.ParseGlob("Templates/*.html")
	if tempErr != nil {
		fmt.Printf("Error loading templates: %s", tempErr.Error())
		os.Exit(2)
	}
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "Home", nil)
	})
	http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "rules", nil)
	})
	http.HandleFunc("/engines", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "engines", nil)
	})
	http.ListenAndServe("localhost:8080", nil)
}
