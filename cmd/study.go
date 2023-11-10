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
	"math/rand"

	"github.com/adamkpickering/clsr/internal/config"
	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/adamkpickering/clsr/internal/scheduler"
	"github.com/adamkpickering/clsr/internal/utils"
	"github.com/adamkpickering/clsr/internal/views"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

var studyFlags = struct {
	DeckName string
}{}

func init() {
	rootCmd.AddCommand(studyCmd)
	studyCmd.Flags().StringVarP(&studyFlags.DeckName, "deck", "d", "", "study a specific deck")
}

var studyCmd = &cobra.Command{
	Use:   "study",
	Short: "Study cards that are due",
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := studyFlags.DeckName
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}
		scheduler := scheduler.NewTwoReviewScheduler(config.DefaultConfig)

		// get a list of decks
		decks := []*models.Deck{}
		if cmd.Flags().Changed("deck") {
			decks, err = utils.GetDecks(deckSource, deckName)
		} else {
			decks, err = utils.GetDecks(deckSource)
		}
		if err != nil {
			return fmt.Errorf("failed to get decks: %w", err)
		}

		// get a list of cards with randomized order
		var cards []*models.Card
		for _, deck := range decks {
			cards = append(cards, deck.Cards...)
		}
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})

		// study the cards
		if err := doStudy(cards, scheduler); err != nil {
			return err
		}

		// write the changes to the decks
		for _, deck := range decks {
			err = deckSource.WriteDeck(deck)
			if err != nil {
				return fmt.Errorf("failed to sync studied deck %q: %w", deck.Name, err)
			}
		}

		return nil
	},
}

// Runs the study TUI until there is an error, we run out of cards
// to review, the user quits, or the user edits a card. If the user
// chooses to edit a card, it allows them to do so, and then resumes
// studying by calling itself. Uses recursion because it is the
// cleanest solution for re-evaluating the set of cards that need to
// be studied each time the user edits a card.
func doStudy(cards []*models.Card, scheduler scheduler.Scheduler) error {
	// get only cards that are due to be studied
	cardsToStudy := make([]*models.Card, 0, len(cards))
	for _, card := range cards {
		isDue, err := scheduler.IsDue(card)
		if err != nil {
			return fmt.Errorf("failed to determine whether card %q is due: %w", card.ID, err)
		}
		if isDue && card.Active {
			cardsToStudy = append(cardsToStudy, card)
		}
	}

	if cardID, err := doStudyFragment(cardsToStudy, scheduler); errors.Is(err, views.ErrExit) {
		return nil
	} else if errors.Is(err, views.ErrEdit) {
		card, err := getCardByID(cardID, cardsToStudy)
		if err != nil {
			return fmt.Errorf("failed to get card from cards to study: %w", err)
		}
		if err := models.EditCardViaEditor(card); err != nil && !errors.Is(err, models.ErrNotModified) {
			return fmt.Errorf("failed to edit card %q: %w", cardID, err)
		}
		if err := doStudy(cardsToStudy, scheduler); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("problem while studying cards: %w", err)
	}
	return nil
}

// This is a separate function because it allows screen.Fini() to be
// called as a deferred function.
func doStudyFragment(cardsToStudy []*models.Card, scheduler scheduler.Scheduler) (string, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return "", fmt.Errorf("failed to instantiate Screen: %w", err)
	}
	defer screen.Fini()
	if err := screen.Init(); err != nil {
		return "", fmt.Errorf("failed to initialize Screen: %w", err)
	}
	ss := &views.StudySession{
		Screen:    screen,
		Cards:     cardsToStudy,
		Scheduler: scheduler,
	}
	return ss.Run()
}

func getCardByID(cardID string, cards []*models.Card) (*models.Card, error) {
	for _, card := range cards {
		if card.ID == cardID {
			return card, nil
		}
	}
	return nil, fmt.Errorf("failed to find card with ID %q in passed card slice", cardID)
}
