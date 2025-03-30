package main

import (
	"aifmt/cmd"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ErrRootCommand = errors.New("root command execution error")
)

var rootCmd = &cobra.Command{
	Use:   "aifmt",
	Short: "Инструмент для форматирования кода с помощью ИИ",
	Long: `AIFMT - это инструмент командной строки, использующий ИИ для форматирования и улучшения вашего кода.
Поддерживает множество языков программирования и моделей ИИ.`,
}

func main() {
	rootCmd.AddCommand(cmd.FmtCmd)
	rootCmd.AddCommand(cmd.SetCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка выполнения команды: %v\n", err)
		os.Exit(1)
	}
}
