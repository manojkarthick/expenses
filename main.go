package main

import (
	"database/sql"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
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
	ExpensesFile = "expenses.csv"
	ExpensesDB   = "./expenses.db"

	// add all question names here
	NameItem     = "item"
	NameCost     = "cost"
	NameLocation = "location"
	NameCategory = "category"
	NameSource   = "source"
	NameNotes    = "notes"
	NameDate     = "date"

	// add all message prompts here
	PromptItem      = "What did you buy?"
	PromptCost      = "Cost?"
	PromptLocation  = "Where did you buy this?"
	PromptCategory  = "Please select the item's categories"
	PromptSource    = "Source of Funds?"
	PromptNotes     = "Any notes?"
	PromptDate      = "Transaction date"
	PromptOtherDate = "Enter the date"

	// add all expense categories here
	CategoryRent        = "Rent"
	CategoryFood        = "Food"
	CategoryUtilities   = "Utilities"
	CategoryMaintenance = "Maintenance"
	CategoryLiving      = "Living"
	CategoryHealth      = "Health"
	CategoryElectronics = "Electronics"
	CategoryMedicines   = "Medicines"
	CategoryHygiene     = "Hygiene"
	CategoryTravel      = "Travel"

	// add expense sources here
	SourceVisa       = "VISA"
	SourceMastercard = "Mastercard"
	SourceChequing   = "Chequing"
	SourceCash       = "Cash"
	SourcePayPal     = "PayPal"
	SourceGiftCards  = "Gift cards"

	// all top-level date options here
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

func main() {

	// the questions to ask
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
				Options: []string{
					CategoryRent, CategoryFood, CategoryUtilities, CategoryMaintenance, CategoryLiving, CategoryHealth, CategoryElectronics, CategoryMedicines, CategoryHygiene, CategoryTravel,
				},
			},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: NameSource,
			Prompt: &survey.Select{
				Message: PromptSource,
				Options: []string{SourceVisa, SourceMastercard, SourceChequing, SourceCash, SourcePayPal, SourceGiftCards},
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

	transaction := TransactionInfo{
		answer:  answer,
		txnId:   GetTransactionId(),
		txnDate: GetTransactionDate(answer.Date),
	}

	// print the answers
	fmt.Printf("Transaction Date: %s\n", transaction.txnId)
	fmt.Printf("Transaction Date: %s\n", transaction.txnDate)
	fmt.Printf("Item: %s\n", transaction.answer.Item)
	fmt.Printf("Cost: %f\n", transaction.answer.Cost)
	fmt.Printf("Location: %s\n", transaction.answer.Location)
	fmt.Printf("Category: %s\n", transaction.answer.Category)
	fmt.Printf("Source: %s\n", transaction.answer.Source)
	fmt.Printf("Notes: %s\n", transaction.answer.Notes)

	file, err := os.OpenFile(ExpensesFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		log.Fatal("Could not access transactionsFile: ", err)
	}
	WriteTransactionToCSV(file, &transaction)
	defer file.Close()

	database, err = sql.Open("sqlite3", ExpensesDB)
	if err != nil {
		log.Fatal("Could not open sqlite database: ", err)
	}
	WriteTransactionToDB(database, &transaction)
	defer database.Close()

}
