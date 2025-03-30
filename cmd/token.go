package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SetTokenCmd = &cobra.Command{
	Use:   "settoken",
	Short: "Set your OpenRouter API token",
	Long: `Set or update your OpenRouter API token.
The token will be saved in the configuration file for future use.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter your OpenRouter API token: ")
		var token string
		_, err := fmt.Scanln(&token)
		if err != nil {
			fmt.Println("Error reading token:", err)
			os.Exit(1)
		}

		token = strings.TrimSpace(token)
		if token == "" {
			fmt.Println("Error: token cannot be empty")
			os.Exit(1)
		}

		viper.Set("api_key", token)
		err = viper.WriteConfig()
		if err != nil {
			fmt.Println("Error saving token:", err)
			os.Exit(1)
		}

		fmt.Println("Token saved successfully!")
	},
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	configDir := filepath.Join(home, ".aifmt")
	configPath := filepath.Join(configDir, "config.yaml")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			fmt.Println("Error creating config directory:", err)
			return
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config not found, creating a default one at", configPath)

			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				fmt.Println("Error writing default config:", err)
			}
		} else {
			fmt.Println("Error reading config:", err)
		}
	}
}
