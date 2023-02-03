/*
Copyright © 2021 ADAM PICKERING

 ermission is hereby granted, free of charge, to any person obtaining a copy
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
	"os"

	"github.com/adamkpickering/clsr/internal/deck_source"
	"github.com/adamkpickering/clsr/internal/models"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var deckDirectory string

var rootCmd = &cobra.Command{
	Use:   "clsr",
	Short: "Learn things efficiently on the CLI using spaced repetition",
	Long: `
clsr is a CLI tool that allows you to manage and study decks
of virtual flash cards. It takes care of scheduling so that you do not
review them more often than necessary. Similar to Anki and other
spaced repetition applications, except cards are always stored in
plain text. By doing this, we gain all the usual benefits of storing
data in plain text, such as scriptability and the ability to commit
data to version control.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	localDeckDirectory, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %w", err))
	}
	deckDirectory = localDeckDirectory

	cobra.OnInitialize(initConfig)
	configHelp := "config file (default is $HOME/.clsr.yaml)"
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", configHelp)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".clsr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".clsr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// If passedDeckNames is empty, reads all decks and returns them as a slice.
// If passedDeckNames is not empty, reads and returns only the decks named in it.
func getDecks(deckSource deck_source.DeckSource, passedDeckNames ...string) ([]*models.Deck, error) {
	var deckNames []string
	if len(passedDeckNames) == 0 {
		var err error
		deckNames, err = deckSource.ListDecks()
		if err != nil {
			return []*models.Deck{}, fmt.Errorf("failed to list decks: %w", err)
		}
	} else {
		deckNames = passedDeckNames
	}

	// read decks
	decks := []*models.Deck{}
	for _, deckName := range deckNames {
		deck, err := deckSource.ReadDeck(deckName)
		if err != nil {
			return []*models.Deck{}, fmt.Errorf("failed to read deck %q: %w", deckName, err)
		}
		decks = append(decks, deck)
	}

	return decks, nil
}
