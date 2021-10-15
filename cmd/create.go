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

	"github.com/adamkpickering/clsr/models"
	"github.com/spf13/cobra"
)

var deckName string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <resource_type>",
	Short: "Create a resource",
	Long:  "\nAllows the user to create clsr resources.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check that the directory has been initialized
		if _, err := os.Stat(deckDirectory); errors.Is(err, os.ErrNotExist) {
			msg := "could not find %s. Please call `clsr init` and try again."
			return fmt.Errorf(msg, deckDirectory)
		}

		// construct DeckSource
		deckSource, err := models.NewFlatFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to construct DeckSource: %s", err)
		}

		resourceType := args[0]
		switch resourceType {

		case "card":
			// load the deck
			deck, err := deckSource.LoadDeck(deckName)
			if err != nil {
				return fmt.Errorf("failed to load deck %s: %s", deckName, err)
			}

			// create the card
			card, err := getCardViaEditor()
			if err != nil {
				return fmt.Errorf("failed to get card: %s", err)
			}

			// add the card to the deck
			deck.AddCard(card)

			// sync the deck
			err = deckSource.SyncDeck(deck)
			if err != nil {
				return fmt.Errorf("failed to sync deck: %s", err)
			}

		case "deck":
			// create the deck
			_, err := deckSource.CreateDeck(deckName)
			if err != nil {
				return fmt.Errorf("failed to create deck %s: %s", deckName, err)
			}

		default:
			return fmt.Errorf("\"%s\" is not a valid resource type", resourceType)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringVarP(&deckName, "deck", "d", "", "the deck to act on")
	createCmd.MarkFlagRequired("deck")
}

// Execs into the user's editor to allow them to fill a new card out with the
// question and answer they want. Returns the card.
func getCardViaEditor() (models.Card, error) {
	// create temporary directory
	tempDir, err := os.MkdirTemp("", "clsr_")
	if err != nil {
		return models.Card{}, fmt.Errorf("failed to create tempdir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	// get interim card file
	cardPath, err := getInterimCardFile(tempDir)
	if err != nil {
		return models.Card{}, fmt.Errorf("failed to create interim card file: %s", err)
	}

	// call the user's editor to let them edit the card
	cmd := exec.Command("vim", cardPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return models.Card{}, fmt.Errorf("editor command error: %s", err)
	}

	// read what the user wrote in the file
	data, err := os.ReadFile(cardPath)
	if err != nil {
		return models.Card{}, fmt.Errorf("failed to read file: %s", err)
	}
	cardFilename := filepath.Base(cardPath)
	cardID := strings.Split(cardFilename, ".")[0]
	card, err := models.ParseCardFromString(string(data), cardID)
	if err != nil {
		return models.Card{}, fmt.Errorf("failed to parse card: %s", err)
	}

	return card, nil
}

// Returns the path to a file that is filled with a card template so that it
// can be opened by the user, edited, and then read again to turn what
// the user wrote into a valid Card object. Should be removed once you are
// finished with it.
func getInterimCardFile(baseDir string) (string, error) {
	// get an interim card
	question := "Write the question here."
	answer := "Write the answer here."
	card := models.NewCard(question, answer)

	// write to temporary file
	err := card.WriteToDir(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err)
	}

	return filepath.Join(baseDir, fmt.Sprintf("%s.txt", card.ID)), nil
}
