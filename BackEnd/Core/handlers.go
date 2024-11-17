package backend

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates *template.Template

func InitTemplates() {
	var err error
	templates, err = template.ParseGlob("Templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func ChooseHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "choose", nil)
}

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "scoreboard", nil)
}

func EnginesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "engines", nil)
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var guessedLetters []string
	if sess.GuessedLetters != "" {
		guessedLetters = strings.Split(sess.GuessedLetters, ",")
	} else {
		guessedLetters = []string{}
	}

	game := &Game{
		WordToGuess:    sess.WordToGuess,
		GuessedLetters: guessedLetters,
		AttemptsLeft:   sess.Attempts,
		IsOver:         sess.IsGameOver,
		HasWon:         sess.HasWon,
	}

	displayWord := game.GetDisplayWord()

	data := map[string]interface{}{
		"PlayerName":      sess.PlayerName,
		"Score":           sess.Score,
		"Attempts":        sess.Attempts,
		"Difficulty":      sess.Difficulty,
		"DisplayWord":     displayWord,
		"GuessedLetters":  guessedLetters,
		"GameOver":        sess.IsGameOver,
		"GameOverMessage": getGameOverMessage(sess.HasWon),
		"WordToGuess":     sess.WordToGuess,
	}

	renderTemplate(w, "game", data)
}

func getGameOverMessage(hasWon bool) string {
	if hasWon {
		return "Congratulations! You won!"
	}
	return "Game Over! Better luck next time!"
}

func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		playerName := r.FormValue("pseudo")
		difficulty := r.FormValue("difficulty")

		err := CreateNewSession(w, r, playerName, difficulty)
		if err != nil {
			http.Error(w, "Error starting game", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/choose", http.StatusSeeOther)
}

func GuessHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if sess.IsGameOver {
			http.Redirect(w, r, "/game", http.StatusSeeOther)
			return
		}

		guess := strings.ToLower(r.FormValue("guess"))
		if len(guess) != 1 {
			http.Redirect(w, r, "/game", http.StatusSeeOther)
			return
		}

		// Fix the GuessedLetters split
		var guessedLetters []string
		if sess.GuessedLetters != "" {
			guessedLetters = strings.Split(sess.GuessedLetters, ",")
		} else {
			guessedLetters = []string{}
		}

		game := &Game{
			WordToGuess:    sess.WordToGuess,
			GuessedLetters: guessedLetters,
			AttemptsLeft:   sess.Attempts,
			IsOver:         sess.IsGameOver,
			HasWon:         sess.HasWon,
		}

		if game.GuessLetter(guess) {
			// Update session with new game state
			sess.GuessedLetters = strings.Join(game.GuessedLetters, ",")
			sess.Attempts = game.AttemptsLeft
			sess.IsGameOver = game.IsOver
			sess.HasWon = game.HasWon

			if game.HasWon {
				sess.Score += calculateScore(sess.Difficulty, game.AttemptsLeft)
			}

			err = SaveSession(w, r, sess)
			if err != nil {
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func calculateScore(difficulty string, attemptsLeft int) int {
	baseScore := 10
	difficultyMultiplier := map[string]int{
		"Easy":   1,
		"Normal": 2,
		"Hard":   3,
		"Insane": 4,
	}

	return baseScore * difficultyMultiplier[difficulty] * attemptsLeft
}
