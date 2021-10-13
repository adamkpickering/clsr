package models

import (
	"fmt"
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
	inputCard := Card{
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
`, inputCard.Version, inputCard.LastReview.Format(dateLayout), inputCard.NextReview.Format(dateLayout), inputCard.Active, inputCard.Question, inputCard.Answer)
	parsedCard, err := parseCardFromString(data, id)
	if err != nil {
		t.Errorf("failed to parse card from string: %s", err)
	}
	//if parsedCard != inputCard {
	//	t.Errorf("parsedCard does not match inputCard")
	//}
	if parsedCard.ID != inputCard.ID {
		t.Errorf("mismatched ID")
	}
	if parsedCard.Version != inputCard.Version {
		t.Errorf("mismatched Version")
	}
	if parsedCard.LastReview != inputCard.LastReview {
		t.Errorf("mismatched LastReview")
	}
	if parsedCard.NextReview != inputCard.NextReview {
		t.Errorf("mismatched NextReview")
	}
	if parsedCard.Active != inputCard.Active {
		t.Errorf("mismatched Active")
	}
	if parsedCard.Question != inputCard.Question {
		t.Errorf("mismatched Question")
	}
	if parsedCard.Answer != inputCard.Answer {
		t.Errorf("mismatched Answer")
	}
	fmt.Printf("inputCard: %q\n", inputCard.Answer)
	fmt.Printf("parsedCard: %q\n", parsedCard.Answer)
}

//type Card struct {
//	ID         string
//	Version    int
//	LastReview time.Time
//	NextReview time.Time
//	Question   string
//	Answer     string
//	Active     bool
//}
