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
	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

func init() {
	importCmd.AddCommand(importAnkiCmd)
}

type ankiExportFileHeaders struct {
	Separator  string
	HTML       bool
	DeckColumn int
}

var importAnkiCmd = &cobra.Command{
	Use:   "anki <path>",
	Short: "Import from Anki",
	Long: `Imports decks from Anki. To export decks from Anki, click on
File > Export. Ensure you have the following things set:

- Export format: "Notes in Plain Text"
- Include: "All Decks"
- Only the "Include deck name" checkbox is checked

Then click Export and save the file. Pass the path to this file
to this command and clsr will do the rest.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// read anki export file and split into lines
		ankiDeckPath := args[0]
		data, err := os.ReadFile(ankiDeckPath)
		if err != nil {
			return fmt.Errorf("failed to read anki data: %w", err)
		}
		trimmedData := strings.TrimSpace(string(data))
		allLines := strings.Split(trimmedData, "\n")

		// parse headers
		headers, dataLines, err := parseHeaderLines(allLines)
		if err != nil {
			return fmt.Errorf("failed to parse header lines: %w", err)
		}
		if headers.Separator != "tab" {
			return fmt.Errorf("%q is not a valid value for separator", headers.Separator)
		}
		if headers.DeckColumn != 1 {
			return fmt.Errorf("%q is not a valid value for deck column", headers.DeckColumn)
		}

		// parse data lines into cards
		cards := make([]*models.Card, 0, len(dataLines))
		for _, line := range dataLines {
			parts := strings.Split(line, "\t")
			if len(parts) != 3 {
				fmt.Printf("line could not be parsed: %q\n", line)
				continue
			}
			card := models.NewCard(parts[1], parts[2], parts[0])
			cards = append(cards, card)
		}

		// get deck source
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}

		// check whether cards' decks already exist
		cardDeckNames := map[string]struct{}{}
		for _, card := range cards {
			cardDeckNames[card.Deck] = struct{}{}
		}
		existingDeckNames, err := deckSource.ListDecks()
		if err != nil {
			return fmt.Errorf("failed to list decks: %w", err)
		}
		for _, existingDeckName := range existingDeckNames {
			_, present := cardDeckNames[existingDeckName]
			if present {
				return fmt.Errorf("deck %q already exists", existingDeckName)
			}
		}

		// create decks
		deckNameToDeck := map[string]*models.Deck{}
		for _, card := range cards {
			deck, ok := deckNameToDeck[card.Deck]
			if !ok {
				deck = models.NewDeck(card.Deck, true)
				deckNameToDeck[card.Deck] = deck
			}
			deck.Cards = append(deck.Cards, card)
		}

		// write decks
		for _, deck := range deckNameToDeck {
			err := deckSource.WriteDeck(deck)
			if err != nil {
				return fmt.Errorf("failed to write deck: %w", err)
			}
		}

		return nil
	},
}

// Parses header lines. Returns the headers and the part of the
// file that is not headers.
func parseHeaderLines(lines []string) (ankiExportFileHeaders, []string, error) {
	headers := ankiExportFileHeaders{}
	for i, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return headers, lines[i:], nil
		}
		switch parts[0] {
		case "#separator":
			headers.Separator = parts[1]
		case "#html":
			value, err := strconv.ParseBool(parts[1])
			if err != nil {
				return headers, []string{}, fmt.Errorf("failed to parse html header %q as bool: %w", parts[1], err)
			}
			headers.HTML = value
		case "#deck column":
			value, err := strconv.ParseInt(parts[1], 10, strconv.IntSize)
			if err != nil {
				return headers, []string{}, fmt.Errorf("failed to parse deck column %q as int: %w", parts[1], err)
			}
			headers.DeckColumn = int(value)
		default:
			return headers, lines[i:], nil
		}
	}
	return headers, []string{}, errors.New("reached end of lines with every line matching")
}

// Tells the caller whether slice has an element that is equal
// to element.
func contains(slice []string, element string) bool {
	m := map[string]struct{}{}
	for _, value := range slice {
		m[value] = struct{}{}
	}
	_, present := m[element]
	return present
}
