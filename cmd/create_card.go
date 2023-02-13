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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/spf13/cobra"
)

const (
	tempFileQuestion = "# Write the question here. This line, as well as the divider below, will be removed.\n"
	tempFileDivider  = "--------------------\n"
	tempFileAnswer   = "# Write the answer here. This line, as well as the above divider, will be removed.\n"
)

var ErrNotModified error = errors.New("temporary file not modified")

var createCardFlags = struct {
	DeckName string
}{}

func init() {
	createCmd.AddCommand(createCardCmd)
	createCardCmd.Flags().StringVarP(&createCardFlags.DeckName, "deck", "d", "", "filter cards by deck")
	createCardCmd.MarkFlagRequired("deck")
}

var createCardCmd = &cobra.Command{
	Use:   "card",
	Short: "Create card",
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := createCardFlags.DeckName
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}

		// read the deck
		deck, err := deckSource.ReadDeck(deckName)
		if err != nil {
			return fmt.Errorf("failed to read deck %q: %w", deckName, err)
		}

		// exec into editor to get Card fields from user
		card, err := getCardFromUserViaEditor(deckName)
		if err == ErrNotModified {
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to get user input: %w", err)
		}

		// add the Card to the Deck and write the Deck
		deck.Cards = append(deck.Cards, card)
		err = deckSource.WriteDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to save deck: %w", err)
		}

		return nil
	},
}

// Returns the editor specified in the EDITOR env var. If EDITOR is not specified,
// or has zero length, defaults to "nano".
func getPreferredEditor() (string, error) {
	value, ok := os.LookupEnv("EDITOR")
	if (!ok || len(value) == 0) && runtime.GOOS == "windows" {
		return "", errors.New("EDITOR environment variable is not set")
	}
	if ok && len(value) > 0 {
		return value, nil
	}
	return "nano", nil
}

// Execs into the user's preferred editor and returns what they enter
// as a new Card. If the user exits without writing any changes,
// error is set to ErrNotModified.
func getCardFromUserViaEditor(deckName string) (*models.Card, error) {
	// create temp directory
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return &models.Card{}, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// write temp file into the temp directory
	initialText := fmt.Sprintf("%s%s%s", tempFileQuestion, tempFileDivider, tempFileAnswer)
	tempFilePath := filepath.Join(tempDir, "clsr_create_card.txt")
	err = os.WriteFile(tempFilePath, []byte(initialText), 0644)
	if err != nil {
		return &models.Card{}, fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tempFilePath)

	// get last modified time of temp file
	info, err := os.Stat(tempFilePath)
	if err != nil {
		return &models.Card{}, fmt.Errorf("failed to get temp file info: %w", err)
	}
	firstModified := info.ModTime()

	// call the user's editor to let them edit the card
	editor, err := getPreferredEditor()
	if err != nil {
		return &models.Card{}, fmt.Errorf("failed to get editor: %w", err)
	}
	cmd := exec.Command(editor, tempFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return &models.Card{}, fmt.Errorf("editor error: %w", err)
	}

	// return if the user did not write the temp file
	info, err = os.Stat(tempFilePath)
	if err != nil {
		return &models.Card{}, fmt.Errorf("failed to get temp file info after potential write: %w", err)
	}
	if !info.ModTime().After(firstModified) {
		return &models.Card{}, ErrNotModified
	}

	// read the contents of the temp file and parse into a Card
	contents, err := os.ReadFile(tempFilePath)
	if err != nil {
		return &models.Card{}, fmt.Errorf("failed to read temp file: %w", err)
	}
	elements := strings.Split(string(contents), tempFileDivider)
	if len(elements) != 2 {
		return &models.Card{}, fmt.Errorf(`splitting on "%s" did not produce exactly 2 elements`, tempFileDivider)
	}
	question := strings.TrimSpace(strings.ReplaceAll(elements[0], tempFileQuestion, ""))
	answer := strings.TrimSpace(strings.ReplaceAll(elements[1], tempFileAnswer, ""))
	card := models.NewCard(question, answer, deckName)

	return card, nil
}
