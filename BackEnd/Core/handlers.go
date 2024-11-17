package backend

import (
	"fmt"
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

func RequireSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := GetSession(r)
		if err != nil || sess == nil {
			log.Printf("Unauthorized access attempt: %v", err)
			http.Redirect(w, r, "/choose", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func ChooseHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "choose", nil)
}

func EnginesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "engines", nil)
}

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	topScores, err := GetTopScores(10) // Get top 10 scores
	if err != nil {
		log.Printf("Error loading scores: %v", err)
		topScores = []ScoreEntry{}
	}

	data := map[string]interface{}{
		"Scores": topScores,
	}
	renderTemplate(w, "scoreboard", data)
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := GetSession(r) // Error checking is done by middleware

	// If game is over, redirect to result page
	if sess.IsGameOver {
		http.Redirect(w, r, "/result", http.StatusSeeOther)
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

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := GetSession(r)
	if err != nil {
		log.Printf("Session error in ResultHandler: %v", err)
		http.Redirect(w, r, "/choose", http.StatusSeeOther)
		return
	}

	// Additional validation for the result page
	if !sess.IsGameOver {
		log.Printf("Unauthorized access to result page: game not over")
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	// Save the score when game is over and player has won
	if sess.IsGameOver && sess.HasWon {
		err := SaveScore(sess.PlayerName, sess.Score)
		if err != nil {
			log.Printf("Error saving score: %v", err)
		}
	}

	data := map[string]interface{}{
		"PlayerName":      sess.PlayerName,
		"Score":           sess.Score,
		"GameOverMessage": getGameOverMessage(sess.HasWon),
		"WordToGuess":     sess.WordToGuess,
	}

	renderTemplate(w, "result", data)
}

// Add a new middleware specifically for the result page
func RequireGameOverSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := GetSession(r)
		if err != nil || sess == nil {
			log.Printf("Session error in result page middleware: %v", err)
			http.Redirect(w, r, "/choose", http.StatusSeeOther)
			return
		}

		// Check if the game is actually over
		if !sess.IsGameOver {
			log.Printf("Unauthorized access: game not over")
			http.Redirect(w, r, "/game", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
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

		// Add validation
		if playerName == "" {
			http.Error(w, "Player name is required", http.StatusBadRequest)
			return
		}

		if difficulty == "" {
			http.Error(w, "Difficulty selection is required", http.StatusBadRequest)
			return
		}

		err := CreateNewSession(w, r, playerName, difficulty)
		if err != nil {
			log.Printf("Session creation error: %v", err)
			http.Error(w, fmt.Sprintf("Error starting game: %v", err), http.StatusInternalServerError)
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
		log.Printf("Session error: %v", err)
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
			// Create new session with all data preserved
			newSess := &Session{
				PlayerName:     sess.PlayerName, // Preserve player name
				Difficulty:     sess.Difficulty, // Preserve difficulty
				Score:          sess.Score,
				Attempts:       game.AttemptsLeft,
				WordToGuess:    game.WordToGuess,
				GuessedLetters: strings.Join(game.GuessedLetters, ","),
				IsGameOver:     game.IsOver,
				HasWon:         game.HasWon,
			}

			if game.HasWon {
				newSess.Score += calculateScore(sess.Difficulty, game.AttemptsLeft)
			}

			err = SaveSession(w, r, newSess)
			if err != nil {
				log.Printf("Error saving session: %v", err)
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
