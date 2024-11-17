package utils

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

func GetWord(difficulty string) string {
	var filename string

	switch difficulty {
	case "Easy":
		filename = "easy.txt"
	case "Normal":
		filename = "normal.txt"
	case "Hard":
		filename = "hard.txt"
	case "Insane":
		filename = "insane.txt"
	default:
		filename = "easy.txt"
	}

	words, err := loadWordsFromFile(filename)
	if err != nil {
		return "hangman"
	}

	rand.Seed(time.Now().UnixNano())
	return words[rand.Intn(len(words))]
}

func loadWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, strings.ToLower(word))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
