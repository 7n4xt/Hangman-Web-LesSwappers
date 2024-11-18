package backend

type Game struct {
	WordToGuess    string
	GuessedLetters []string
	AttemptsLeft   int
	IsOver         bool
	HasWon         bool
}

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

type ScoreEntry struct {
	PlayerName string `json:"playerName"`
	Score      int    `json:"score"`
	Difficulty string `json:"difficulty"`
}
