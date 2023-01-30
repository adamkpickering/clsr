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
	"strings"

	"github.com/adamkpickering/clsr/pkg/deck_source"
	"github.com/adamkpickering/clsr/pkg/models"
	"github.com/spf13/cobra"
)

var ErrNotModified error = errors.New("temporary file not modified")

const tempFileDivider = "--------------------"

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <resource_type> [<resource_name>]",
	Short: "Create a resource",
	Long:  "\nAllows the user to create clsr resources.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check that the directory has been initialized
		if _, err := os.Stat(deckDirectory); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("Could not find %s. Please invoke `clsr init`.", deckDirectory)
		}

		// construct DeckSource
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to construct DeckSource: %w", err)
		}

		resourceType := args[0]
		switch resourceType {

		case "card":
			// read the deck
			if len(deckName) == 0 {
				return errors.New("--deck or -d is required for this command")
			}
			deck, err := deckSource.ReadDeck(deckName)
			if err != nil {
				return fmt.Errorf("failed to read deck %q: %w", deckName, err)
			}

			// exec into editor to get Card fields from user
			card, err := getCardFromUserViaEditor()
			fmt.Println(card, err)
			if err == ErrNotModified {
				return nil
			} else if err != nil {
				return fmt.Errorf("failed to get user input: %w", err)
			}

			// add the Card to the Deck and write the Deck
			deck.Cards = append(deck.Cards, card)
			fmt.Println(deck)
			err = deckSource.WriteDeck(deck)
			if err != nil {
				return fmt.Errorf("failed to save deck: %w", err)
			}

		case "deck":
			// get deckName
			if len(args) == 2 {
				deckName = args[1]
			}
			if len(deckName) == 0 {
				msg := "must specify deck name either as positional arg or in --deck/-d flag"
				return errors.New(msg)
			}

			// check for an existing deck of this name
			_, err := deckSource.ReadDeck(deckName)
			if err == nil {
				return fmt.Errorf(`deck "%s" already exists`, deckName)
			}

			// create the deck
			deck := models.NewDeck(deckName)
			err = deckSource.WriteDeck(deck)
			if err != nil {
				return fmt.Errorf("failed to write deck %q: %w", deckName, err)
			}

		default:
			return fmt.Errorf("%q is not a valid resource type", resourceType)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

// Returns the editor specified in the EDITOR env var. If EDITOR is not specified,
// or has zero length, defaults to "nano".
func getPreferredEditor() string {
	value, ok := os.LookupEnv("EDITOR")
	if ok && len(value) > 0 {
		return value
	}
	return "nano"
}

func getTempFileContents() string {
	tempFileQuestion := "# Write the question here. All lines starting with # will be ignored."
	tempFileAnswer := "# Write the answer here. All lines starting with # will be ignored."
	return fmt.Sprintf("%s\n%s\n%s\n", tempFileQuestion, tempFileDivider, tempFileAnswer)
}

// Execs into the user's preferred editor and returns what they enter
// as a new Card. If the user exits without writing any changes,
// error is set to ErrNotModified.
func getCardFromUserViaEditor() (*models.Card, error) {
	// create temp directory
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return &models.Card{}, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// write temp file into the temp directory
	initialText := getTempFileContents()
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
	cmd := exec.Command(getPreferredEditor(), tempFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return &models.Card{}, fmt.Errorf("editor error: %w", err)
	}

	// return if the user did not write the temp file
	info, err = os.Stat(tempFilePath)
	if err != nil {
		fmtString := "failed to get temp file info after potential write: %w"
		return &models.Card{}, fmt.Errorf(fmtString, err)
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
	card := models.NewCard(elements[0], elements[1])

	return card, nil
}
