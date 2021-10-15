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
				return fmt.Errorf("failed to get build card: %s", err)
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

func getCardViaEditor() (models.Card, error) {
	return models.NewCard("test question", "test answer"), nil
}
