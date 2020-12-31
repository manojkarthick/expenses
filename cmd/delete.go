package cmd

import (
	"database/sql"
	"encoding/csv"
	"expenses/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var txnIds []string

const deleteStatementSQL = `
	DELETE FROM expenses WHERE txnId = ?
`

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete expenses by transaction id",
	Run: func(cmd *cobra.Command, args []string) {
		//if no transaction ids passed, ignore
		if len(txnIds) == 0 {
			logger.Warn("No transaction IDs passed, ignored.")
		} else {
			// some transaction ids passed
			logger.Infof("Deleting transaction IDs: %s", strings.Join(txnIds, ", "))

			// if csv is not disabled, proceed
			if !config.DisableCSV {
				logger.Debug("###############################################")
				logger.Debug("Starting CSV operations..")
				records, err := utils.ReadCSVFile(config.CsvName)
				if err != nil {
					logger.Fatalf("Could not read CSV file: %s", config.CsvName)
				}
				logger.Debugf("Read csv file: %s", config.CsvName)

				file, err := os.OpenFile(config.CsvName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
				if err != nil {
					logger.Fatalf("Could not access transactions file: %v", err)
				}
				writer := csv.NewWriter(file)

				var txnId string
				rowsAffected := 0
				for _, record := range records {
					txnId = record[0]
					if utils.Contains(txnIds, txnId) {
						logger.Debugf("Found transaction ID: %s. Skipping write.", txnId)
						rowsAffected += 1
					} else {
						result := record
						if err := writer.Write(result); err != nil {
							logger.Fatalf("Could not write to csv: %v", err)
						}
					}
				}
				writer.Flush()
				logger.Debugf("Successfully wrote to CSV file: %s", config.CsvName)
				logger.Infof("%d rows affected in csv %s", rowsAffected, config.CsvName)

			} else {
				// csv is disabled, ignore and warn
				logger.Warnf("CSV has been disabled. Ignoring")
			}

			// if database is not disabled, proceed
			if !config.DisableDb {
				logger.Debug("###############################################")
				logger.Debug("Starting DB operations..")

				database, err := sql.Open("sqlite3", config.DbName)
				if err != nil {
					logger.Fatalf("Could not open SQLite database %s: %v: ", config.DbName, err)
				}
				var rowsAffected int64 = 0
				for _, txnId := range txnIds {
					txn, err := database.Begin()
					if err != nil {
						logger.Fatal(err)
					}

					deleteStatement, err := txn.Prepare(deleteStatementSQL)
					if err != nil {
						logger.Fatal(err)
					}
					if err := deleteStatement.Close(); err != nil {
						logger.Fatalf("Could not complete delete to database %s: %v", config.DbName, err)
					}

					rs, err := deleteStatement.Exec(txnId)
					if err != nil {
						logger.Fatal(err)
					}
					err = txn.Commit()
					if err != nil {
						logger.Fatalf("Could not commit transaction to database %s: %v", config.DbName, err)
					}
					localRowsAffected, err := rs.RowsAffected()
					if err != nil {
						logger.Fatal(err)
					}

					if localRowsAffected == 1 {
						logger.Debugf("Successfully committed transaction to database: %s to delete: %s", config.DbName, txnId)
					} else {
						logger.Debugf("Could not find transaction with ID: %s. Skipping.", txnId)
					}
					rowsAffected += localRowsAffected
				}

				defer func() {
					if err := database.Close(); err != nil {
						logger.Fatalf("Could not close SQLite database %s: %v", config.DbName, err)
					}
				}()

				logger.Infof("%d rows affected in csv %s", rowsAffected, config.DbName)

			} else {
				// database is disabled, ignore and warn
				logger.Warnf("Database has been disabled. Ignoring")
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVar(&txnIds, "transactions", []string{}, "Transaction IDs to delete")
}
