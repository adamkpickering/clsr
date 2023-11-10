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
