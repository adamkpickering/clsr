package models

// A Deck is a collection of Cards. It is persisted to the filesystem
// as a directory containing files that represent its Cards.
type Deck struct {
	Name  string
	Cards []*Card
}

func NewDeck(name string) *Deck {
	deck := &Deck{
		Name: name,
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

func (d *Deck) CountActiveCards() int {
	count := 0
	for _, card := range d.Cards {
		if card.Active {
			count += 1
		}
	}
	return count
}

func (d *Deck) CountInactiveCards() int {
	count := 0
	for _, card := range d.Cards {
		if !card.Active {
			count += 1
		}
	}
	return count
}
