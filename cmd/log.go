/*
Copyright © 2020 Manoj Karthick Selva Kumar <manojkarthick@ymail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"expenses/utils"
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

var database *sql.DB

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log your expenses to the file/database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("log called")

		var transaction TransactionInfo
		transaction.askQuestionsToUser()

		file, err := os.OpenFile(config.CsvName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
		if err != nil {
			log.Fatal("Could not access transactionsFile: ", err)
		}
		transaction.WriteTransactionToCSV(file)
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatalf("Could not close CSV file %s: %v", config.CsvName, err)
			}
		}()

		database, err = sql.Open("sqlite3", config.DbName)
		if err != nil {
			log.Fatalf("Could not open SQLite database %s: %v: ", config.DbName, err)
		}
		transaction.WriteTransactionToDB(database)
		defer func() {
			if err := database.Close(); err != nil {
				log.Fatalf("Could not close SQLite database %s: %v", config.DbName, err)
			}
		}()

	},
}

func init() {
	rootCmd.AddCommand(logCmd)
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
			log.Fatal(err)
		}
	} else if inputDate == DateYesterday {
		result = time.Now().AddDate(0, 0, -1).Format(DateFormat)
	} else if inputDate == DateDayBefore {
		result = time.Now().AddDate(0, 0, -2).Format(DateFormat)
	} else if inputDate == DateToday {
		result = time.Now().Format(DateFormat)
	} else {
		fmt.Println("Could not understand the date.")
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
		fmt.Println(err.Error())
		return
	}

	t.answer = answer
	t.txnDate = GetTransactionDate(answer.Date)
	t.txnId = utils.GetTransactionId()
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
	if err := writer.Write(t.FormatAsString()); err != nil {
		log.Fatal("Could not write to csv.")
	}
	writer.Flush()
}

func (t *TransactionInfo) WriteTransactionToDB(database *sql.DB) {
	_, err := database.Exec(createStatementSQL)
	if err != nil {
		log.Printf("%q: %s\n", err, createStatementSQL)
		return
	}

	txn, err := database.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insertStatement, err := txn.Prepare(insertStatementSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := insertStatement.Close(); err != nil {
			log.Fatalf("Could not complete insert to database %s: %v", config.DbName, err)
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
		log.Fatal(err)
	}
	err = txn.Commit()
	if err != nil {
		log.Fatalf("Could not commit transaction to database %s: %v", config.DbName, err)
	}

}
