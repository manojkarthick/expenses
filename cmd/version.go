package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "devel"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show application version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: ", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
