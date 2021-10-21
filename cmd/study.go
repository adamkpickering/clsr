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

	"github.com/adamkpickering/clsr/models"
	"github.com/adamkpickering/clsr/views"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

// studyCmd represents the study command
var studyCmd = &cobra.Command{
	Use:   "study",
	Short: "Study cards that are due",
	Long:  "\nUsed to study any cards that need to be studied.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// get a DeckSource
		deckSource, err := models.NewFlatFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to get DeckSource: %w", err)
		}

		// get a list of cards to study
		if len(deckName) == 0 {
			return errors.New("--deck or -d is required for this command")
		}
		deck, err := deckSource.LoadDeck(deckName)
		if err != nil {
			return fmt.Errorf("failed to load deck: %w", err)
		}
		cardsToStudy := []*models.Card{}
		for _, card := range deck.Cards {
			if card.IsDue() && card.Active {
				cardsToStudy = append(cardsToStudy, card)
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

		// write the changes to the deck
		err = deckSource.SyncDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to sync studied deck %q: %w", deck.Name, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(studyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// studyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// studyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
