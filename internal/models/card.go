package models

import (
	"fmt"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Card struct {
	ID       string      `json:"id"`
	Deck     string      `json:"-"`
	Version  int         `json:"version"`
	Active   bool        `json:"active"`
	Modified bool        `json:"-"`
	Question string      `json:"question"`
	Answer   string      `json:"answer"`
	Reviews  ReviewSlice `json:"reviews"`
}

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

func NewCard(question string, answer string, deck string) *Card {
	return &Card{
		ID:       randomString(10),
		Deck:     deck,
		Version:  0,
		Question: question,
		Answer:   answer,
		Active:   true,
		Modified: true,
		Reviews:  ReviewSlice{},
	}
}

func (card *Card) Copy() *Card {
	newCard := *card
	newReviews := make(ReviewSlice, len(card.Reviews))
	copy(newReviews, card.Reviews)
	newCard.Reviews = newReviews
	return &newCard
}

func (card *Card) String() string {
	return fmt.Sprintf("%s\n%s%s\n", card.Question, tempFileDivider, card.Answer)
}
