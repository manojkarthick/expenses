package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

const selectStatementSQL = `
	SELECT txnId, txnDate, item, cost, location, category, "source" FROM expenses
`

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Display contents from your expenses database",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.DisableDb {
			database, err := sql.Open("sqlite3", config.DbName)
			if err != nil {
				logger.Fatalf("Could not open SQLite database %s: %v: ", config.DbName, err)
			}
			rows, err := database.Query(selectStatementSQL)
			defer func() {
				if err := database.Close(); err != nil {
					logger.Fatal(err)
				}
			}()
			if err != nil {
				logger.Fatalf("Could not read from database: %s", config.DbName)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Transaction ID", "Date", "Item", "Cost", "Location", "Category", "Source"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetBorder(true)

			var txnId string
			var txnDate string
			var item string
			var cost float64
			var location string
			var category string
			var source string
			for rows.Next() {
				rows.Scan(&txnId, &txnDate, &item, &cost, &location, &category, &source)
				table.Append([]string{txnId, txnDate, item, fmt.Sprintf("%f", cost), location, category, source})
			}
			table.Render()

		} else {
			logger.Warnf("Database has been disabled, please re-enable and run this command again.")
		}
	},
}

func init() {
	showCmd.AddCommand(dbCmd)
}
