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
	"time"

	"github.com/adamkpickering/clsr/internal/config"
	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/adamkpickering/clsr/internal/scheduler"
	"github.com/adamkpickering/clsr/internal/utils"
	"github.com/spf13/cobra"
)

type CardRow struct {
	ID           string
	Deck         string
	Active       bool
	ReviewCount  int
	LastReviewed string
	NextReview   string
	Question     string
	Answer       string
}

var listCardFlags = struct {
	DeckNames []string
}{}

func init() {
	listCmd.AddCommand(listCardCmd)
	listCardCmd.Flags().StringSliceVarP(&listCardFlags.DeckNames, "decks", "d", []string{}, "only list cards from these decks")
}

var listCardCmd = &cobra.Command{
	Use:   "cards",
	Short: "List cards",
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := listCardFlags.DeckNames
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}
		scheduler := scheduler.NewTwoReviewScheduler(config.DefaultConfig)

		// get a list of cards
		cards, err := utils.GetCards(deckSource, deckName...)
		if err != nil {
			return fmt.Errorf("failed to get cards: %w", err)
		}

		// convert cards to CardRows
		var cardRows []CardRow
		for _, card := range cards {
			cardRow, err := cardToCardRow(card, scheduler)
			if err != nil {
				return fmt.Errorf("failed to convert Card %q to CardRow: %w", card.ID, err)
			}
			cardRows = append(cardRows, cardRow)
		}

		return printCardTable(cardRows)
	},
}

func printCardTable(cardRows []CardRow) error {
	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	_, err := fmt.Fprintln(writer, "ID\tDeck\tActive\tReview Count\tLast Reviewed\tNext Review")
	if err != nil {
		return fmt.Errorf("failed to write header row: %w", err)
	}
	for _, cardRow := range cardRows {
		_, err = fmt.Fprintf(writer, "%s\t%s\t%t\t%d\t%s\t%s\n",
			cardRow.ID,
			cardRow.Deck,
			cardRow.Active,
			cardRow.ReviewCount,
			cardRow.LastReviewed,
			cardRow.NextReview,
		)
		if err != nil {
			return fmt.Errorf("failed to write row for card %q: %w", cardRow.ID, err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	return nil
}

func cardToCardRow(card *models.Card, scheduler scheduler.Scheduler) (CardRow, error) {
	row := CardRow{
		ID:          card.ID,
		Deck:        card.Deck,
		Active:      card.Active,
		ReviewCount: len(card.Reviews),
		Question:    card.Question,
		Answer:      card.Answer,
	}

	// deal with NextReview
	nextReview, err := scheduler.GetNextReview(card)
	if err != nil {
		return CardRow{}, fmt.Errorf("failed to get next review: %w", err)
	}
	due, err := scheduler.IsDue(card)
	if err != nil {
		return CardRow{}, fmt.Errorf("failed during check whether card is due: %w", err)
	}
	if due {
		row.NextReview = "due"
	} else {
		row.NextReview = utils.GetReadableTimeDifference(time.Now(), nextReview)
	}

	// deal with LastReviewed
	if len(card.Reviews) == 0 {
		row.LastReviewed = "never"
	} else {
		readableTimeDifference := utils.GetReadableTimeDifference(card.Reviews[0].Datetime, time.Now())
		lastReviewed := fmt.Sprintf("%s ago", readableTimeDifference)
		row.LastReviewed = lastReviewed
	}
	return row, nil
}
