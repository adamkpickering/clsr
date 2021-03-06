package models

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

var cardFmtString string = `
# here is a comment line
ID = %s
Version = %d
LastReview = %s

NextReview=%s
garbage garbage = asdf
Active = %t
---
%s
---
%s
`

func TestCardUnmarshalText(t *testing.T) {
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
	data := fmt.Sprintf(cardFmtString, inputCard.ID, inputCard.Version, inputCard.LastReview.Format(DateLayout), inputCard.NextReview.Format(DateLayout), inputCard.Active, inputCard.Question, inputCard.Answer)
	parsedCard := &Card{}
	err := parsedCard.UnmarshalText([]byte(data))
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
	oldCard := NewCard("fake question", "fake answer", "fake_deck")
	oldCard.Modified = false
	cardPath := filepath.Join(tempdir, GetCardFilename(oldCard))
	err := WriteCardToFile(cardPath, oldCard)
	if err != nil {
		t.Errorf("failed to write card to file: %s", err)
	}
	newCard, err := ReadCardFromFile(cardPath)
	if err != nil {
		t.Errorf("failed to parse card from file: %s", err)
	}
	if *newCard != *oldCard {
		t.Error("new and old cards do not match")
		t.Logf("old card: %#v\n", *oldCard)
		t.Logf("new card: %#v\n", *newCard)
	}
}

func TestGetCurrentReviewInterval(t *testing.T) {
	type testCase struct {
		LastReview time.Time
		NextReview time.Time
		Result     int
	}
	testCases := []testCase{
		{
			LastReview: time.Date(2021, 06, 24, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 06, 25, 0, 0, 0, 0, time.Local),
			Result:     1,
		},
		{
			LastReview: time.Date(2021, 11, 6, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 11, 8, 0, 0, 0, 0, time.Local),
			Result:     2,
		},
		{
			LastReview: time.Date(2021, 11, 8, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 11, 6, 0, 0, 0, 0, time.Local),
			Result:     -2,
		},
		{
			LastReview: time.Time{},
			NextReview: time.Date(2021, 11, 6, 0, 0, 0, 0, time.Local),
			Result:     1,
		},
	}
	for _, test := range testCases {
		card := NewCard("fake question", "fake answer", "fake_deck")
		card.LastReview = test.LastReview
		card.NextReview = test.NextReview
		interval := card.GetCurrentReviewInterval()
		if interval != test.Result {
			t.Errorf("Got unexpected interval %d (%d expected)", interval, test.Result)
			t.Logf("card: %#v\n", *card)
		}
	}
}

func TestGetMultipliedReviewInterval(t *testing.T) {
	type testCase struct {
		LastReview time.Time
		NextReview time.Time
		Multiplier float64
		Result     int
	}
	testCases := []testCase{
		{
			LastReview: time.Date(2021, 06, 24, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 06, 25, 0, 0, 0, 0, time.Local),
			Multiplier: 1.0,
			Result:     1,
		},
		{
			LastReview: time.Date(2021, 06, 24, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 06, 25, 0, 0, 0, 0, time.Local),
			Multiplier: 3.3,
			Result:     3,
		},
		{
			LastReview: time.Date(2021, 06, 24, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 06, 25, 0, 0, 0, 0, time.Local),
			Multiplier: 3.8,
			Result:     4,
		},
		{
			LastReview: time.Date(2021, 06, 24, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 06, 25, 0, 0, 0, 0, time.Local),
			Multiplier: 0.0,
			Result:     0,
		},
		{
			LastReview: time.Date(2021, 11, 6, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 11, 8, 0, 0, 0, 0, time.Local),
			Multiplier: 0.0,
			Result:     0,
		},
		{
			LastReview: time.Date(2021, 11, 6, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 11, 8, 0, 0, 0, 0, time.Local),
			Multiplier: 1.0,
			Result:     2,
		},
		{
			LastReview: time.Date(2021, 11, 6, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 11, 8, 0, 0, 0, 0, time.Local),
			Multiplier: 3.5,
			Result:     7,
		},
		{
			LastReview: time.Date(2021, 06, 24, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 06, 25, 0, 0, 0, 0, time.Local),
			Multiplier: -1.0,
			Result:     -1,
		},
		{
			LastReview: time.Date(2021, 11, 6, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 11, 8, 0, 0, 0, 0, time.Local),
			Multiplier: -0.5,
			Result:     -1,
		},
		{
			LastReview: time.Date(2021, 11, 6, 0, 0, 0, 0, time.Local),
			NextReview: time.Date(2021, 11, 8, 0, 0, 0, 0, time.Local),
			Multiplier: 0.5,
			Result:     1,
		},
	}
	for _, test := range testCases {
		card := NewCard("fake question", "fake answer", "fake_deck")
		card.LastReview = test.LastReview
		card.NextReview = test.NextReview
		interval := card.GetMultipliedReviewInterval(test.Multiplier)
		if interval != test.Result {
			t.Errorf("Got unexpected interval %d (%d expected)", interval, test.Result)
			t.Logf("card: %#v\n", *card)
		}
	}
}

