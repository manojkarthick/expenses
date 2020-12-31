package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the current configuration used by expenses",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Print("Current configuration: \n\n")

		fmt.Printf("Disable Database: %t\n", config.DisableDb)
		if !config.DisableDb {
			fmt.Printf("Database Name: %s\n", config.DbName)
		}

		fmt.Printf("Disable CSV: %t\n", config.DisableCSV)
		if !config.DisableCSV {
			fmt.Printf("CSV File Name: %s\n", config.CsvName)
		}

		fmt.Printf("Disable Result: %t\n", config.DisableResult)
		fmt.Printf("Categories: %s\n", strings.Join(config.Categories, ", "))
		fmt.Printf("Funds: %s\n", strings.Join(config.Funds, ", "))

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
