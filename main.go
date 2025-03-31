package main

import (
	"fmt"
	"os"

	"github.com/seelentov/aifmt/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aifmt",
	Short: "Инструмент для форматирования кода с помощью ИИ",
	Long:  `AIFMT - это инструмент командной строки, использующий ИИ для форматирования и улучшения вашего кода.\nПоддерживает множество языков программирования и моделей ИИ.`,
}

func main() {
	// Инициализация конфигурации перед выполнением команд
	cmd.InitConfig()

	// Добавление команд в корневую команду
	rootCmd.AddCommand(cmd.FmtCmd, cmd.SetCmd)

	// Выполнение корневой команды
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка выполнения команды:", err)
		os.Exit(1)
	}
}