func TestSetNextReview(t *testing.T) {
	card := NewCard("fake question", "fake answer", "fake_deck")
	if !card.LastReview.IsZero() {
		t.Error("card.LastReview is not zero valued")
	}
	card.SetNextReview(1.0)
	if interval := card.GetCurrentReviewInterval(); interval != 1 {
		t.Errorf("got interval %d (expected 1)", interval)
	}
	card.SetNextReview(2.0)
	if interval := card.GetCurrentReviewInterval(); interval != 2 {
		t.Errorf("got interval %d (expected 2)", interval)
	}
	card.SetNextReview(4.0)
	if interval := card.GetCurrentReviewInterval(); interval != 8 {
		t.Errorf("got interval %d (expected 8)", interval)
	}
	card.SetNextReview(4.0)
	if interval := card.GetCurrentReviewInterval(); interval != 32 {
		t.Errorf("got interval %d (expected 32)", interval)
	}
}

// Tests that any lines starting with # in card .txt files are ignored,
// whether they are in the header fields, the question field or the
// answer field.
func TestHonorComments(t *testing.T) {
	// create a card string with comments
	questionComment := "# this is a comment line"
	questionContent := "this is the actual question"
	question := questionComment + "\n" + questionContent
	answerComment := "# this is a comment line"
	answerContent := "this is the actual answer"
	answer := answerComment + "\n" + answerContent
	inputCard := NewCard(question, answer, "fake_deck")
	inputCardString := fmt.Sprintf(cardFmtString, inputCard.ID, inputCard.Version, inputCard.LastReview.Format(DateLayout), inputCard.NextReview.Format(DateLayout), inputCard.Active, inputCard.Question, inputCard.Answer)

	// parse the card and ensure that comments are not there
	parsedCard := &Card{}
	err := parsedCard.UnmarshalText([]byte(inputCardString))
	if err != nil {
		t.Errorf("failed to parse card from string: %s", err)
	}
	if parsedCard.Question != questionContent {
		fmtString := "parsed card question %q does not match expected (%q)"
		t.Errorf(fmtString, parsedCard.Question, questionContent)
	}
	if parsedCard.Answer != answerContent {
		fmtString := "parsed card answer %q does not match expected (%q)"
		t.Errorf(fmtString, parsedCard.Answer, answerContent)
	}
}

func TestDaysUntilDue(t *testing.T) {
	testCases := []int{-4, -1, 0, 1, 5}
	for _, testCase := range testCases {
		year, month, day := time.Now().Date()
		today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		card := NewCard("fake question", "fake answer", "fake_deck")
		card.NextReview = today.AddDate(0, 0, testCase)
		if due := card.DaysUntilDue(); due != testCase {
			t.Errorf("Got %d days until due (%d expected)", due, testCase)
		}
	}
}
