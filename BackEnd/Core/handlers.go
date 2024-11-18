package backend

import (
	"fmt"
	"hangmanWeb/BackEnd/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
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
	topScores, err := GetTopScores(10)
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
	sess, _ := GetSession(r)

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
	elapsedTime := time.Since(sess.StartTime).Seconds()

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
		"ElapsedTime":     int(elapsedTime),
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
	if !sess.IsGameOver {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}
	baseScore := calculateBaseScore(sess.Difficulty, sess.Attempts)
	timeBonus := 0
	errorBonus := 0

	if sess.HasWon {
		elapsedSeconds := int(time.Since(sess.StartTime).Seconds())
		timeBonus = calculateTimeBonus(elapsedSeconds, sess.Difficulty)
		errorBonus = calculateErrorScore(sess.TotalGuesses, sess.WrongGuesses, sess.Difficulty)
		finalScore := baseScore + timeBonus + errorBonus
		err = SaveScore(sess.PlayerName, finalScore, sess.Difficulty)
		if err != nil {
			log.Printf("Error saving score: %v", err)
		}
	}
	accuracy := 0.0
	if sess.TotalGuesses > 0 {
		accuracy = float64(sess.TotalGuesses-sess.WrongGuesses) / float64(sess.TotalGuesses) * 100
	}

	data := map[string]interface{}{
		"PlayerName":      sess.PlayerName,
		"BaseScore":       baseScore,
		"TimeBonus":       timeBonus,
		"ErrorBonus":      errorBonus,
		"TotalScore":      baseScore + timeBonus + errorBonus,
		"TotalGuesses":    sess.TotalGuesses,
		"WrongGuesses":    sess.WrongGuesses,
		"Accuracy":        accuracy,
		"GameOverMessage": getGameOverMessage(sess.HasWon),
		"WordToGuess":     sess.WordToGuess,
		"HasWon":          sess.HasWon,
		"Attempts":        sess.Attempts,
		"Difficulty":      sess.Difficulty,
	}
	renderTemplate(w, "result", data)
}

func RequireGameOverSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := GetSession(r)
		if err != nil || sess == nil {
			log.Printf("Session error in result page middleware: %v", err)
			http.Redirect(w, r, "/choose", http.StatusSeeOther)
			return
		}

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
		return utils.GetRandomPhrase("BackEnd/utils/won.txt")
	}
	return utils.GetRandomPhrase("BackEnd/utils/lose.txt")
}

func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		playerName := r.FormValue("pseudo")
		difficulty := r.FormValue("difficulty")

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
		if guess == "" {
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

		sess.TotalGuesses++
		isCorrectGuess := false

		if len(guess) > 1 {
			if guess == game.WordToGuess {
				game.HasWon = true
				game.IsOver = true
				isCorrectGuess = true
			} else {
				game.AttemptsLeft -= 2
				sess.WrongGuesses++
				if game.AttemptsLeft <= 0 {
					game.AttemptsLeft = 0
					game.IsOver = true
				}
			}
		} else {
			isCorrectGuess = strings.Contains(game.WordToGuess, guess)
			if !isCorrectGuess {
				sess.WrongGuesses++
			}
			game.GuessLetter(guess)
		}
		baseScore := 0
		timeBonus := 0
		errorBonus := 0

		if game.HasWon {
			elapsedSeconds := int(time.Since(sess.StartTime).Seconds())
			baseScore = calculateBaseScore(sess.Difficulty, game.AttemptsLeft)
			timeBonus = calculateTimeBonus(elapsedSeconds, sess.Difficulty)
			errorBonus = calculateErrorScore(sess.TotalGuesses, sess.WrongGuesses, sess.Difficulty)
		}

		newSess := &Session{
			PlayerName:     sess.PlayerName,
			Difficulty:     sess.Difficulty,
			Score:          baseScore + timeBonus + errorBonus,
			TimeBonus:      timeBonus,
			ErrorBonus:     errorBonus,
			TotalGuesses:   sess.TotalGuesses,
			WrongGuesses:   sess.WrongGuesses,
			Attempts:       game.AttemptsLeft,
			WordToGuess:    game.WordToGuess,
			GuessedLetters: strings.Join(game.GuessedLetters, ","),
			IsGameOver:     game.IsOver,
			HasWon:         game.HasWon,
			StartTime:      sess.StartTime,
		}

		err = SaveSession(w, r, newSess)
		if err != nil {
			log.Printf("Error saving session: %v", err)
			http.Error(w, "Error saving session", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func calculateBaseScore(difficulty string, attemptsLeft int) int {
	difficultyMultiplier := map[string]int{
		"Easy":   1,
		"Normal": 2,
		"Hard":   3,
		"Insane": 4,
	}
	return 10 * difficultyMultiplier[difficulty] * attemptsLeft
}

func calculateTimeBonus(elapsedSeconds int, difficulty string) int {
	difficultyMultiplier := map[string]int{
		"Easy":   1,
		"Normal": 2,
		"Hard":   3,
		"Insane": 4,
	}

	var bonus int
	switch {
	case elapsedSeconds <= 10:
		bonus = 200
	case elapsedSeconds <= 20:
		bonus = 150
	case elapsedSeconds <= 30:
		bonus = 100
	case elapsedSeconds <= 45:
		bonus = 50
	case elapsedSeconds <= 60:
		bonus = 25
	default:
		bonus = 10
	}

	return bonus * difficultyMultiplier[difficulty]
}

func calculateErrorScore(totalGuesses int, wrongGuesses int, difficulty string) int {
	difficultyMultiplier := map[string]int{
		"Easy":   1,
		"Normal": 2,
		"Hard":   3,
		"Insane": 4,
	}
	correctGuesses := totalGuesses - wrongGuesses
	accuracy := float64(correctGuesses) / float64(totalGuesses)

	var bonus int
	switch {
	case accuracy >= 0.9:
		bonus = 200
	case accuracy >= 0.8:
		bonus = 150
	case accuracy >= 0.7:
		bonus = 100
	case accuracy >= 0.6:
		bonus = 75
	case accuracy >= 0.5:
		bonus = 50
	case accuracy >= 0.4:
		bonus = 25
	case accuracy >= 0.3:
		bonus = 10
	default:
		bonus = 0
	}

	return bonus * difficultyMultiplier[difficulty]
}
