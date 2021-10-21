package models

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type parseState int

const (
	header parseState = iota
	question
	answer
)

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

const dateLayout = "2006-01-02"

var idRegex *regexp.Regexp
var versionRegex *regexp.Regexp
var lastReviewRegex *regexp.Regexp
var nextReviewRegex *regexp.Regexp
var activeRegex *regexp.Regexp
var dividerRegex *regexp.Regexp
var commentRegex *regexp.Regexp

func init() {
	idRegex = regexp.MustCompile(`^ID *=`)
	versionRegex = regexp.MustCompile(`^Version *=`)
	lastReviewRegex = regexp.MustCompile(`^LastReview *=`)
	nextReviewRegex = regexp.MustCompile(`^NextReview *=`)
	activeRegex = regexp.MustCompile(`^Active *=`)
	dividerRegex = regexp.MustCompile(`^---`)
	commentRegex = regexp.MustCompile(`^ *#`)
}

type Card struct {
	ID         string
	Version    int
	LastReview time.Time
	NextReview time.Time
	Active     bool
	Question   string
	Answer     string
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

func NewCard(question string, answer string) *Card {
	year, month, day := time.Now().Date()
	return &Card{
		ID:         randomString(10),
		Version:    0,
		LastReview: time.Time{},
		NextReview: time.Date(year, month, day, 0, 0, 0, 0, time.Local),
		Question:   question,
		Answer:     answer,
		Active:     true,
	}
}

func (card *Card) MarshalText() ([]byte, error) {
	// process card into a map
	outputCard := map[string]string{}
	outputCard["ID"] = card.ID
	outputCard["Version"] = fmt.Sprintf("%d", card.Version)
	outputCard["LastReview"] = card.LastReview.Format(dateLayout)
	outputCard["NextReview"] = card.NextReview.Format(dateLayout)
	outputCard["Active"] = fmt.Sprintf("%t", card.Active)
	outputCard["Question"] = card.Question
	outputCard["Answer"] = card.Answer

	// fill buffer with output
	buffer := &bytes.Buffer{}
	err := CardTemplate.Execute(buffer, outputCard)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to execute CardTemplate: %w", err)
	}

	return buffer.Bytes(), nil
}

func (card *Card) UnmarshalText(text []byte) error {
	// trim and split the string into lines
	cardAsString := string(text)
	lines := strings.Split(strings.TrimSpace(cardAsString), "\n")

	// instantiate variables to be filled
	var id string
	var version int
	var lastReview, nextReview time.Time
	var active bool

	var state parseState = header
	var question_lines, answer_lines []string
	for _, line := range lines {
		switch {

		case state == header:
			if idRegex.MatchString(line) {
				rawID := strings.Split(line, "=")[1]
				id = strings.TrimSpace(rawID)

			} else if versionRegex.MatchString(line) {
				rawVersion := strings.Split(line, "=")[1]
				trimmedVersion := strings.TrimSpace(rawVersion)
				value, err := strconv.ParseInt(trimmedVersion, 10, strconv.IntSize)
				if err != nil {
					return fmt.Errorf("failed to parse Version: %w", err)
				}
				version = int(value)

			} else if lastReviewRegex.MatchString(line) {
				var value time.Time
				rawValue := strings.Split(line, "=")[1]
				trimmedValue := strings.TrimSpace(rawValue)
				possibleZeroTime, err := time.Parse(dateLayout, trimmedValue)
				if err != nil {
					return fmt.Errorf("failed to parse LastReview as zero time: %w", err)
				}
				if possibleZeroTime.IsZero() {
					value = possibleZeroTime
				} else {
					value, err = time.ParseInLocation(dateLayout, trimmedValue, time.Local)
					if err != nil {
						return fmt.Errorf("failed to parse LastReview: %w", err)
					}
				}
				lastReview = value

			} else if nextReviewRegex.MatchString(line) {
				rawValue := strings.Split(line, "=")[1]
				trimmedValue := strings.TrimSpace(rawValue)
				value, err := time.ParseInLocation(dateLayout, trimmedValue, time.Local)
				if err != nil {
					return fmt.Errorf("failed to parse NextReview: %w", err)
				}
				nextReview = value

			} else if activeRegex.MatchString(line) {
				rawValue := strings.Split(line, "=")[1]
				trimmedValue := strings.TrimSpace(rawValue)
				value, err := strconv.ParseBool(trimmedValue)
				if err != nil {
					return fmt.Errorf("failed to parse Active: %w", err)
				}
				active = value

			} else if dividerRegex.MatchString(line) {
				state = question
			}

		case state == question:
			if commentRegex.MatchString(line) {
				continue
			} else if dividerRegex.MatchString(line) {
				state = answer
			} else {
				question_lines = append(question_lines, line)
			}

		case state == answer:
			if commentRegex.MatchString(line) {
				continue
			} else {
				answer_lines = append(answer_lines, line)
			}
		}
	}

	card.ID = id
	card.Version = version
	card.LastReview = lastReview
	card.NextReview = nextReview
	card.Active = active
	card.Question = strings.Join(question_lines, "\n")
	card.Answer = strings.Join(answer_lines, "\n")
	return nil
}

// Returns the current review interval, that is, the number of days
// between the date the card was last reviewed and the date the card
// should be reviewed next. Review interval is in days.
func (card *Card) GetCurrentReviewInterval() int {
	// if card has not been studied yet, return 1
	if card.LastReview.IsZero() {
		return 1
	}

	// determine interval if card has been studied
	rawDifference := card.NextReview.Sub(card.LastReview)
	difference := rawDifference.Round(24 * time.Hour)
	days := int(difference / (24 * time.Hour))
	return days
}

// Given a multiplier that represents the effect on the review interval
// following a card review, returns the number of days until that card
// should be reviewed again.
func (card *Card) GetMultipliedReviewInterval(multiplier float64) int {
	currentInterval := card.GetCurrentReviewInterval()
	newInterval := int(math.Round(float64(currentInterval) * multiplier))
	return newInterval
}

// Given a multiplier that represents the effect on the review interval
// following a card review, modifies the fields of the Card that pertain
// to review dates such that they reflect the new interval, with the last
// review set to today.
func (card *Card) SetNextReview(multiplier float64) {
	// get the next interval between reviews
	newInterval := card.GetMultipliedReviewInterval(multiplier)

	// set new last review and new next review
	card.LastReview = time.Now().Truncate(24 * time.Hour)
	card.NextReview = card.LastReview.AddDate(0, 0, newInterval)
}

func (card *Card) IsDue() bool {
	return card.NextReview.Before(time.Now())
}
