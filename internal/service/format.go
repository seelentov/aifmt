package service

import (
	"fmt"

	"github.com/seelentov/aifmt/internal/entity"
	"github.com/seelentov/aifmt/pkg/api"
)

type AIFormatCodeRequest struct {
	Code    string           `json:"code"`
	Updates []*entity.Update `json:"updates"`
}

func FormatCode(content, language, model, token string, ctx []*entity.File) (string, []*entity.Update, error) {
	format := "Исправь этот код: ```%s\n%s\n```. Устрани ошибки, прокомментируй его, проведи оптимизацию. Комментарии должны быть на русском языке. Ответ должен быть в формате json: {code:(новый код), updates:(массив изменений)[{code:(часть кода, которую ты решил изменить), description:(причина изменения)}]}."
	p := fmt.Sprintf(string(format), language, content)

	var res *AIFormatCodeRequest

	dialog := make([]*entity.Message, 0)
	dialog = append(dialog, &entity.Message{Text: p, IsUser: true})

	if len(ctx) > 1 {
		ctxPr := "Так же учти и другие файлы этого же проекта. Я пришлю тебе список в виде отдельных сообщений: "
		dialog = append(dialog, &entity.Message{Text: ctxPr, IsUser: true})

		for _, file := range ctx {
			filePr := fmt.Sprintf("%s:\n```%s\n%s\n```", file.Path, language, file.Content)
			dialog = append(dialog, &entity.Message{Text: filePr, IsUser: true})
		}
	}

	if err := api.GetAnswer(token, model, dialog, &res); err != nil {
		return "", nil, nil
	}

	return res.Code, res.Updates, nil
}
