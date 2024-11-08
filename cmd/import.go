package cmd

import (
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import decks from other spaced repetition tools",
}

func init() {
	rootCmd.AddCommand(importCmd)
}
