package models

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func TestParseCardFromString(t *testing.T) {
	id := "1234asdf"
	question := `
here is the question
it has

multiple lines

`
	answer := `
this is the answer

if you got it
good job`
	inputCard := &Card{
		ID:         id,
		Version:    0,
		LastReview: time.Date(2021, 06, 12, 0, 0, 0, 0, time.UTC),
		NextReview: time.Date(2021, 06, 13, 0, 0, 0, 0, time.UTC),
		Active:     true,
		Question:   question,
		Answer:     answer,
	}
	data := fmt.Sprintf(`
Version = %d
LastReview = %s

NextReview=%s
garbage garbage = asdf
Active = %t
---
%s
---
%s
`, inputCard.Version, inputCard.LastReview.Format(time.RFC3339), inputCard.NextReview.Format(time.RFC3339), inputCard.Active, inputCard.Question, inputCard.Answer)
	parsedCard, err := ParseCardFromString(data, id)
	if err != nil {
		t.Errorf("failed to parse card from string: %s", err)
	}
	if parsedCard.ID != inputCard.ID {
		t.Error("mismatched ID")
	}
	if parsedCard.Version != inputCard.Version {
		t.Error("mismatched Version")
	}
	if parsedCard.LastReview != inputCard.LastReview {
		t.Error("mismatched LastReview")
	}
	if parsedCard.NextReview != inputCard.NextReview {
		t.Error("mismatched NextReview")
	}
	if parsedCard.Active != inputCard.Active {
		t.Error("mismatched Active")
	}
	if parsedCard.Question != inputCard.Question {
		t.Error("mismatched Question")
	}
	if parsedCard.Answer != inputCard.Answer {
		t.Error("mismatched Answer")
	}
}

func TestCardWriteRead(t *testing.T) {
	tempdir := t.TempDir()
	oldCard := NewCard("fake question", "fake answer")
	filename := fmt.Sprintf("%s.txt", oldCard.ID)
	err := oldCard.WriteToDir(tempdir)
	if err != nil {
		t.Errorf("failed to write card to file: %s", err)
	}
	newCard, err := parseCardFromFile(filepath.Join(tempdir, filename))
	if err != nil {
		t.Errorf("failed to parse card from file: %s", err)
	}
	if *newCard != *oldCard {
		t.Error("new and old cards do not match")
	}
	fmt.Printf("old card: %#v\n", oldCard)
	fmt.Printf("new card: %#v\n", newCard)
}
