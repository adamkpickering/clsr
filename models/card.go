package models

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

// Returns a string of length n that is comprised of random letters
// and numbers.
// From https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type Card struct {
	ID         string
	Version    int
	LastReview time.Time
	NextReview time.Time
	Question   string
	Answer     string
	Active     bool
}

func NewCard(question string, answer string) Card {
	// generate ID
	rand.Seed(time.Now().UnixNano())
	id := randomString(10)

	// build card
	return Card{
		ID:         id,
		Version:    0,
		LastReview: time.Time{},
		NextReview: time.Now(),
		Question:   question,
		Answer:     answer,
		Active:     true,
	}
}
