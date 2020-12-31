package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"time"

	"github.com/manojkarthick/expenses/utils"
	"github.com/AlecAivazis/survey/v2"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type Answer struct {
	Item     string
	Cost     float64
	Location string
	Category string
	Source   string
	Notes    string
	Date     string
}

type TransactionInfo struct {
	answer  Answer
	txnId   string
	txnDate string
}

const (
	// names of fields
	NameItem     = "item"
	NameCost     = "cost"
	NameLocation = "location"
	NameCategory = "category"
	NameSource   = "source"
	NameNotes    = "notes"
	NameDate     = "date"

	// prompt questions
	PromptItem      = "What did you buy?"
	PromptCost      = "Cost?"
	PromptLocation  = "Where did you buy this?"
	PromptCategory  = "Please select the item's categories"
	PromptSource    = "Source of Funds?"
	PromptNotes     = "Any notes?"
	PromptDate      = "Transaction date"
	PromptOtherDate = "Enter the date"

	// date options
	DateToday     = "Today"
	DateYesterday = "Yesterday"
	DateDayBefore = "Day before"
	DateOther     = "Other"

	createStatementSQL = `
	CREATE TABLE IF NOT EXISTS expenses (
		txnId TEXT PRIMARY KEY,
		txnDate TEXT NOT NULL,
		item TEXT NOT NULL,
		cost DOUBLE NOT NULL,
		location TEXT NOT NULL,
		category TEXT NOT NULL,
		source TEXT NOT NULL,
		notes TEXT
	)`

	insertStatementSQL = `
	INSERT INTO expenses (txnId, txnDate, item, cost, location, category, source, notes)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	DateFormat = "2006/01/02"
)

// addCmd represents the log command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Log your expenses to the file/database",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debugf("Running log command")

		var transaction TransactionInfo
		transaction.askQuestionsToUser()

		if !config.DisableCSV {
			file, err := os.OpenFile(config.CsvName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
			if err != nil {
				logger.Fatalf("Could not access transactions file: %v", err)
			}
			transaction.WriteTransactionToCSV(file)
			defer func() {
				if err := file.Close(); err != nil {
					logger.Fatalf("Could not close CSV file %s: %v", config.CsvName, err)
				}
			}()
		}

		if !config.DisableDb {
			database, err := sql.Open("sqlite3", config.DbName)
			if err != nil {
				logger.Fatalf("Could not open SQLite database %s: %v: ", config.DbName, err)
			}
			transaction.WriteTransactionToDB(database)
			defer func() {
				if err := database.Close(); err != nil {
					logger.Fatalf("Could not close SQLite database %s: %v", config.DbName, err)
				}
			}()

		}

		if !config.DisableResult {
			transaction.RenderTransactionTable()
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func GetTransactionDate(inputDate string) string {
	var result string

	if inputDate == DateOther {
		prompt := &survey.Input{Message: PromptOtherDate}
		if err := survey.AskOne(prompt, &result, survey.WithValidator(func(val interface{}) error {
			_, err := time.Parse(DateFormat, val.(string))
			if err != nil {
				return err
			}
			return nil
		}), survey.WithValidator(survey.Required)); err != nil {
			logger.Fatal(err)
		}
	} else if inputDate == DateYesterday {
		result = time.Now().AddDate(0, 0, -1).Format(DateFormat)
	} else if inputDate == DateDayBefore {
		result = time.Now().AddDate(0, 0, -2).Format(DateFormat)
	} else if inputDate == DateToday {
		result = time.Now().Format(DateFormat)
	} else {
		logger.Fatalf("Could not understand the given date")
	}

	return result
}

func (t *TransactionInfo) askQuestionsToUser() {
	var expenseQuestions = []*survey.Question{
		{
			Name: NameItem,
			Prompt: &survey.Input{
				Message: PromptItem,
			},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: NameCost,
			Prompt: &survey.Input{
				Message: PromptCost,
			},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: NameLocation,
			Prompt: &survey.Input{
				Message: PromptLocation,
			},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: NameCategory,
			Prompt: &survey.Select{
				Message: PromptCategory,
				Options: config.Categories,
			},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: NameSource,
			Prompt: &survey.Select{
				Message: PromptSource,
				Options: config.Funds,
			},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: NameNotes,
			Prompt: &survey.Input{
				Message: PromptNotes,
			},
			Transform: survey.Title,
		},
		{
			Name: NameDate,
			Prompt: &survey.Select{
				Message: PromptDate,
				Options: []string{DateToday, DateYesterday, DateDayBefore, DateOther},
			},
			Transform: survey.Title,
			Validate:  survey.Required,
		},
	}

	answer := Answer{}

	err := survey.Ask(expenseQuestions, &answer)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	txnId := utils.GetTransactionId()
	logger.Debugf("Generated transaction ID: %s", txnId)

	t.answer = answer
	t.txnDate = GetTransactionDate(answer.Date)
	t.txnId = txnId
}

func (t *TransactionInfo) FormatAsString() []string {
	cost := strconv.FormatFloat(t.answer.Cost, 'f', -1, 64)
	result := []string{
		t.txnId, t.txnDate, t.answer.Item, cost, t.answer.Location, t.answer.Category, t.answer.Source, t.answer.Notes,
	}
	return result
}

func (t *TransactionInfo) WriteTransactionToCSV(file *os.File) {
	writer := csv.NewWriter(file)
	rowString := t.FormatAsString()
	logger.Debugf("Writing to CSV: %s", rowString)
	if err := writer.Write(rowString); err != nil {
		logger.Fatalf("Could not write to csv: %v", err)
	}
	writer.Flush()
	logger.Debugf("Successfully wrote to CSV file: %s", config.CsvName)
}

func (t *TransactionInfo) WriteTransactionToDB(database *sql.DB) {
	logger.Debugf("Executing Create table statement: %v", createStatementSQL)
	_, err := database.Exec(createStatementSQL)
	if err != nil {
		logger.Fatalf("Could not create database table: %q\n", err)
		return
	}

	txn, err := database.Begin()
	if err != nil {
		logger.Fatal(err)
	}

	insertStatement, err := txn.Prepare(insertStatementSQL)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debugf("Inserted into database: %s", insertStatement)
	defer func() {
		if err := insertStatement.Close(); err != nil {
			logger.Fatalf("Could not complete insert to database %s: %v", config.DbName, err)
		}
	}()

	_, err = insertStatement.Exec(
		t.txnId,
		t.txnDate,
		t.answer.Item,
		t.answer.Cost,
		t.answer.Location,
		t.answer.Category,
		t.answer.Source,
		t.answer.Notes,
	)
	if err != nil {
		logger.Fatal(err)
	}
	err = txn.Commit()
	if err != nil {
		logger.Fatalf("Could not commit transaction to database %s: %v", config.DbName, err)
	}
	logger.Debugf("Successfully committed transaction to database: %s", config.DbName)

}

func (t *TransactionInfo) RenderTransactionTable() {
	tableData := make([][]string, 8)
	tableData = append(tableData, []string{"ID", t.txnId})
	tableData = append(tableData, []string{"Date", t.txnDate})
	tableData = append(tableData, []string{"Item", t.answer.Item})
	tableData = append(tableData, []string{"Cost", fmt.Sprintf("%f", t.answer.Cost)})
	tableData = append(tableData, []string{"Location", t.answer.Location})
	tableData = append(tableData, []string{"Category", t.answer.Category})
	tableData = append(tableData, []string{"Source", t.answer.Source})
	tableData = append(tableData, []string{"Notes", t.answer.Notes})

	fmt.Println("")

	// print the answers
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Info", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	for _, v := range tableData {
		table.Append(v)
	}
	table.Render()
}
