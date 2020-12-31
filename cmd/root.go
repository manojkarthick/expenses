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
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	config   Configuration
	defaults = map[string]interface{}{
		"db":       "expenses.db",
		"csv":      "expenses.csv",
		"currency": "$",
		"categories": []string{
			"Rent/Mortgage",
			"Food",
			"Utilities",
			"Maintenance",
			"Living",
			"Health",
			"Electronics",
			"Hygiene",
			"Travel",
			"Education",
		},
		"funds": []string{
			"VISA",
			"Mastercard",
			"Chequing",
			"Savings",
			"Cash",
			"PayPal",
		},
	}
)

type Configuration struct {
	DbName     string   `mapstructure:"db"`
	CsvName    string   `mapstructure:"csv"`
	Currency   string   `mapstructure:"currency"`
	Categories []string `mapstructure:"categories"`
	Funds      []string `mapstructure:"funds"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "expenses",
	Short: "A simple command line utility to log your expenses",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.expenses.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".expenses" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".expenses")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()

	// If a config file is found, read it in.
	if err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using default values instead")
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Some error occurred, cannot proceed: %v", err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Could not decode viper configuration struct")
	}

}
