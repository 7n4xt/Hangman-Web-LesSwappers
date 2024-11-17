package backend

import (
	"hangmanWeb/BackEnd/utils"
	"strings"
)

func NewGame(difficulty string) *Game {
	word := utils.GetWord(difficulty)
	return &Game{
		WordToGuess:    strings.ToLower(word),
		GuessedLetters: []string{},
		AttemptsLeft:   6,
		IsOver:         false,
		HasWon:         false,
	}
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
