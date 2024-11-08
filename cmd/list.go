package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
}

func init() {
	rootCmd.AddCommand(listCmd)
}
