package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/seelentov/aifmt/internal/entity"
)

type response struct {
	Choices []*choice `json:"choices"`
}

type choice struct {
	Message *message `json:"message"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GetAnswer отправляет запрос к API OpenRouter и возвращает ответ
func GetAnswer(token string, model string, dialog []*entity.Message, target interface{}) error {
	rb := struct {
		Model       string     `json:"model"`
		Messages    []*message `json:"messages"`
		Temperature float64    `json:"temperature"`
	}{
		Model:       model,
		Temperature: 0.3,
	}

	// Преобразуем диалог в формат, понятный API
	for _, item := range dialog {
		role := "assistant"
		if item.IsUser {
			role = "user"
		}

		rb.Messages = append(rb.Messages, &message{
			Role:    role,
			Content: item.Text,
		})
	}

	bodyBytes, err := json.Marshal(rb)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга тела запроса: %w", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	resBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка в ответе: %v %s", resp.StatusCode, resBodyBytes)
	}

	tempTarget := &response{}
	err = json.Unmarshal(resBodyBytes, tempTarget)
	if err != nil {
		return fmt.Errorf("ошибка анмаршалинга ответа: %w", err)
	}

	msg := tempTarget.Choices[len(tempTarget.Choices)-1].Message.Content

	// Обработка ответа в зависимости от типа целевого объекта
	if reflect.TypeOf(target).String() == "*string" {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(msg))
		return nil
	}

	msg = strings.TrimPrefix(msg, "```json\n")
	msg = strings.TrimSuffix(msg, "\n```")

	err = json.Unmarshal([]byte(msg), &target)
	if err != nil {
		return fmt.Errorf("ошибка анмаршалинга: %w: %s", err, msg[0:20]+"...")
	}

	return nil
}