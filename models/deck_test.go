package models

import (
	"path/filepath"
	"testing"
)

func TestNewDeck(t *testing.T) {
	tempdir := t.TempDir()
	name := "test_deck"
	deck_path := filepath.Join(tempdir, name)
	deck := NewDeck(deck_path)
	if deck.Name != name {
		t.Errorf("deck has name %s but should have name %s", deck.Name, name)
	}
}
