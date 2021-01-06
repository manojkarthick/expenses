package cmd

import (
	"fmt"
	"github.com/manojkarthick/expenses/utils"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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
			left := tablewriter.ALIGN_LEFT
			right := tablewriter.ALIGN_RIGHT
			table.SetColumnAlignment([]int{left, left, left, right, left, left, left, left, left})
			table.SetBorder(true)

			var total float64
			log.Debug("Starting table render")
			for _, record := range records {
				table.Append(record[0:7])
				cost, err := strconv.ParseFloat(record[3], 64)
				if err != nil {
					logger.Fatalf("Unable to parse cost value: %s", record[3])
				}
				total += cost
			}
			if showTotal {
				// Add footer
				table.SetFooter([]string{"", "", "Total", fmt.Sprintf("%.0f", total), "", "", ""})
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
