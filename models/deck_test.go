package models

import (
	"path/filepath"
	"testing"
)

func TestLoadDeck(t *testing.T) {
	deck, err := LoadDeck("testdata/deck1/")
	if err != nil {
		t.Errorf("failed to load deck: %s", err)
	}

	if length := len(deck.Cards); length != 2 {
		t.Errorf("read %d, not 2, cards for deck", length)
	}
}

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
