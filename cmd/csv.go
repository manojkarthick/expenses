package cmd

import (
	"github.com/manojkarthick/expenses/utils"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Display contents from your expenses CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.DisableCSV {
			records, err := utils.ReadCSVFile(config.CsvName)
			if err != nil {
				logger.Fatalf("Could not read CSV File: %v", err)
			}
			log.Debugf("Read csv file: %s", config.CsvName)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Transaction ID", "Date", "Item", "Cost", "Location", "Category", "Source"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetBorder(true)

			log.Debug("Starting table render")
			for _, record := range records {
				table.Append(record[0:7])
			}
			table.Render()
		} else {
			logger.Warnf("CSV has been disabled, please re-enable and run this command again.")
		}

	},
}

func init() {
	showCmd.AddCommand(csvCmd)
}
