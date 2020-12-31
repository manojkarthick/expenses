package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the current configuration used by expenses",
	Run: func(cmd *cobra.Command, args []string) {

		logger.Info("Current configuration:")

		logger.Infof("Disable Database: %t", config.DisableDb)
		if !config.DisableDb {
			logger.Infof("Database Name: %s", config.DbName)
		}

		logger.Infof("Disable CSV: %t", config.DisableCSV)
		if !config.DisableCSV {
			logger.Infof("CSV File Name: %s", config.CsvName)
		}

		logger.Info("Disable Result: %t", config.DisableResult)
		logger.Infof("Categories: %s", strings.Join(config.Categories, ", "))
		logger.Infof("Funds: %s", strings.Join(config.Funds, ", "))

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
