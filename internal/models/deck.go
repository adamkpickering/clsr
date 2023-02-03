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

// The below commented code is probably better implemented elsewhere.

// func (deck *Deck) AddCard(card *Card) {
// 	deck.Cards = append(deck.Cards, card)
// }

// func (deck *Deck) CountCardsDue() int {
// 	count := 0
// 	for _, card := range deck.Cards {
// 		if card.IsDue() {
// 			count += 1
// 		}
// 	}
// 	return count
// }

// func (deck *Deck) CountActiveCards() int {
// 	count := 0
// 	for _, card := range deck.Cards {
// 		if card.Active {
// 			count += 1
// 		}
// 	}
// 	return count
// }

// func (deck *Deck) CountInactiveCards() int {
// 	count := 0
// 	for _, card := range deck.Cards {
// 		if !card.Active {
// 			count += 1
// 		}
// 	}
// 	return count
// }
