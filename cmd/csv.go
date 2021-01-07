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
			table.SetHeader([]string{"Transaction ID", "Date", "Location", "Category", "Source", "Item", "Cost"})
			left := tablewriter.ALIGN_LEFT
			right := tablewriter.ALIGN_RIGHT
			table.SetColumnAlignment([]int{left, left, left, left, left, left, right})
			table.SetBorder(true)

			var total float64
			log.Debug("Starting table render")
			for _, record := range records {
				table.Append(reorderRecord(record[0:7]))
				cost, err := strconv.ParseFloat(record[6], 64)
				if err != nil {
					logger.Fatalf("Unable to parse cost value: %s", record[3])
				}
				total += cost
			}
			if showTotal {
				// Add footer
				table.SetFooter([]string{"", "", "", "", "", "Total", fmt.Sprintf("%.2f", total)})
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

// reorderRecord will essentially move "item" and "cost" as the last two elements
func reorderRecord(record []string) []string {
	item, cost := record[2], record[3]
	copy(record[2:], record[4:])
	record[5], record[6] = item, cost

	return record
}
