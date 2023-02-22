package deck_source

import (
	"github.com/adamkpickering/clsr/internal/models"
	"testing"
)

func TestJSONFileDeckSource(t *testing.T) {
	t.Run("ReadDeck", func(t *testing.T) {
		deckSource, err := NewJSONFileDeckSource("testdata")
		if err != nil {
			t.Errorf("failed to instantiate FlatFileDeckSource: %s", err)
		}
		deck, err := deckSource.ReadDeck("test_deck")
		if err != nil {
			t.Errorf("failed to read deck: %s", err)
		}
		if length := len(deck.Cards); length != 2 {
			t.Errorf("read %d, not 2, cards for deck", length)
		}
	})

	t.Run("WriteDeck", func(t *testing.T) {
		// create a deck in temporary directory
		testDeckName := "test_deck"
		tempDir := t.TempDir()
		deckSource, err := NewJSONFileDeckSource(tempDir)
		if err != nil {
			t.Fatalf("failed to create deck source: %s", err)
		}
		initialDeck := models.NewDeck(testDeckName, true)
		card1 := models.NewCard("card1 question", "card1 answer", testDeckName)
		card2 := models.NewCard("card2 question", "card2 answer", testDeckName)
		initialDeck.Cards = []*models.Card{card1, card2}
		err = deckSource.WriteDeck(initialDeck)
		if err != nil {
			t.Fatalf("failed to write initial deck: %s", err)
		}

		// read the deck from the tempdir
		deck, err := deckSource.ReadDeck(testDeckName)
		if err != nil {
			t.Fatalf("failed to read deck: %s", err)
		}
		if len(deck.Cards) < 2 {
			t.Fatal("failed to read the expected number of cards")
		}
		t.Logf("%#v", deck.Cards[0])
		t.Logf("%#v", deck.Cards[1])
	})
}
