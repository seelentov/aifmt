package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SetCmd - команда для установки значений в конфигурации
var SetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Установка значения в конфигурации",
	Long: `Установка или обновление значения в конфигурационном файле.
Пример: aifmt set api_token ваш_токен`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		// Установка значения в конфигурации
		viper.Set(key, value)

		// Сохранение конфигурации
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Ошибка сохранения конфигурации: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Значение '%s' успешно установлено для ключа '%s'\n", value, key)
	},
}

// InitConfig - инициализация конфигурации
func InitConfig() {
	// Получение домашней директории пользователя
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Ошибка получения домашней директории: %v\n", err)
		os.Exit(1)
	}

	// Определение пути к конфигурационной директории и файлу
	configDir := filepath.Join(home, ".aifmt")
	configPath := filepath.Join(configDir, "config.yaml")

	// Создание конфигурационной директории, если она не существует
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			fmt.Printf("Ошибка создания конфигурационной директории: %v\n", err)
			os.Exit(1)
		}
	}

	// Настройка Viper для работы с конфигурационным файлом
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Чтение конфигурационного файла или создание нового, если он отсутствует
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Конфигурационный файл не найден, создается новый в", configPath)

			// Установка значений по умолчанию
			viper.Set("comments_language", "Русский")
			viper.Set("api_key", "")
			viper.Set("max_retry", 5)
			viper.Set("channels", 10)

			// Запись конфигурации в файл
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				fmt.Printf("Ошибка записи конфигурации: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Ошибка чтения конфигурации: %v\n", err)
			os.Exit(1)
		}
	}
}

// init - инициализация команды и конфигурации
func init() {
	InitConfig()
}
