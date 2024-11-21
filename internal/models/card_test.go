package models

import (
	"testing"
)

func TestCard(t *testing.T) {
	t.Run("Copy", func(t *testing.T) {
		oldQuestion := "Old question"
		oldAnswer := "Old answer"
		deckName := "test_deck"
		oldCard := NewCard(oldQuestion, oldAnswer, deckName)
		oldCard.Reviews = append(oldCard.Reviews, NewReview(Hard))

		// get a new card with different field values
		newCard := oldCard.Copy()
		t.Logf("length of newCard.Reviews: %d", len(newCard.Reviews))
		newCard.ID = "asdfqwerxc"
		newCard.Question = "New question"
		newCard.Answer = "New answer"
		newCard.Reviews[0].Result = Normal
		newCard.Reviews = append(newCard.Reviews, NewReview(Easy))

		// check that the fields of oldCard have not been changed
		if newCard.ID == oldCard.ID {
			t.Errorf("newCard.ID matches oldCard.ID")
		}
		if newCard.Question == oldCard.Question {
			t.Errorf("newCard.Question matches oldCard.Question")
		}
		if newCard.Answer == oldCard.Answer {
			t.Errorf("newCard.Answer matches oldCard.Answer")
		}
		if newCard.Reviews[0] == oldCard.Reviews[0] {
			t.Errorf("first elements of newCard.Reviews and oldCard.Reviews are the same")
		}
		if len(newCard.Reviews) == len(oldCard.Reviews) {
			t.Errorf("len(newCard.Reviews) matches len(oldCard.Reviews)")
		}
	})
}
