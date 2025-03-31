package service

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// setup загружает переменные окружения из файла .env
func setup() error {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("Не удалось загрузить .env файл: ", err)
		return err
	}
	return nil
}

// tearDown выполняет очистку после завершения тестов
func tearDown() error {
	// В данном случае очистка не требуется, но функция оставлена для будущих расширений
	return nil
}

// TestMain является точкой входа для запуска тестов
func TestMain(m *testing.M) {
	// Инициализация перед запуском тестов
	if err := setup(); err != nil {
		log.Fatal("Ошибка при инициализации: ", err)
		os.Exit(1)
	}

	// Запуск тестов
	exitCode := m.Run()

	// Очистка после завершения тестов
	if err := tearDown(); err != nil {
		log.Fatal("Ошибка при завершении: ", err)
		os.Exit(1)
	}

	// Возврат кода завершения
	os.Exit(exitCode)
}
