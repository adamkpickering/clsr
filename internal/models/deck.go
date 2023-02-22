package models

// A Deck is a collection of Cards that are all related.
type Deck struct {
	Name    string  `json:"name"`
	Version int     `json:"version"`
	Active  bool    `json:"active"`
	Cards   []*Card `json:"cards"`
}

func NewDeck(name string, active bool) *Deck {
	deck := &Deck{
		Name:    name,
		Version: 0,
		Active:  active,
		Cards:   []*Card{},
	}
	return deck
}

func (deck *Deck) Copy() *Deck {
	copiedDeck := NewDeck(deck.Name, deck.Active)
	copiedDeck.Cards = make([]*Card, 0, len(deck.Cards))
	for _, card := range deck.Cards {
		copiedDeck.Cards = append(copiedDeck.Cards, card.Copy())
	}
	return copiedDeck
}
