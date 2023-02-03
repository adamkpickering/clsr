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

	"github.com/adamkpickering/clsr/internal/config"
	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/adamkpickering/clsr/internal/scheduler"
	"github.com/adamkpickering/clsr/internal/utils"
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
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}
		scheduler := scheduler.NewTwoReviewScheduler(config.DefaultConfig)

		// get all decks
		decks, err := utils.GetDecks(deckSource)
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
			dueCount, err := countCardsDue(deck, scheduler)
			if err != nil {
				return fmt.Errorf("failed to count due cards: %w", err)
			}
			activeCount, inactiveCount := countActiveCards(deck)
			row := []*simpletable.Cell{
				{Text: deck.Name},
				{Text: fmt.Sprintf("%d", dueCount)},
				{Text: fmt.Sprintf("%d", activeCount)},
				{Text: fmt.Sprintf("%d", inactiveCount)},
				{Text: fmt.Sprintf("%d", len(deck.Cards))},
			}
			table.Body.Cells = append(table.Body.Cells, row)
		}
		fmt.Println(table.String())

		return nil
	},
}

func countCardsDue(deck *models.Deck, scheduler scheduler.Scheduler) (int, error) {
	count := 0
	for _, card := range deck.Cards {
		cardIsDue, err := scheduler.IsDue(card)
		if err != nil {
			return 0, fmt.Errorf("failed to check if card is due: %w", err)
		}
		if cardIsDue {
			count += 1
		}
	}
	return count, nil
}

// Returns two counts. The first is the count of active cards,
// and the second is the count of inactive cards.
func countActiveCards(deck *models.Deck) (int, int) {
	activeCount := 0
	inactiveCount := 0
	for _, card := range deck.Cards {
		if card.Active {
			activeCount += 1
		} else {
			inactiveCount += 1
		}
	}
	return activeCount, inactiveCount
}
