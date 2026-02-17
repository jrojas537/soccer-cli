package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var ApiKey string
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "soccer-cli",
	Short: "A CLI to check soccer scores and game details.",
	Long: `soccer-cli is a command-line interface to interact with the API-Football
service, allowing you to retrieve scores, game details, and squad information
directly from your terminal.`,
	Version: version,
}

func logError(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+message+"\n", args...)
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logError("%v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/soccer-cli/config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		configPath := home + "/.config/soccer-cli"
		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logError("Config file not found. Please create one at ~/.config/soccer-cli/config.yaml\nExample:\napikey: YOUR_API_KEY_HERE")
		} else {
			logError("reading config file: %v", err)
		}
	}

	ApiKey = viper.GetString("apikey")
	if ApiKey == "" {
		logError("API key is missing from the config file.")
	}
}
