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

	"github.com/adamkpickering/clsr/pkg/config"
	"github.com/adamkpickering/clsr/pkg/deck_source"
	"github.com/adamkpickering/clsr/pkg/models"
	"github.com/adamkpickering/clsr/pkg/scheduler"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
)

const DateLayout = "2006-01-02"

type cardRow struct {
	ID           string
	Deck         string
	Active       string
	ReviewCount  string
	LastReviewed string
	NextReview   string
	Question     string
	Answer       string
}

var listCardFlags = struct {
	DeckName string
}{}

func init() {
	listCmd.AddCommand(listCardCmd)
	listCardCmd.Flags().StringVarP(&listCardFlags.DeckName, "deck", "d", "", "filter cards by deck")
}

var listCardCmd = &cobra.Command{
	Use:   "cards",
	Short: "List cards",
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := listCardFlags.DeckName

		// get a DeckSource
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to get DeckSource: %w", err)
		}
		scheduler := scheduler.NewTwoReviewScheduler(config.DefaultConfig)

		// get a list of decks
		decks := []*models.Deck{}
		if cmd.Flags().Changed("deck") {
			decks, err = getDecks(deckSource, deckName)
		} else {
			decks, err = getDecks(deckSource)
		}
		if err != nil {
			return fmt.Errorf("failed to get decks: %w", err)
		}

		// get a list of cards
		var cardRows []cardRow
		for _, deck := range decks {
			for _, card := range deck.Cards {
				cardRow, err := cardToCardRow(card, scheduler)
				if err != nil {
					return fmt.Errorf("failed to convert Card %q to cardRow: %w", card.ID, err)
				}
				cardRows = append(cardRows, cardRow)
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
				{Text: "Review Count"},
				{Text: "Last Reviewed"},
				{Text: "Next Review"},
				// {Text: "Due"},
			},
		}
		for _, cardRow := range cardRows {
			row := []*simpletable.Cell{
				{Text: cardRow.ID},
				{Text: cardRow.Deck},
				{Text: cardRow.Active},
				{Text: cardRow.ReviewCount},
				{Text: cardRow.LastReviewed},
				{Text: cardRow.NextReview},
				// {Text: fmt.Sprintf("%dd", cardRow.DaysUntilDue())},
			}
			table.Body.Cells = append(table.Body.Cells, row)
		}
		fmt.Println(table.String())

		return nil
	},
}

func cardToCardRow(card *models.Card, scheduler scheduler.Scheduler) (cardRow, error) {
	nextReview, err := scheduler.GetNextReview(card)
	if err != nil {
		return cardRow{}, fmt.Errorf("failed to get next review: %w", err)
	}
	cardRow := cardRow{
		ID:          card.ID,
		Deck:        card.Deck,
		Active:      fmt.Sprintf("%t", card.Active),
		ReviewCount: fmt.Sprintf("%d", len(card.Reviews)),
		NextReview:  nextReview.Format(DateLayout),
		Question:    card.Question,
		Answer:      card.Answer,
	}
	if len(card.Reviews) == 0 {
		cardRow.LastReviewed = "never"
	} else {
		cardRow.LastReviewed = card.Reviews[0].Datetime.Format(DateLayout)
	}
	return cardRow, nil
}
