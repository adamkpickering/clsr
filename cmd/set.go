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
	"errors"
	"fmt"
	"os"

	"github.com/adamkpickering/clsr/models"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <resource_type> <resource_id> <state>",
	Short: "Set the state of various resource types",
	Long: `
Sets the state of various resource types.

There are two states for cards: "active", and "inactive".
"active" means that it will show up when studying cards.
"inactive" means that it will not.`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check that the directory has been initialized
		if _, err := os.Stat(deckDirectory); errors.Is(err, os.ErrNotExist) {
			msg := "could not find %s. Please call `clsr init` and try again."
			return fmt.Errorf(msg, deckDirectory)
		}

		// construct DeckSource
		deckSource, err := models.NewFlatFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to construct DeckSource: %w", err)
		}

		// get deck
		if len(deckName) == 0 {
			return errors.New("must specify --deck/-d")
		}
		deck, err := deckSource.LoadDeck(deckName)
		if err != nil {
			return fmt.Errorf("failed to load deck %q: %w", deckName, err)
		}

		// get the card from the deck
		cardID := args[1]
		card := &models.Card{}
		for _, deckCard := range deck.Cards {
			if deckCard.ID == cardID {
				card = deckCard
				break
			}
		}
		if len(card.ID) == 0 {
			return fmt.Errorf("could not find card with ID %q", cardID)
		}

		resourceType := args[0]
		adjective := args[2]
		switch resourceType {
		case "card":
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
		default:
			return fmt.Errorf("invalid resource type %q", resourceType)
		}

		// sync deck
		err = deckSource.SyncDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to sync deck %q: %w", deck.Name, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
