package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage soccer-cli configuration.",
	Long:  `Provides subcommands to view and set configuration properties for soccer-cli.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration property.",
	Long:  `Sets a configuration property (e.g., 'apikey') in the config file.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if key == "apikey" {
			viper.Set(key, value) // viper.Set (not vipher.Set)
			err := viper.WriteConfig()
			if err != nil {
				logError("writing config file: %v", err)
			}
			fmt.Printf("Successfully set %s.\n", key)
		} else {
			logError("Unsupported config key: %s", key)
		}
	},
}
