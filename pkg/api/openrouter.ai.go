package api

import (
	"aifmt/internal/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
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

func GetAnswer(token string, model string, dialog []*entity.Message, target interface{}) error {
	rb := struct {
		Model       string        `json:"model"`
		Messages    []interface{} `json:"messages"`
		Temperature float64       `json:"temperature"`
	}{
		Model:       model,
		Temperature: 0.1,
	}

	for _, item := range dialog {
		role := "assistant"
		if item.IsUser {
			role = "user"
		}

		rb.Messages = append(rb.Messages, map[string]string{
			"role":    role,
			"content": item.Text,
		})
	}

	bodyBytes, err := json.Marshal(rb)
	if err != nil {
		return fmt.Errorf("error marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("error create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error do request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in response: %v %s", resp.StatusCode, bodyBytes)
	}

	tempTarget := &response{}
	err = json.Unmarshal(bodyBytes, tempTarget)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	msg := tempTarget.Choices[len(tempTarget.Choices)-1].Message.Content

	if reflect.TypeOf(target).String() == "*string" {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(msg))
		return nil
	}

	msg = strings.TrimPrefix(msg, "```json\n")
	msg = strings.TrimSuffix(msg, "\n```")

	err = json.Unmarshal([]byte(msg), &target)
	if err != nil {
		return fmt.Errorf("error unmarshal: %w", err)
	}

	return nil
}
