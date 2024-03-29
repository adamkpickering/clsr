package deck_source

import (
	"github.com/adamkpickering/clsr/internal/models"
)

type DeckSource interface {
	ReadDeck(name string) (*models.Deck, error)
	WriteDeck(deck *models.Deck) error
	ListDecks() ([]string, error)
	DeleteDeck(name string) error
}
