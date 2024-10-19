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
	http.ListenAndServe("localhost:8000", nil)
}
