package main

import (
	"aifmt/cmd"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aifmt",
	Short: "AI-powered code formatting tool",
	Long: `AIFMT is a CLI tool that uses AI to format and improve your code.
It supports multiple programming languages and AI models.`,
}

func main() {
	// Добавляем команды в корневой командный интерфейс
	rootCmd.AddCommand(cmd.FmtCmd)
	rootCmd.AddCommand(cmd.SetTokenCmd)

	// Выполняем корневую команду и обрабатываем возможные ошибки
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}