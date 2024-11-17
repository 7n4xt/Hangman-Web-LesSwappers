package backend

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
)

var (
	secretKey = []byte("secret-key")
)

func init() {
	gob.Register(Session{})
}

func CreateNewSession(w http.ResponseWriter, r *http.Request, playerName, difficulty string) error {
	if playerName == "" {
		return fmt.Errorf("player name cannot be empty")
	}

	game := NewGame(difficulty)
	sess := &Session{
		PlayerName:     playerName,
		Difficulty:     difficulty,
		Score:          0,
		Attempts:       6,
		WordToGuess:    game.WordToGuess,
		GuessedLetters: "",
		IsGameOver:     false,
		HasWon:         false,
	}

	return SaveSession(w, r, sess)
}

func GetSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("hangman-session")
	if err != nil {
		return nil, err
	}

	// Decode base64 string
	decoded, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	var sess Session
	err = json.Unmarshal(decoded, &sess)
	if err != nil {
		return nil, err
	}

	// Validate session data
	if sess.PlayerName == "" {
		return nil, fmt.Errorf("invalid session: player name is empty")
	}

	return &sess, nil
}

func SaveSession(w http.ResponseWriter, r *http.Request, sess *Session) error {
	sessionData, err := json.Marshal(sess)
	if err != nil {
		return err
	}

	// Encode to base64
	encoded := base64.URLEncoding.EncodeToString(sessionData)

	cookie := &http.Cookie{
		Name:     "hangman-session",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600,
	}

	http.SetCookie(w, cookie)
	return nil
}

func ClearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "hangman-session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

const scoresFile = "scores.json"

func SaveScore(playerName string, score int) error {
	// Read existing scores
	scores, err := LoadScores()
	if err != nil {
		scores = []ScoreEntry{}
	}

	// Add new score
	newScore := ScoreEntry{
		PlayerName: playerName,
		Score:      score,
	}
	scores = append(scores, newScore)

	// Sort scores in descending order
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Write to JSON file
	file, err := json.MarshalIndent(scores, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(scoresFile, file, 0644)
}

func LoadScores() ([]ScoreEntry, error) {
	var scores []ScoreEntry

	// Check if file exists
	if _, err := os.Stat(scoresFile); os.IsNotExist(err) {
		return []ScoreEntry{}, nil
	}

	// Read file
	data, err := os.ReadFile(scoresFile)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	err = json.Unmarshal(data, &scores)
	if err != nil {
		return nil, err
	}

	return scores, nil
}

func GetTopScores(limit int) ([]ScoreEntry, error) {
	scores, err := LoadScores()
	if err != nil {
		return nil, err
	}

	if limit > len(scores) {
		limit = len(scores)
	}
	return scores[:limit], nil
}
