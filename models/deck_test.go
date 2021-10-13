package models

import "testing"

func TestLoadDeck(t *testing.T) {
	deck, err := LoadDeck("testdata/deck1/")
	if err != nil {
		t.Errorf("failed to load deck: %s", err)
	}

	if length := len(deck.Cards); length != 2 {
		t.Errorf("read %d, not 2, cards for deck", length)
	}
}
