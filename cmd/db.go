package cmd

import (
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Display contents from your expenses database",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.DisableDb {
			// TODO
		}
		logger.Errorf("Database has been disabled, please re-enable and run this command again.")
	},
}

func init() {
	showCmd.AddCommand(dbCmd)
}
