package backend

import (
	"hangmanWeb/BackEnd/utils"
	"math/rand"
	"strings"
	"time"
)

func NewGame(difficulty string) *Game {
	word := utils.GetWord(difficulty)
	game := &Game{
		WordToGuess:    strings.ToLower(word),
		GuessedLetters: []string{},
		AttemptsLeft:   6,
		IsOver:         false,
		HasWon:         false,
	}

	// Get 2 random letters from the word
	initialLetters := getInitialLetters(game.WordToGuess, 2)
	game.GuessedLetters = initialLetters

	return game
}

// New function to get initial letters
func getInitialLetters(word string, count int) []string {
	// Create a slice of available positions
	letterPositions := make([]int, len(word))
	for i := range letterPositions {
		letterPositions[i] = i
	}

	// Shuffle the positions
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(letterPositions), func(i, j int) {
		letterPositions[i], letterPositions[j] = letterPositions[j], letterPositions[i]
	})

	// Get unique letters
	selectedLetters := make(map[string]bool)
	result := []string{}

	// Select random letters
	for _, pos := range letterPositions {
		letter := string(word[pos])
		if !selectedLetters[letter] {
			selectedLetters[letter] = true
			result = append(result, letter)
			if len(result) >= count {
				break
			}
		}
	}

	return result
}

func (g *Game) GetDisplayWord() []string {
	display := make([]string, len(g.WordToGuess))
	for i, letter := range g.WordToGuess {
		if containsLetter(g.GuessedLetters, string(letter)) {
			display[i] = string(letter)
		} else {
			display[i] = "_"
		}
	}
	return display
}

func (g *Game) IsWordGuessed() bool {
	for _, letter := range g.WordToGuess {
		if !containsLetter(g.GuessedLetters, string(letter)) {
			return false
		}
	}
	return true
}

func (g *Game) GuessLetter(letter string) bool {
	letter = strings.ToLower(letter)

	if containsLetter(g.GuessedLetters, letter) {
		return false
	}

	g.GuessedLetters = append(g.GuessedLetters, letter)

	if !strings.Contains(g.WordToGuess, letter) {
		g.AttemptsLeft--
	}

	g.HasWon = g.IsWordGuessed()
	g.IsOver = g.HasWon || g.AttemptsLeft <= 0

	return true
}

func containsLetter(slice []string, letter string) bool {
	for _, l := range slice {
		if l == letter {
			return true
		}
	}
	return false
}
