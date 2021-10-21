package models

import (
	"path/filepath"
)

// A Deck is a collection of Cards. It is persisted to the filesystem
// as a directory containing files that represent its Cards.
type Deck struct {
	Name  string
	Cards []*Card
}

func NewDeck(path string) *Deck {
	deck := &Deck{
		Name: filepath.Base(path),
	}
	return deck
}

func (d *Deck) AddCard(card *Card) {
	d.Cards = append(d.Cards, card)
}

func (d *Deck) CountCardsDue() int {
	count := 0
	for _, card := range d.Cards {
		if card.IsDue() {
			count += 1
		}
	}
	return count
}
