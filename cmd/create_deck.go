package cmd

import (
	"fmt"

	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.AddCommand(createDeckCmd)
}

var createDeckCmd = &cobra.Command{
	Use:   "deck <name>",
	Short: "Create deck",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}

		// check for an existing deck of this name
		_, err = deckSource.ReadDeck(deckName)
		if err == nil {
			return fmt.Errorf(`deck %q already exists`, deckName)
		}

		// create the deck
		deck := models.NewDeck(deckName, true)
		err = deckSource.WriteDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to write deck %q: %w", deckName, err)
		}

		return nil
	},
}
