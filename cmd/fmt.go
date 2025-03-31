package cmd

import (
	"aifmt/internal/entity"
	"aifmt/internal/service"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var FmtCmd = &cobra.Command{
	Use:   "fmt [флаги] [файлы...]",
	Short: "Форматирование кода с помощью ИИ",
	Long: `Форматирование одного или нескольких файлов с кодом с использованием ИИ.
Вы можете указать язык программирования и модель ИИ для использования.
Если токен не настроен, вам будет предложено его установить.`,
	Example: `  # Форматирование Go файла
  aifmt fmt -l go main.go
  
  # Форматирование нескольких Python файлов с указанной моделью
  aifmt fmt -l python --model claude-2 *.py
  
  # Форматирование с автоопределением языка
  aifmt fmt script.js
  
  # Форматирование с учетом контекста других файлов
  aifmt fmt -w -l go *.go`,
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("api_key")
		if token == "" {
			fmt.Println("API токен не настроен. Пожалуйста, сначала выполните 'aifmt set api_key ваш_токен'.")
			os.Exit(1)
		}

		language, _ := cmd.Flags().GetString("language")
		if language == "" {
			fmt.Println("Ошибка: не указан язык программирования")
			os.Exit(1)
		}

		model, _ := cmd.Flags().GetString("model")
		withCtx, _ := cmd.Flags().GetBool("with-context")

		if len(args) == 0 {
			fmt.Println("Ошибка: не указаны файлы для обработки")
			cmd.Help()
			os.Exit(1)
		}

		// Собираем контекстные файлы если указан флаг
		var ctx []*entity.File
		if withCtx {
			for _, pattern := range args {
				files, err := filepath.Glob(pattern)
				if err != nil {
					fmt.Printf("Ошибка при разборе шаблона %s: %v\n", pattern, err)
					continue
				}

				for _, file := range files {
					content, err := os.ReadFile(file)
					if err != nil {
						fmt.Printf("Ошибка чтения контекстного файла %s: %v\n", file, err)
						continue
					}
					ctx = append(ctx, &entity.File{
						Content: string(content),
						Path:    file,
					})
				}
			}
			fmt.Printf("Загружено %d файлов для контекста\n", len(ctx))
		}

		// Обрабатываем каждый файл
		for _, pattern := range args {
			files, err := filepath.Glob(pattern)
			if err != nil {
				fmt.Printf("Ошибка при разборе шаблона %s: %v\n", pattern, err)
				continue
			}

			for _, file := range files {
				fmt.Printf("Обработка %s (Язык: %s, Модель: %s, Контекст: %v)...\n",
					file, language, model, withCtx)

				content, err := os.ReadFile(file)
				if err != nil {
					fmt.Printf("Ошибка чтения файла %s: %v\n", file, err)
					continue
				}

				u, upds, err := service.FormatCode(string(content), language, model, token, ctx)
				if err != nil {
					fmt.Printf("Ошибка при форматировании %s: %v\n", file, err)
					continue
				}
				if u == "" {
					fmt.Printf("Ошибка: ответ ИИ пуст\n")
					continue
				}

				// Выводим предложенные изменения
				for _, upd := range upds {
					fmt.Printf("```%s\n%s\n```\n%s\n\n", language, upd.Code, upd.Description)
				}

				// Записываем изменения в файл
				if err := os.WriteFile(file, []byte(u), 0644); err != nil {
					fmt.Printf("Ошибка записи в %s: %v\n", file, err)
					continue
				}

				fmt.Printf("Файл %s успешно обновлен\n", file)
			}
		}
	},
}

func init() {
	FmtCmd.Flags().StringP("language", "l", "", "Язык программирования файлов")
	FmtCmd.Flags().StringP("model", "m", "deepseek/deepseek-chat:free", "Модель ИИ для форматирования")
	FmtCmd.Flags().BoolP("with-context", "w", false, "Использовать контекст других файлов при форматировании")
}
