package cmd

import (
	"fmt"

	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/adamkpickering/clsr/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	editCmd.AddCommand(editCardCmd)
}

var editCardCmd = &cobra.Command{
	Use:   "card <card_id>",
	Short: "Edit card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// search for the card in all decks
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}
		decks, err := utils.GetDecks(deckSource)
		if err != nil {
			return fmt.Errorf("failed to get decks: %w", err)
		}
		card, deck, err := getCardFromDecks(decks, args[0])
		if err != nil {
			return fmt.Errorf("failed to find card %q: %w", args[0], err)
		}

		// edit the card
		if err := models.EditCardViaEditor(card); err == models.ErrNotModified {
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to edit card: %w", err)
		}

		// write changed card to deck
		if err = deckSource.WriteDeck(deck); err != nil {
			return fmt.Errorf("failed to write deck %q: %w", deck.Name, err)
		}

		return nil
	},
}
