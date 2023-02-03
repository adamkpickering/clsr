package models

// A Deck is a collection of Cards that are all related.
type Deck struct {
	Name    string  `json:"name"`
	Version int     `json:"version"`
	Cards   []*Card `json:"cards"`
}

func NewDeck(name string) *Deck {
	deck := &Deck{
		Name:    name,
		Version: 0,
		Cards:   []*Card{},
	}
	return deck
}

func (deck *Deck) Copy() *Deck {
	copiedDeck := NewDeck(deck.Name)
	copiedDeck.Cards = make([]*Card, 0, len(deck.Cards))
	for _, card := range deck.Cards {
		copiedDeck.Cards = append(copiedDeck.Cards, card.Copy())
	}
	return copiedDeck
}
