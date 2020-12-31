package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	config   Configuration
	verbose  bool
	logger   = log.New()
	defaults = map[string]interface{}{
		"dbName":     "expenses.db",
		"disableDb":  false,
		"csvName":    "expenses.csv",
		"disableCSV": false,
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
		"disableResult": false,
	}
)

type Configuration struct {
	DbName        string   `mapstructure:"dbName"`
	DisableDb     bool     `mapstructure:"disableDb"`
	CsvName       string   `mapstructure:"csvName"`
	DisableCSV    bool     `mapstructure:"disableCSV"`
	Categories    []string `mapstructure:"categories"`
	Funds         []string `mapstructure:"funds"`
	DisableResult bool     `mapstructure:"disableResult"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "expenses",
	Short: "A simple command line utility to log your expenses",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.SetOutput(os.Stdout)
		if verbose {
			logger.SetLevel(log.DebugLevel)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	logger.SetLevel(log.InfoLevel)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.expenses.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "use verbose logging")
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
		logger.Debugf("Using home directory: %s", home)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".expenses" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".expenses")
	}

	err := viper.ReadInConfig()

	// If a config file is found, read it in.
	if err == nil {
		logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Debug("Config file not found, using default values instead")
		} else {
			// Config file was found but another error was produced
			logger.Fatalf("Some error occurred, cannot proceed: %v", err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Fatal("Could not decode viper configuration struct: %v", err)
	}

}
