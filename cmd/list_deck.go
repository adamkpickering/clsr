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
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listDeckCmd)
}

var listDeckCmd = &cobra.Command{
	Use:   "decks",
	Short: "List decks",
	RunE: func(cmd *cobra.Command, args []string) error {
		// get DeckSource
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate DeckSource: %w", err)
		}

		// get all decks
		decks, err := getDecks(deckSource)
		if err != nil {
			return fmt.Errorf("failed to get decks: %w", err)
		}

		// display decks in table
		table := simpletable.New()
		table.SetStyle(simpletable.StyleCompactClassic)
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Text: "Deck"},
				// {Text: "Cards Due"},
				// {Text: "Active Cards"},
				// {Text: "Inactive Cards"},
				{Text: "Total Cards"},
			},
		}
		for _, deck := range decks {
			row := []*simpletable.Cell{
				{Text: deck.Name},
				// {Text: fmt.Sprintf("%d", deck.CountCardsDue())},
				// {Text: fmt.Sprintf("%d", deck.CountActiveCards())},
				// {Text: fmt.Sprintf("%d", deck.CountInactiveCards())},
				{Text: fmt.Sprintf("%d", len(deck.Cards))},
			}
			table.Body.Cells = append(table.Body.Cells, row)
		}
		fmt.Println(table.String())

		return nil
	},
}
