package cmd

import (
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show contents of the CSV or SQLite database created by the program",
}

func init() {
	rootCmd.AddCommand(showCmd)
}
