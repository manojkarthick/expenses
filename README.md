# Expenses [![GitHub release](https://img.shields.io/github/release/manojkarthick/expenses.svg)](https://github.com/manojkarthick/expenses/releases/) ![build](https://github.com/manojkarthick/expenses/workflows/release/badge.svg)
An interactive command line expense logger. Answer a series of questions to log your expenses. Currently writes to CSV and SQLite backends.

## Installation
You can download pre-built binaries [here](https://github.com/manojkarthick/expenses/releases).

View the tool in action here:
 
 [![asciicast](https://asciinema.org/a/381989.svg)](https://asciinema.org/a/381989)

### Building from source
You can use the Makefile present in the directory to build the project. Make sure you have Golang v1.14.x installed on your system.

```
git clone https://github.com/manojkarthick/expenses.git
cd expenses
make dev-build
./expenses help
```

### Running
To add a new expense, run `expenses add` and answer the questions that follow. Example below:

```
$ expenses add
? What did you buy? Cookies
? Cost? 2.99
? Where did you buy this? Trader Joes
? Please select the item's categories Food
? Source of Funds? Chequing
? Any notes?
? Transaction date Today

```

### Available Commands
Run `expenses help` to view the list of commands available:

```
A simple command line utility to log your expenses

Usage:
  expenses [command]

Available Commands:
  add         Log your expenses to the file/database
  config      Show the current configuration used by expenses
  delete      Delete expenses by transaction id
  help        Help about any command
  show        Show contents of the CSV or SQLite database created by the program
  version     Show application version information

Flags:
      --config string   config file (default is $HOME/.expenses.yaml)
  -h, --help            help for expenses
      --verbose         use verbose logging

Use "expenses [command] --help" for more information about a command.
```

* To add a new expense: `expenses add`
* View the current version: `expenses version`
* View the configuration used: `expenses config`
* Show the expenses: `expenses show db` or `expenses show csv`
* Delete expenses: `expenses delete --transaction <transaction_ids>`


### Sample Output

To view the expenses currently logged, run `expenses show db` or `expenses show csv`.

```
+--------------------------------------+------------+---------+----------+-------------+----------+----------+
|            TRANSACTION ID            |    DATE    |  ITEM   |   COST   |  LOCATION   | CATEGORY |  SOURCE  |
+--------------------------------------+------------+---------+----------+-------------+----------+----------+
| fb96cf93-b096-4137-ad8b-b60a3bf08045 | 2020/12/31 | Cookies | 2.990000 | Trader Joes | Food     | Chequing |
| 3d275a44-76a2-4162-942d-d8901cae2c82 | 2020/12/30 | Coffee  | 4.990000 | Starbucks   | Food     | Cash     |
+--------------------------------------+------------+---------+----------+-------------+----------+----------+
```

### Configuration

Expenses allows you to configure values such as the database name, csv file name, categories for expenses, source of funds, etc.

Create a file under your home directory called `.expenses.yaml`. You can modify the following fields in the config file:

1. `dbName`: Name of the SQLite3 database to commit to (`default = expenses.db`)
2. `disableDb`: Set to `true` if you don't want to write to database (`default = false`)
3. `csvName`: Name of the CSV file to write to (`default = expenses.csv`)
4. `disableCSV`: Set to `true` if you don't want to write to the CSV file (`default = false`)
5. `categories`: Provide an alternate list of categories for your expenses to show during the interactive prompt (`default = Rent/Mortgage, Food, Utilities, Maintenance, Living, Health, Electronics, Hygiene, Travel, Education`)
6. `funds`: Provide possible source of funds for the expense (`default = VISA, Mastercard, Chequing, Savings, Cash, PayPal`)
 












