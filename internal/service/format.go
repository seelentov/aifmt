package service

import (
	"fmt"
	"strings"

	"github.com/seelentov/aifmt/internal/entity"
	"github.com/seelentov/aifmt/pkg/api"
)

type AIFormatCodeRequest struct {
	Code    string           `json:"code"`
	Updates []*entity.Update `json:"updates"`
}

func FormatCode(content, language, model, token string, comment bool, commentsLanguage string, ctx []*entity.File) (string, []*entity.Update, error) {
	format := strings.Builder{}
	format.WriteString("Исправь этот код: ```%s\n%s\n```. Устрани ошибки, проведи оптимизацию. В твоем ответе обязательно должен быть только json объект, без текста до или после в следующем формате: {code:(новый код), updates:(массив изменений)[{code:(часть кода, которую ты решил изменить), description:(причина изменения)}]}!.")
	if comment {
		format.WriteString("Так же закоментируй код. Язык должен быть: %s")
	} else {
		format.WriteString("Не добавляй в код новых комментариев, оставь уже имеющиеся")
	}

	p := fmt.Sprintf(string(format.String()), language, content, commentsLanguage)

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
		return "", nil, err
	}

	return res.Code, res.Updates, nil
}
