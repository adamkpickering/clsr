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
