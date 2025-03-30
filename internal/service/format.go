package service

import (
	"aifmt/internal/entity"
	"aifmt/pkg/api"
	"fmt"
)

type AIFormatCodeRequest struct {
	Code    string           `json:"code"`
	Updates []*entity.Update `json:"updates"`
}

func FormatCode(content, language, model, token string) (string, []*entity.Update, error) {
	format := "Исправь этот код: ```%s%s```. Устрани ошибки, прокомментируй его, проведи оптимизацию. Комментарии должны быть на русском языке. Ответ должен быть в формате json: {code:(новый код), updates:(массив изменений)[{code:(часть кода, которую ты решил изменить), description:(причина изменения)}]}."

	p := fmt.Sprintf(string(format), language, content)

	var res *AIFormatCodeRequest

	if err := api.GetAnswer(token, model, []*entity.Message{{Text: p, IsUser: true}}, &res); err != nil {
		return "", nil, nil
	}

	return res.Code, res.Updates, nil
}
