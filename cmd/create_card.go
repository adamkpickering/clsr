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
			return fmt.Errorf("failed to get user input: %w", err)
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
