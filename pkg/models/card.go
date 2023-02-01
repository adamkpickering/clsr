package models

import (
	"fmt"
	"github.com/adamkpickering/clsr/pkg/config"
	"math/rand"
	"time"
)

type parseState int

const (
	header parseState = iota
	question
	answer
)

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

type Card struct {
	ID       string    `json:"id"`
	Deck     string    `json:"-"`
	Version  int       `json:"version"`
	Active   bool      `json:"active"`
	Modified bool      `json:"-"`
	Question string    `json:"question"`
	Answer   string    `json:"answer"`
	Reviews  []*Review `json:"reviews"`
}

// Returns a string of length n that is comprised of random letters
// and numbers.
// From https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randomString(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
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
		Reviews:  []*Review{},
	}
}

func (card *Card) NextReview() (time.Time, error) {
	var timeBetweenReviews time.Duration
	reviewsLength := len(card.Reviews)
	if reviewsLength == 0 {
		return time.Now(), nil
	} else if reviewsLength == 1 {
		timeBetweenReviews = 24 * time.Hour
	} else {
		lastReview := card.Reviews[0].Datetime
		lastLastReview := card.Reviews[1].Datetime
		timeBetweenReviews = lastReview.Sub(lastLastReview)
	}

	multiplier, err := card.GetMultiplier(config.DefaultConfig)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get multiplier: %w", err)
	}

	timeUntilNextReview := float64(timeBetweenReviews) * multiplier
	lastReview := card.Reviews[0].Datetime
	return lastReview.Add(time.Duration(timeUntilNextReview)), nil
}

func (card *Card) GetMultiplier(config *config.Config) (float64, error) {
	if len(card.Reviews) == 0 {
		return 0.0, nil
	}
	result := card.Reviews[0].Result
	switch result {
	case Failed:
		return config.Multipliers.Failed, nil
	case Hard:
		return config.Multipliers.Hard, nil
	case Normal:
		return config.Multipliers.Normal, nil
	case Easy:
		return config.Multipliers.Easy, nil
	}
	return 0.0, fmt.Errorf("got unexpected result %q", result)
}

func (card *Card) IsDue() (bool, error) {
	nextReview, err := card.NextReview()
	if err != nil {
		return false, fmt.Errorf("failed to get next review: %w", err)
	}
	return time.Now().After(nextReview), nil
}

// // Returns the current review interval, that is, the number of days
// // between the date the card was last reviewed and the date the card
// // should be reviewed next. Review interval is in days.
// func (card *Card) GetCurrentReviewInterval() int {
// 	// if card has not been studied yet, return 1
// 	if card.LastReview.IsZero() {
// 		return 1
// 	}

// 	// determine interval if card has been studied
// 	rawDifference := card.NextReview.Sub(card.LastReview)
// 	difference := rawDifference.Round(24 * time.Hour)
// 	days := int(difference / (24 * time.Hour))

// 	// without this we can't ever get past 0 days
// 	if days == 0 {
// 		return 1
// 	}

// 	return days
// }

// // Given a multiplier that represents the effect on the review interval
// // following a card review, returns the number of days until that card
// // should be reviewed again.
// func (card *Card) GetMultipliedReviewInterval(multiplier float64) int {
// 	currentInterval := card.GetCurrentReviewInterval()
// 	newInterval := int(math.Round(float64(currentInterval) * multiplier))
// 	return newInterval
// }

// func (card *Card) DaysUntilDue() int {
// 	year, month, day := time.Now().Date()
// 	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
// 	timeUntilDue := card.NextReview.Sub(today)
// 	daysUntilDue := int(timeUntilDue / (24 * time.Hour))
// 	return daysUntilDue
// }
