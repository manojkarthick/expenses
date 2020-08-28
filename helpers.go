package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
	"time"
)

func FormatAsString(t *TransactionInfo) []string {
	cost := strconv.FormatFloat(t.answer.Cost, 'f', -1, 64)
	result := []string{
		t.txnId, t.txnDate, t.answer.Item, cost, t.answer.Location, t.answer.Category, t.answer.Source, t.answer.Notes,
	}
	return result
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

func GetTransactionId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func WriteTransactionToCSV(file *os.File, t *TransactionInfo) {
	writer := csv.NewWriter(file)
	if err := writer.Write(FormatAsString(t)); err != nil {
		log.Fatal("Could not write to csv.")
	}
	writer.Flush()
}

func WriteTransactionToDB(database *sql.DB, transaction *TransactionInfo) {
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
	defer insertStatement.Close()

	_, err = insertStatement.Exec(transaction.txnId, transaction.txnDate, transaction.answer.Item, transaction.answer.Cost, transaction.answer.Location, transaction.answer.Category, transaction.answer.Source, transaction.answer.Notes)
	if err != nil {
		log.Fatal(err)
	}
	txn.Commit()
}
