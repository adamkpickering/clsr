package models

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	name := "test_deck"
	deck := NewDeck(name)
	if deck.Name != name {
		t.Errorf("deck has name %s but should have name %s", deck.Name, name)
	}
}
