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
		LastReview: time.Date(2021, 06, 12, 0, 0, 0, 0, time.Local),
		NextReview: time.Date(2021, 06, 13, 0, 0, 0, 0, time.Local),
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
`, inputCard.Version, inputCard.LastReview.Format(dateLayout), inputCard.NextReview.Format(dateLayout), inputCard.Active, inputCard.Question, inputCard.Answer)
	parsedCard, err := ParseCardFromString(data, id)
	if err != nil {
		t.Errorf("failed to parse card from string: %s", err)
	}
	if *inputCard != *parsedCard {
		t.Error("inputCard does not match parsedCard")
		t.Logf("inputCard: %#v\n", *inputCard)
		t.Logf("parsedCard: %#v\n", *parsedCard)
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
		t.Logf("old card: %#v\n", *oldCard)
		t.Logf("new card: %#v\n", *newCard)
	}
}
