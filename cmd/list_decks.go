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
	"text/tabwriter"

	"github.com/adamkpickering/clsr/internal/config"
	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/adamkpickering/clsr/internal/scheduler"
	"github.com/adamkpickering/clsr/internal/utils"
	"github.com/spf13/cobra"
)

var listDeckFlags = struct {
	All bool
}{}

func init() {
	listCmd.AddCommand(listDeckCmd)
	listDeckCmd.Flags().BoolVarP(&listDeckFlags.All, "all", "a", false, "list all decks, not just active ones")
}

var listDeckCmd = &cobra.Command{
	Use:   "decks",
	Short: "List decks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}
		allDecks, err := utils.GetDecks(deckSource)
		if err != nil {
			return fmt.Errorf("failed to get decks: %w", err)
		}
		// filter decks if necessary
		decks := make([]*models.Deck, 0, len(allDecks))
		if listDeckFlags.All {
			decks = allDecks
		} else {
			for _, deck := range allDecks {
				if deck.Active {
					decks = append(decks, deck)
				}
			}
		}
		return printDeckTable(decks)
	},
}

func printDeckTable(decks []*models.Deck) error {
	scheduler := scheduler.NewTwoReviewScheduler(config.DefaultConfig)
	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	_, err := fmt.Fprintln(writer, "Deck\tCards Due\tActive Cards\tInactive Cards\tTotal Cards\tActive")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	for _, deck := range decks {
		dueCount, err := countCardsDue(deck, scheduler)
		if err != nil {
			return fmt.Errorf("failed to count due cards: %w", err)
		}
		activeCount, inactiveCount := countActiveCards(deck)
		_, err = fmt.Fprintf(writer, "%s\t%d\t%d\t%d\t%d\t%t\n",
			deck.Name,
			dueCount,
			activeCount,
			inactiveCount,
			len(deck.Cards),
			deck.Active,
		)
		if err != nil {
			return fmt.Errorf("failed to write row for deck %q: %w", deck.Name, err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	return nil
}

func countCardsDue(deck *models.Deck, scheduler scheduler.Scheduler) (int, error) {
	count := 0
	for _, card := range deck.Cards {
		cardIsDue, err := scheduler.IsDue(card)
		if err != nil {
			return 0, fmt.Errorf("failed to check if card is due: %w", err)
		}
		if cardIsDue && card.Active {
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
