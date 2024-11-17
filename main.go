package main

import (
	"fmt"
	"hangmanWeb/backend/core"
	"log"
	"net/http"
)

func main() {
	// Initialize templates first
	backend.InitTemplates()

	fileServer := http.FileServer(http.Dir("./Style"))
	http.Handle("/Style/", http.StripPrefix("/Style/", fileServer))

	http.HandleFunc("/", backend.IndexHandler)
	http.HandleFunc("/start-game", backend.StartGameHandler)
	http.HandleFunc("/guess", backend.GuessHandler)
	http.HandleFunc("/index", backend.IndexHandler)
	http.HandleFunc("/choose", backend.ChooseHandler)
	http.HandleFunc("/game", backend.GameHandler)
	http.HandleFunc("/scoreboard", backend.ScoreboardHandler)
	http.HandleFunc("/engines", backend.EnginesHandler)

	fmt.Printf("Server starting on port 8080...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
