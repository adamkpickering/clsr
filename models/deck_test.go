package models

import (
	"path/filepath"
	"testing"
)

func TestNewDeck(t *testing.T) {
	tempdir := t.TempDir()
	name := "test_deck"
	deck_path := filepath.Join(tempdir, name)

	deck, err := NewDeck(deck_path)
	if err != nil {
		t.Errorf("failed to create new deck: %s", err)
	}

	if deck.Name != name {
		t.Errorf("deck has name %s but should have name %s", deck.Name, name)
	}
}
