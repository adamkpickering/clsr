package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
}

func init() {
	rootCmd.AddCommand(createCmd)
}
