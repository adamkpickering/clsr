package cmd

import (
	"fmt"

	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/spf13/cobra"
)

func init() {
	setCmd.AddCommand(setDeckCmd)
}

var setDeckCmd = &cobra.Command{
	Use:   "deck <deck_name> (active|inactive)",
	Short: "Set whether a deck is active or inactive",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}

		// get deck
		deck, err := deckSource.ReadDeck(args[0])
		if err != nil {
			return fmt.Errorf("failed to get deck: %w", err)
		}

		// modify the deck
		switch args[1] {
		case "active":
			if deck.Active {
				return nil
			}
			deck.Active = true
		case "inactive":
			if !deck.Active {
				return nil
			}
			deck.Active = false
		default:
			return fmt.Errorf("invalid adjective %q", args[1])
		}

		// write deck
		err = deckSource.WriteDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to write deck %q: %w", deck.Name, err)
		}

		return nil
	},
}
