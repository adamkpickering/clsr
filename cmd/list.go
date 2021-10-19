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
	"os"

	"github.com/adamkpickering/clsr/models"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list (decks|cards)",
	Short: "List various resources",
	Long:  "\nLists resources.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType := args[0]
		switch {
		case resourceType == "decks":
			return listDecks()
		default:
			return fmt.Errorf("unrecognized resource \"%s\"", resourceType)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listDecks() error {
	// get list of decks
	entries, err := os.ReadDir(deckDirectory)
	if err != nil {
		return fmt.Errorf("failed to read deck directory: %s", err)
	}
	deckNames := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			deckNames = append(deckNames, entry.Name())
		}
	}

	// load decks
	deckSource, err := models.NewFlatFileDeckSource(deckDirectory)
	if err != nil {
		return fmt.Errorf("failed to instantiate DeckSource: %s", err)
	}
	decks := []*models.Deck{}
	for _, deckName := range deckNames {
		deck, err := deckSource.LoadDeck(deckName)
		if err != nil {
			return fmt.Errorf("failed to load deck \"%s\": %s", deckName, err)
		}
		decks = append(decks, deck)
	}

	// display decks in table
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactClassic)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "Deck"},
			{Text: "Cards Due"},
			{Text: "Total Cards"},
		},
	}
	for _, deck := range decks {
		row := []*simpletable.Cell{
			{Text: deck.Name},
			{Text: fmt.Sprintf("%d", deck.CountCardsDue())},
			{Text: fmt.Sprintf("%d", len(deck.Cards))},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}
	fmt.Println(table.String())

	return nil
}
