package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Установка значения в конфигурации",
	Long: `Установка или обновление значения в конфигурационном файле.
Пример: aifmt set api_token ваш_токен`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		viper.Set(key, value)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("Ошибка сохранения конфигурации: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Значение '%s' успешно установлено для ключа '%s'\n", value, key)
	},
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Ошибка получения домашней директории: %v\n", err)
		return
	}

	configDir := filepath.Join(home, ".aifmt")
	configPath := filepath.Join(configDir, "config.yaml")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			fmt.Printf("Ошибка создания конфигурационной директории: %v\n", err)
			return
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Конфигурационный файл не найден, создается новый в", configPath)

			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				fmt.Printf("Ошибка записи конфигурации: %v\n", err)
			}
		} else {
			fmt.Printf("Ошибка чтения конфигурации: %v\n", err)
		}
	}
}

func init() {
	initConfig()
}
