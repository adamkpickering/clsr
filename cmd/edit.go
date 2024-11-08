package cmd

import (
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit resources",
}

func init() {
	rootCmd.AddCommand(editCmd)
}
