package cmd

import (
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set attributes of resources",
}

func init() {
	rootCmd.AddCommand(setCmd)
}
