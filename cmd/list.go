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

	"github.com/adamkpickering/clsr/models"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list (decks|cards)",
	Short: "List various resources",
	Long:  "\nLists resources.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType := args[0]
		switch resourceType {
		case "cards":
			return listCards()
		case "decks":
			return listDecks()
		default:
			return fmt.Errorf("unrecognized resource %q", resourceType)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listCards() error {
	// get a DeckSource
	deckSource, err := models.NewFlatFileDeckSource(deckDirectory)
	if err != nil {
		return fmt.Errorf("failed to get DeckSource: %w", err)
	}

	// get a list of decks
	var decks []*models.Deck
	if len(deckName) > 0 {
		deck, err := deckSource.LoadDeck(deckName)
		if err != nil {
			return fmt.Errorf("failed to load deck %q: %w", deckName, err)
		}
		decks = append(decks, deck)
	} else {
		deckNames, err := deckSource.ListDecks()
		if err != nil {
			return fmt.Errorf("failed to get deck names: %w", err)
		}
		for _, deckName := range deckNames {
			deck, err := deckSource.LoadDeck(deckName)
			if err != nil {
				return fmt.Errorf("failed to load deck %q: %w", deckName, err)
			}
			decks = append(decks, deck)
		}
	}

	// get a list of cards
	var cards []*models.Card
	for _, deck := range decks {
		for _, card := range deck.Cards {
			cards = append(cards, card)
		}
	}

	// display cards in table
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactClassic)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "Deck"},
			{Text: "Active"},
			{Text: "Last Reviewed"},
			{Text: "Next Review"},
			{Text: "Due"},
		},
	}
	for _, card := range cards {
		var lastReviewed string
		if card.LastReview.IsZero() {
			lastReviewed = "never"
		} else {
			lastReviewed = card.LastReview.Format(models.DateLayout)
		}
		row := []*simpletable.Cell{
			{Text: card.ID},
			{Text: card.Deck},
			{Text: fmt.Sprintf("%t", card.Active)},
			{Text: lastReviewed},
			{Text: card.NextReview.Format(models.DateLayout)},
			{Text: fmt.Sprintf("%dd", card.DaysUntilDue())},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}
	fmt.Println(table.String())

	return nil
}

func listDecks() error {
	// get DeckSource
	deckSource, err := models.NewFlatFileDeckSource(deckDirectory)
	if err != nil {
		return fmt.Errorf("failed to instantiate DeckSource: %w", err)
	}

	// get all decks
	decks, err := getAllDecks(deckSource)
	if err != nil {
		return fmt.Errorf("failed to get decks: %w", err)
	}

	// display decks in table
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactClassic)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "Deck"},
			{Text: "Cards Due"},
			{Text: "Active Cards"},
			{Text: "Inactive Cards"},
			{Text: "Total Cards"},
		},
	}
	for _, deck := range decks {
		row := []*simpletable.Cell{
			{Text: deck.Name},
			{Text: fmt.Sprintf("%d", deck.CountCardsDue())},
			{Text: fmt.Sprintf("%d", deck.CountActiveCards())},
			{Text: fmt.Sprintf("%d", deck.CountInactiveCards())},
			{Text: fmt.Sprintf("%d", len(deck.Cards))},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}
	fmt.Println(table.String())

	return nil
}
