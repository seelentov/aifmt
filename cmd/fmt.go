package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/seelentov/aifmt/internal/entity"
	"github.com/seelentov/aifmt/internal/service"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var wg sync.WaitGroup

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
		comments, _ := cmd.Flags().GetBool("comments")
		report, _ := cmd.Flags().GetBool("report")
		skip, _ := cmd.Flags().GetBool("skip")
		maxRetries := viper.GetInt("max_retry")

		commentsLanguage := viper.GetString("comments_language")
		if report && commentsLanguage == "" {
			fmt.Println("Язык комментариев не настроен. Пожалуйста, сначала выполните 'aifmt set comments_language язык'.")
			os.Exit(1)
		}

		if len(args) == 0 {
			fmt.Println("Ошибка: не указаны файлы для обработки")
			cmd.Help()
			os.Exit(1)
		}

		// Собираем контекстные файлы, если указан флаг
		var ctx []*entity.File
		var allUpds []*entity.Update

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

		repname := time.Now().Format("report_2006-01-02_15:04:05.json")

		// Обрабатываем каждый файл
		for _, pattern := range args {
			files, err := filepath.Glob(pattern)
			if err != nil {
				fmt.Printf("Ошибка при разборе шаблона %s: %v\n", pattern, err)
				continue
			}

			for _, file := range files {
				wg.Add(1)
				go func(file string) {
					fmt.Printf("Обработка %s (Язык: %s, Модель: %s, Контекст: %v)...\n",
						file, language, model, withCtx)

					content, err := os.ReadFile(file)
					if err != nil {
						fmt.Printf("Ошибка чтения файла %s: %v\n", file, err)
						if !skip {
							fmt.Println("Попытка повторного чтения файла...")
							if err := retryOperation(maxRetries, func() error {
								content, err = os.ReadFile(file)
								return err
							}); err != nil {
								fmt.Printf("Не удалось прочитать файл %s после %d попыток: %v\n", file, maxRetries, err)
								return
							}
						} else {
							return
						}
					}

					var u string
					var upds []*entity.Update
					var formatErr error

					// Функция для форматирования кода
					formatFunc := func() error {
						u, upds, formatErr = service.FormatCode(string(content), language, model, token, comments, commentsLanguage, ctx)
						return formatErr
					}

					if err := formatFunc(); err != nil {
						fmt.Printf("Ошибка при форматировании %s: %v\n", file, err)
						if !skip {
							fmt.Println("Попытка повторного форматирования...")
							if err := retryOperation(maxRetries, formatFunc); err != nil {
								fmt.Printf("Не удалось отформатировать файл %s после %d попыток: %v\n", file, maxRetries, err)
								return
							}
						} else {
							return
						}
					}

					if u == "" {
						fmt.Printf("Ошибка: ответ ИИ пуст.\n")
						if !skip {
							fmt.Println("Попытка повторного форматирования из-за пустого ответа...")
							retryCount := 0
							for u == "" && retryCount < maxRetries {
								retryCount++
								if err := formatFunc(); err != nil {
									fmt.Printf("Попытка %d: ошибка форматирования: %v\n", retryCount, err)
									continue
								}
								if u != "" {
									break
								}
								fmt.Printf("Попытка %d: ответ ИИ все еще пуст\n", retryCount)
								time.Sleep(time.Second * time.Duration(retryCount)) // Увеличиваем задержку между попытками
							}
							if u == "" {
								fmt.Printf("Не удалось получить непустой ответ для файла %s после %d попыток\n", file, maxRetries)
								return
							}
						}
					}

					for i := range upds {
						upds[i].Path = file
					}

					// Выводим предложенные изменения
					for _, upd := range upds {
						fmt.Printf("%s:\n```%s\n%s\n```\n%s\n\n", file, language, upd.Code, upd.Description)
					}

					if report {
						allUpds = append(allUpds, upds...)
					}

					// Записываем изменения в файл
					if err := os.WriteFile(file, []byte(u), 0644); err != nil {
						fmt.Printf("Ошибка записи в %s: %v\n", file, err)
						if !skip {
							fmt.Println("Попытка повторной записи файла...")
							if err := retryOperation(maxRetries, func() error {
								return os.WriteFile(file, []byte(u), 0644)
							}); err != nil {
								fmt.Printf("Не удалось записать файл %s после %d попыток: %v\n", file, maxRetries, err)
								return
							}
						}
					}

					fmt.Printf("Файл %s успешно обновлен\n", file)

					if report {
						writetoReport(allUpds, repname)
					}

					wg.Done()
				}(file)
			}
		}

		wg.Wait()
	},
}

var wtrMutex sync.Mutex

func writetoReport(upds []*entity.Update, repname string) {
	wtrMutex.Lock()
	rep, err := json.Marshal(upds)
	if err != nil {
		fmt.Printf("Ошибка при приведении изменений в строку JSON: %s\n", err)
		return
	}

	if err := os.WriteFile(repname, rep, 0644); err != nil {
		fmt.Printf("Ошибка записи в %s: %v\n", repname, err)
		return
	}

	fmt.Printf("Отчет о форматировании записан в %s\n", repname)
	wtrMutex.Unlock()
}

// retryOperation выполняет операцию с повторными попытками при ошибках
func retryOperation(maxRetries int, op func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = op(); err == nil {
			return nil
		}
		fmt.Printf("Попытка %d из %d: %v\n", i+1, maxRetries, err)
		time.Sleep(time.Second * time.Duration(i+1)) // Увеличиваем задержку между попытками
	}
	return fmt.Errorf("достигнуто максимальное количество попыток (%d): %v", maxRetries, err)
}

func init() {
	FmtCmd.Flags().StringP("language", "l", "", "Язык программирования файлов")
	FmtCmd.Flags().StringP("model", "m", "deepseek/deepseek-chat:free", "Модель ИИ для форматирования")
	FmtCmd.Flags().BoolP("with-context", "w", false, "Использовать контекст других файлов при форматировании")
	FmtCmd.Flags().BoolP("comments", "c", false, "Добавить в код комментарии. Язык комментариев настраивается в конфигурации")
	FmtCmd.Flags().BoolP("report", "r", false, "Запись результатов форматирования в файл")
	FmtCmd.Flags().BoolP("skip", "s", false, "Не повторять попытки при ошибках обработки файлов")
}