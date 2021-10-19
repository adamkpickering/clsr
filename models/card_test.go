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
	}
	for _, test := range testCases {
		card := NewCard("fake question", "fake answer")
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
		card := NewCard("fake question", "fake answer")
		card.LastReview = test.LastReview
		card.NextReview = test.NextReview
		interval := card.GetMultipliedReviewInterval(test.Multiplier)
		if interval != test.Result {
			t.Errorf("Got unexpected interval %d (%d expected)", interval, test.Result)
			t.Logf("card: %#v\n", *card)
		}
	}
}

//func TestSetNextReview(t *testing.T) {
//	type testCase struct {
//		Multiplier float64
//		Result     int
//		CardFunc   func() *Card
//	}
//	generateCard := func() *Card {
//		card := NewCard("fake question", "fake answer")
//		// remove the zero value from LastReview
//		card.SetNextReview(1.0)
//		return card
//	}
//	testCases := []testCase{
//		{Multiplier: 1.0, Result: 1, CardFunc: generateCard},
//		{Multiplier: 0.0, Result: 0, CardFunc: generateCard},
//		{Multiplier: 4.0, Result: 4, CardFunc: generateCard},
//		{Multiplier: 4.5, Result: 5, CardFunc: generateCard},
//		{Multiplier: 4.3, Result: 4, CardFunc: generateCard},
//		{Multiplier: -3.2, Result: 0, CardFunc: generateCard},
//	}
//	for _, myCase := range testCases {
//		card := myCase.CardFunc()
//		card.SetNextReview(myCase.Multiplier)
//		interval := card.GetCurrentReviewInterval()
//		if interval != myCase.Result {
//			t.Errorf("got incorrect review interval %d (expected %d)", interval, myCase.Result)
//		}
//	}
//}
