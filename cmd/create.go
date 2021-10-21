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
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/adamkpickering/clsr/models"
	"github.com/spf13/cobra"
)

var ErrNotModified error = errors.New("temporary file not modified")

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <resource_type> [<resource_name>]",
	Short: "Create a resource",
	Long:  "\nAllows the user to create clsr resources.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check that the directory has been initialized
		if _, err := os.Stat(deckDirectory); errors.Is(err, os.ErrNotExist) {
			msg := "could not find %s. Please call `clsr init` and try again."
			return fmt.Errorf(msg, deckDirectory)
		}

		// construct DeckSource
		deckSource, err := models.NewFlatFileDeckSource(deckDirectory)
		if err != nil {
			return fmt.Errorf("failed to construct DeckSource: %w", err)
		}

		resourceType := args[0]
		switch resourceType {

		case "card":
			// load the deck
			if len(deckName) == 0 {
				return errors.New("--deck or -d is required for this command")
			}

			deck, err := deckSource.LoadDeck(deckName)
			if err != nil {
				return fmt.Errorf("failed to load deck %q: %w", deckName, err)
			}

			// build []byte to show to user in editor
			initialText := bytes.Buffer{}
			comment := `# The lines until the first "---" contain important ` +
				"card metadata.\n# Edit them at your own risk.\n"
			if _, err = initialText.WriteString(comment); err != nil {
				return fmt.Errorf("failed to write comment to initial card text: %w", err)
			}
			question := "# Write the question here. All lines starting with # will be ignored."
			answer := "# Write the answer here. All lines starting with # will be ignored."
			card := models.NewCard(question, answer)
			cardText, err := card.MarshalText()
			if err != nil {
				return fmt.Errorf("failed to marshal card: %w", err)
			}
			if _, err = initialText.Write(cardText); err != nil {
				return fmt.Errorf("failed to write card to initial card text: %w", err)
			}

			// exec into editor to allow user to input card fields
			userText, err := getInputFromUserViaEditor(initialText.Bytes())
			if err == ErrNotModified {
				return nil
			} else if err != nil {
				return fmt.Errorf("failed to get user input: %w", err)
			}

			// parse the returned data into a card
			card = &models.Card{}
			err = card.UnmarshalText(userText)
			if err != nil {
				return fmt.Errorf("failed to parse user-input data: %w", err)
			}

			// add the card to the deck and sync the deck
			deck.AddCard(card)
			err = deckSource.SyncDeck(deck)
			if err != nil {
				return fmt.Errorf("failed to sync deck: %w", err)
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

			// create the deck
			_, err := deckSource.CreateDeck(deckName)
			if err != nil {
				return fmt.Errorf("failed to create deck %q: %w", deckName, err)
			}

		default:
			return fmt.Errorf("%q is not a valid resource type", resourceType)
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
}

// Creates a temporary file that is filled with the data passed in text. Returns
// the path to the file, as well as the time that the file was last modified.
// Should be removed once you are finished with it.
func getInterimCardFile(baseDir string, text []byte) (string, time.Time, error) {
	// write interim file into the temporary directory
	fd, err := os.CreateTemp(baseDir, "*.txt")
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to open temporary file: %w", err)
	}
	defer fd.Close()
	path := fd.Name()

	// write to file
	_, err = fd.Write(text)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to write marshalled card: %w", err)
	}

	// get last modified time of file
	info, err := os.Stat(path)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get temp file info: %w", err)
	}
	lastModified := info.ModTime()

	return path, lastModified, nil
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

// Execs into the user's chosen editor (specified by EDITOR env var,
// default "nano") and allows them to make changes to a block of text.
// The initial text that is displayed to the user is given by inputText.
// The text that they have written is returned in the first argument.
// If the user exits without writing any changes, error is set to
// ErrNotModified.
func getInputFromUserViaEditor(initialText []byte) ([]byte, error) {
	// create temporary directory
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return []byte{}, fmt.Errorf("failed to create tempdir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// get interim file
	path, firstModified, err := getInterimCardFile(tempDir, initialText)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to create interim card file: %w", err)
	}

	// call the user's editor to let them edit the card
	cmd := exec.Command(getPreferredEditor(), path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return []byte{}, fmt.Errorf("editor command error: %w", err)
	}

	// read what the user wrote in the file if they modified it
	info, err := os.Stat(path)
	if err != nil {
		fmtString := "failed to get temp file info after potential write: %w"
		return []byte{}, fmt.Errorf(fmtString, err)
	}
	if info.ModTime().After(firstModified) {
		data, err := os.ReadFile(path)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to read file: %w", err)
		}
		return data, nil
	}
	return []byte{}, ErrNotModified
}
