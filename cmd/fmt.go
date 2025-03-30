package cmd

import (
	"aifmt/internal/service"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var FmtCmd = &cobra.Command{
	Use:   "fmt [flags] [files...]",
	Short: "Format code using AI",
	Long: `Format one or more code files using AI.
You can specify the programming language and AI model to use.
If no token is configured, you'll be prompted to set one.`,
	Example: `  # Format a Go file
  aifmt fmt -l go main.go
  
  # Format multiple Python files with a specific model
  aifmt fmt -l python --model claude-2 *.py
  
  # Format with default language detection
  aifmt fmt script.js`,
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("api_key")
		if token == "" {
			fmt.Println("No API token configured. Please run 'aifmt settoken' first.")
			os.Exit(1)
		}

		language, _ := cmd.Flags().GetString("language")

		if language == "" {
			fmt.Println("Error: arg language not set")
			os.Exit(1)
		}

		model, _ := cmd.Flags().GetString("model")

		if len(args) == 0 {
			fmt.Println("Error: no files specified")
			cmd.Help()
			os.Exit(1)
		}

		for _, pattern := range args {
			files, err := filepath.Glob(pattern)
			if err != nil {
				fmt.Printf("Error expanding pattern %s: %v\n", pattern, err)
				continue
			}

			for _, file := range files {
				fmt.Printf("Processing %s (Language: %s, Model: %s)...\n", file, language, model)

				content, err := os.ReadFile(file)
				if err != nil {
					fmt.Printf("Error read file %s: %v\n", file, err)
					continue
				}

				u, upds, err := service.FormatCode(string(content), language, model, token)
				if err != nil {
					fmt.Printf("Error formatting %s: %v\n", file, err)
					continue
				}
				if u == "" {
					fmt.Printf("Error: ai response is empty\n")
					continue
				}

				for _, upd := range upds {
					fmt.Printf("%s\n%s\n\n", upd.Code, upd.Description)
				}

				if err := os.WriteFile(file, []byte(u), 0644); err != nil {
					fmt.Printf("Error write %s: %v\n", file, err)
					continue
				}

				fmt.Printf("Successfully update %s\n", file)
			}
		}
	},
}

func init() {
	FmtCmd.Flags().StringP("language", "l", "", "Programming language of the files")
	FmtCmd.Flags().StringP("model", "m", "deepseek/deepseek-chat:free", "AI model to use for formatting")
}
