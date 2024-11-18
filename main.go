package main

import (
	"fmt"
	"hangmanWeb/BackEnd/Core"
	"log"
	"net/http"
)

func main() {
	backend.InitTemplates()
	fileServer := http.FileServer(http.Dir("./Style"))
	http.Handle("/Style/", http.StripPrefix("/Style/", fileServer))
	http.HandleFunc("/", backend.IndexHandler)
	http.HandleFunc("/result", backend.RequireSession(backend.RequireGameOverSession(backend.ResultHandler)))
	http.HandleFunc("/game", backend.RequireSession(backend.GameHandler))
	http.HandleFunc("/start-game", backend.StartGameHandler)
	http.HandleFunc("/guess", backend.GuessHandler)
	http.HandleFunc("/index", backend.IndexHandler)
	http.HandleFunc("/choose", backend.ChooseHandler)
	http.HandleFunc("/scoreboard", backend.ScoreboardHandler)
	http.HandleFunc("/engines", backend.EnginesHandler)
	fmt.Printf("Server starting on port 8080...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
