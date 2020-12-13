package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"path/filepath"
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

type Configuration struct {
	DbName string `json:"db"`
	CsvName string `json:"csv"`
}

const (
	ConfigFile          = ".expense-logger.json"
	DefaultExpensesFile = "expenses.csv"
	DefaultExpensesDB   = "./expenses.db"

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

	// parse config
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("could not read user's home directory. exiting")
		os.Exit(1)
	}

	configFilePath := filepath.Join(home, ConfigFile)

	var expensesDB string
	var expensesFile string
	if _, err := os.Stat(configFilePath); err == nil {
		file, _ := os.Open(configFilePath)
		defer file.Close()

		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		decodeErr := decoder.Decode(&configuration)
		if decodeErr != nil {
			fmt.Errorf("could not decode configuration file %s", configFilePath)
			os.Exit(1)
		}

		expensesDB = configuration.DbName
		expensesFile = configuration.CsvName

	} else if os.IsNotExist(err) {
		expensesDB = DefaultExpensesDB
		expensesFile = DefaultExpensesFile
	} else {
		fmt.Errorf("unknown error encountered. exiting")
		os.Exit(1)
	}

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

	err = survey.Ask(expenseQuestions, &answer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	transaction := TransactionInfo{
		answer:  answer,
		txnId:   GetTransactionId(),
		txnDate: GetTransactionDate(answer.Date),
	}

	tableData := make([][]string, 8)
	tableData = append(tableData, []string{"ID", transaction.txnId})
	tableData = append(tableData, []string{"Date", transaction.txnDate})
	tableData = append(tableData, []string{"Item", transaction.answer.Item})
	tableData = append(tableData, []string{"Cost", fmt.Sprintf("%f", transaction.answer.Cost)})
	tableData = append(tableData, []string{"Location", transaction.answer.Location})
	tableData = append(tableData, []string{"Category", transaction.answer.Category})
	tableData = append(tableData, []string{"Source", transaction.answer.Source})
	tableData = append(tableData, []string{"Notes", transaction.answer.Notes})

	// print the answers
	fmt.Printf("\n\n")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Info", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	for _, v := range tableData {
		table.Append(v)
	}
	table.Render()


	file, err := os.OpenFile(expensesFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		log.Fatal("Could not access transactionsFile: ", err)
	}
	WriteTransactionToCSV(file, &transaction)
	defer file.Close()

	database, err = sql.Open("sqlite3", expensesDB)
	if err != nil {
		log.Fatal("Could not open sqlite database: ", err)
	}
	WriteTransactionToDB(database, &transaction)
	defer database.Close()

}
