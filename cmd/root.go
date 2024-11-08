package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deckDirectory string

var rootCmd = &cobra.Command{
	Use:   "clsr",
	Short: "Learn things efficiently on the CLI using spaced repetition",
	Long: `clsr allows you to manage and study decks of virtual flash cards.
It schedules cards according to the principle of spaced repetition
so that you learn most efficiently.`,
}

func init() {
	defaultDeckDirectory, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %w", err))
	}
	rootCmd.PersistentFlags().StringVarP(&deckDirectory, "data-directory", "p", defaultDeckDirectory, "Path to the data directory")
	rootCmd.PersistentFlags().Lookup("data-directory").DefValue = ""
}

func Execute() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	cobra.CheckErr(rootCmd.Execute())
}

func SetVersionInfo(version string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{ .Version }}\n")
}
