/*
Copyright © 2021 ADAM PICKERING

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/adamkpickering/clsr/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	setCmd.AddCommand(setCardCmd)
}

var setCardCmd = &cobra.Command{
	Use:   "card <card_id> (active|inactive)",
	Short: "Set whether a card is active or inactive",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}

		// search all the decks for the card
		decks, err := utils.GetDecks(deckSource)
		if err != nil {
			return fmt.Errorf("failed to get decks: %w", err)
		}
		card, deck, err := getCardFromDecks(decks, args[0])
		if err != nil {
			return fmt.Errorf("failed to find card %q: %w", args[0], err)
		}

		// modify the card
		adjective := args[1]
		switch adjective {
		case "active":
			if !card.Active {
				card.Active = true
				card.Modified = true
			}
		case "inactive":
			if card.Active {
				card.Active = false
				card.Modified = true
			}
		default:
			return fmt.Errorf("invalid adjective %q", adjective)
		}

		// write changed card to deck
		err = deckSource.WriteDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to write deck %q: %w", deck.Name, err)
		}

		return nil
	},
}

func getCardFromDecks(decks []*models.Deck, cardID string) (*models.Card, *models.Deck, error) {
	card := &models.Card{}
	deck := &models.Deck{}
	for _, thisDeck := range decks {
		for _, thisCard := range thisDeck.Cards {
			if thisCard.ID == cardID {
				if len(card.ID) != 0 {
					return nil, nil, fmt.Errorf("found a second card with id %q", cardID)
				}
				card = thisCard
				deck = thisDeck
			}
		}
	}
	if card.ID == "" {
		return nil, nil, fmt.Errorf("could not find card with ID %q", cardID)
	}
	return card, deck, nil
}
