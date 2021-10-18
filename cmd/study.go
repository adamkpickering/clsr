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
	"strings"
	"time"

	"github.com/adamkpickering/clsr/models"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

// studyCmd represents the study command
var studyCmd = &cobra.Command{
	Use:   "study",
	Short: "Study cards that are due",
	Long:  "Used to study any cards that need to be studied.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// get a DeckSource
		deckSource, err := models.NewFlatFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to get DeckSource: %s", err)
		}

		// get a list of cards to study
		deck, err := deckSource.LoadDeck(deckName)
		if err != nil {
			return fmt.Errorf("failed to load deck: %s", err)
		}
		now := time.Now().UTC()
		cardsToStudy := []*models.Card{}
		for _, card := range deck.Cards {
			if card.NextReview.Before(now) {
				cardsToStudy = append(cardsToStudy, card)
			}
		}

		// initialize tcell
		screen, err := tcell.NewScreen()
		if err != nil {
			return fmt.Errorf("failed to instantiate Screen: %s", err)
		}
		defer screen.Fini()
		if err := screen.Init(); err != nil {
			return fmt.Errorf("failed to initialize Screen: %s", err)
		}

		// study cards
		ss := &StudySession{
			screen: screen,
			cards:  cardsToStudy,
		}
		err = ss.Run()
		if err != nil && err != ErrExit {
			return fmt.Errorf("problem while studying cards: %s", err)
		}

		// write the changes to the deck
		err = deckSource.SyncDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to sync studied deck \"%s\": %s", deck.Name, err)
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
	studyCmd.Flags().StringVarP(&deckName, "deck", "d", "", "the deck to study")
	studyCmd.MarkFlagRequired("deck")
}

var ErrExit error = errors.New("exit study session")

var StyleDefault tcell.Style

type studyState int

const (
	questionState studyState = iota
	questionAndAnswerState
)

type StudySession struct {
	screen tcell.Screen
	cards  []*models.Card
}

func (ss StudySession) Run() error {
	for _, card := range ss.cards {
		err := ss.studyCard(card)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss StudySession) studyCard(card *models.Card) error {
	state := questionState
	for {
		// render screen
		ss.screen.Clear()
		switch state {
		case questionState:
			ss.renderQuestionOnly(card)
		case questionAndAnswerState:
			ss.renderQuestionAndAnswer(card)
		}
		ss.screen.Show()

		// poll for event and act on it
		eventInterface := ss.screen.PollEvent()
		switch event := eventInterface.(type) {
		case *tcell.EventResize:
			ss.screen.Sync()
		case *tcell.EventKey:
			key := event.Key()

			// allow user to exit cleanly and prematurely
			if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
				return ErrExit
			}

			// handle different keys depending on different states
			switch state {
			case questionState:
				if key == tcell.KeyRune {
					if keyRune := event.Rune(); key == tcell.KeyEnter || keyRune == ' ' {
						state = questionAndAnswerState
					}
				}
			case questionAndAnswerState:
				if key == tcell.KeyRune {
					keyRune := event.Rune()
					multiplier, err := getMultiplierFromRune(keyRune)
					if err != nil {
						continue
					}
					card.SetNextReview(multiplier)
					return nil
				}
			}
		}
	}
}

func (ss StudySession) processString(rawString string) []string {
	trimmedString := strings.TrimSpace(rawString)
	stringLines := strings.Split(trimmedString, "\n")
	return stringLines
}

func (ss StudySession) renderQuestionOnly(card *models.Card) {
	lines := ss.processString(card.Question)
	for lineIndex, line := range lines {
		for i, runeValue := range line {
			ss.screen.SetContent(i, lineIndex, runeValue, nil, StyleDefault)
		}
	}
}

func (ss StudySession) renderQuestionAndAnswer(card *models.Card) {
	// build lines var, which represents lines to print to screen
	lines := []string{}
	for _, questionLine := range ss.processString(card.Question) {
		lines = append(lines, questionLine)
	}
	lines = append(lines, "\n")
	lines = append(lines, "------")
	lines = append(lines, "\n")
	for _, answerLine := range ss.processString(card.Answer) {
		lines = append(lines, answerLine)
	}

	// print to screen
	for lineIndex, line := range lines {
		for i, runeValue := range line {
			ss.screen.SetContent(i, lineIndex, runeValue, nil, StyleDefault)
		}
	}
}

func getMultiplierFromRune(key rune) (float64, error) {
	valueMap := map[rune]float64{
		'1': 0.0,
		'2': 1.0,
		'3': 1.5,
		'4': 2.0,
	}
	multiplier, ok := valueMap[key]
	if ok {
		return multiplier, nil
	}
	return 0.0, fmt.Errorf("invalid key \"%c\"", key)
}
