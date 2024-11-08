package cmd

import (
	"fmt"

	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/spf13/cobra"
)

var createCardFlags = struct {
	DeckName string
}{}

func init() {
	createCmd.AddCommand(createCardCmd)
	createCardCmd.Flags().StringVarP(&createCardFlags.DeckName, "deck", "d", "", "filter cards by deck")
	createCardCmd.MarkFlagRequired("deck")
}

var createCardCmd = &cobra.Command{
	Use:   "card",
	Short: "Create card",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := createCardFlags.DeckName
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}

		// read the deck
		deck, err := deckSource.ReadDeck(deckName)
		if err != nil {
			return fmt.Errorf("failed to read deck %q: %w", deckName, err)
		}

		// exec into editor to get Card fields from user
		card := models.NewCard("", "", deckName)
		if err := models.EditCardViaEditor(card); err == models.ErrNotModified {
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to edit card: %w", err)
		}

		// add the Card to the Deck and write the Deck
		deck.Cards = append(deck.Cards, card)
		err = deckSource.WriteDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to save deck: %w", err)
		}

		return nil
	},
}
