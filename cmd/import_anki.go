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
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	importCmd.AddCommand(importAnkiCmd)
}

var importAnkiCmd = &cobra.Command{
	Use:   "anki",
	Short: "Import from Anki",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := "test_deck"

		// read file
		ankiDeckPath := args[0]
		data, err := os.ReadFile(ankiDeckPath)
		if err != nil {
			return fmt.Errorf("failed to read anki data: %w", err)
		}

		// process file into lines
		trimmedData := strings.TrimSpace(string(data))
		lines := strings.Split(trimmedData, "\n")
		dataLines := lines[2:]

		// split lines
		splitLines := make([][]string, 0, len(dataLines))
		for _, line := range dataLines {
			splitLine := strings.Split(line, "\t")
			if len(splitLine) != 2 {
				fmt.Printf("found problem line: %q\n", line)
			} else {
				splitLines = append(splitLines, splitLine)
			}
		}

		// parse lines into cards
		cards := make([]*models.Card, 0, len(splitLines))
		for _, splitLine := range splitLines {
			card := models.NewCard(splitLine[0], splitLine[1], deckName)
			cards = append(cards, card)
		}

		// get deck source
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to get deck source: %w", err)
		}

		// create deck and write it
		deck := models.NewDeck(deckName)
		deck.Cards = cards
		err = deckSource.WriteDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to write deck: %w", err)
		}

		return nil
	},
}
