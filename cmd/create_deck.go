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

	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.AddCommand(createDeckCmd)
}

var createDeckCmd = &cobra.Command{
	Use:   "deck <name>",
	Short: "Create deck",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		deckSource, err := deck_source.NewJSONFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to instantiate deck source: %w", err)
		}

		// check for an existing deck of this name
		_, err = deckSource.ReadDeck(deckName)
		if err == nil {
			return fmt.Errorf(`deck %q already exists`, deckName)
		}

		// create the deck
		deck := models.NewDeck(deckName)
		err = deckSource.WriteDeck(deck)
		if err != nil {
			return fmt.Errorf("failed to write deck %q: %w", deckName, err)
		}

		return nil
	},
}
