package backend

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"net/http"
)

var (
	secretKey = []byte("secret-key")
)

type Session struct {
	PlayerName     string
	Difficulty     string
	Score          int
	Attempts       int
	WordToGuess    string
	GuessedLetters string
	IsGameOver     bool
	HasWon         bool
}

func init() {
	gob.Register(Session{})
}

func CreateNewSession(w http.ResponseWriter, r *http.Request, playerName, difficulty string) error {
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
