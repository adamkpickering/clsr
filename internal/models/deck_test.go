package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeck(t *testing.T) {
	t.Run("NewDeck", func(t *testing.T) {
		name := "test_deck"
		deck := NewDeck(name, true)
		if deck.Name != name {
			t.Errorf("deck has name %s but should have name %s", deck.Name, name)
		}
	})

	t.Run("JSONMarshal", func(t *testing.T) {
		deckName := "test_deck"
		deck := NewDeck(deckName, true)
		card1 := NewCard("this is question 1", "this is question 1", deckName)
		card2 := NewCard("this is question 2", "this is question 2", deckName)
		deck.Cards = append(deck.Cards, card1, card2)
		result, err := json.MarshalIndent(deck, "", "  ")
		if err != nil {
			t.Fatalf("error while marshaling Deck: %s", err)
		}
		fmt.Printf("marshaled deck to %v\n", string(result))
	})

	t.Run("JSONUnmarshal", func(t *testing.T) {
		input := []byte(`{
  "name": "test_deck",
  "version": 0,
  "cards": [
    {
      "id": "e8h3ku7j00",
      "version": 0,
      "active": true,
      "modified": true,
      "question": "this is question 1",
      "answer": "this is question 1",
      "reviews": []
    },
    {
      "id": "1wsrnwciil",
      "version": 0,
      "active": true,
      "modified": true,
      "question": "this is question 2",
      "answer": "this is question 2",
      "reviews": []
    }
  ]
}`)
		deck := Deck{}
		err := json.Unmarshal(input, &deck)
		if err != nil {
			t.Fatalf("error while unmarshaling Deck: %s", err)
		}
		fmt.Printf("unmarshaled deck to %#v\n", deck)
	})
}
