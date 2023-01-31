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

	"github.com/adamkpickering/clsr/pkg/deck_source"
	"github.com/adamkpickering/clsr/pkg/models"
	"github.com/adamkpickering/clsr/views"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

var studyCmd = &cobra.Command{
	Use:   "study",
	Short: "Study cards that are due",
	Long:  "\nUsed to study any cards that need to be studied.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// get a DeckSource
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to get DeckSource: %w", err)
		}

		// get a list of decks
		decks := []*models.Deck{}
		if len(deckName) > 0 {
			deck, err := deckSource.ReadDeck(deckName)
			if err != nil {
				return fmt.Errorf("failed to read deck %q: %w", deckName, err)
			}
			decks = append(decks, deck)
		} else {
			decks, err = getAllDecks(deckSource)
			if err != nil {
				return fmt.Errorf("failed to get decks: %w", err)
			}
		}

		// get a list of cards to study
		var cardsToStudy []*models.Card
		for _, deck := range decks {
			for _, card := range deck.Cards {
				isDue, err := card.IsDue()
				if err != nil {
					return fmt.Errorf("failed to determine whether card %q is due: %w", card.ID, err)
				}
				if isDue && card.Active {
					cardsToStudy = append(cardsToStudy, card)
				}
			}
		}

		// initialize tcell
		screen, err := tcell.NewScreen()
		if err != nil {
			return fmt.Errorf("failed to instantiate Screen: %w", err)
		}
		defer screen.Fini()
		if err := screen.Init(); err != nil {
			return fmt.Errorf("failed to initialize Screen: %w", err)
		}

		// study cards
		ss := &views.StudySession{
			Screen: screen,
			Cards:  cardsToStudy,
		}
		err = ss.Run()
		if err != nil && err != views.ErrExit {
			return fmt.Errorf("problem while studying cards: %w", err)
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

func init() {
	rootCmd.AddCommand(studyCmd)
}
