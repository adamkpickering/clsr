package models

import (
	"fmt"
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

const dateLayout = "2006-01-02T15:04:05Z"

var versionRegex *regexp.Regexp
var lastReviewRegex *regexp.Regexp
var nextReviewRegex *regexp.Regexp
var activeRegex *regexp.Regexp
var dividerRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile(`^Version *=`)
	lastReviewRegex = regexp.MustCompile(`^LastReview *=`)
	nextReviewRegex = regexp.MustCompile(`^NextReview *=`)
	activeRegex = regexp.MustCompile(`^Active *=`)
	dividerRegex = regexp.MustCompile(`^---`)
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

type Card struct {
	ID         string
	Version    int
	LastReview time.Time
	NextReview time.Time
	Active     bool
	Question   string
	Answer     string
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

//func ParseCardFromFile(path string) (Card, error) {
//	// check that path exists
//	// read card into string
//	// parse the card and return
//}

func parseCardFromString(data string, id string) (Card, error) {
	card := Card{
		ID: id,
	}

	// trim and split the string into lines
	lines := strings.Split(strings.TrimSpace(data), "\n")

	var state parseState = header
	var question_lines, answer_lines []string
	for _, line := range lines {
		switch {

		case state == header:
			if versionRegex.MatchString(line) {
				raw_version := strings.Split(line, "=")[1]
				trimmed_version := strings.TrimSpace(raw_version)
				version, err := strconv.ParseInt(trimmed_version, 10, strconv.IntSize)
				if err != nil {
					return Card{}, fmt.Errorf("failed to parse Version: %s", err)
				}
				card.Version = int(version)

			} else if lastReviewRegex.MatchString(line) {
				raw_value := strings.Split(line, "=")[1]
				trimmed_value := strings.TrimSpace(raw_value)
				value, err := time.Parse(dateLayout, trimmed_value)
				if err != nil {
					return Card{}, fmt.Errorf("failed to parse LastReview: %s", err)
				}
				card.LastReview = value

			} else if nextReviewRegex.MatchString(line) {
				raw_value := strings.Split(line, "=")[1]
				trimmed_value := strings.TrimSpace(raw_value)
				value, err := time.Parse(dateLayout, trimmed_value)
				if err != nil {
					return Card{}, fmt.Errorf("failed to parse NextReview: %s", err)
				}
				card.NextReview = value

			} else if activeRegex.MatchString(line) {
				raw_value := strings.Split(line, "=")[1]
				trimmed_value := strings.TrimSpace(raw_value)
				value, err := strconv.ParseBool(trimmed_value)
				if err != nil {
					return Card{}, fmt.Errorf("failed to parse Active: %s", err)
				}
				card.Active = value

			} else if dividerRegex.MatchString(line) {
				state = question
			}

		case state == question:
			if dividerRegex.MatchString(line) {
				state = answer
			} else {
				question_lines = append(question_lines, line)
			}

		case state == answer:
			answer_lines = append(answer_lines, line)
		}
	}

	card.Question = strings.Join(question_lines, "\n")
	card.Answer = strings.Join(answer_lines, "\n")

	return card, nil
}
