package entity

// Message представляет структуру сообщения, используемого в диалоге с ИИ.
type Message struct {
	Text   string `json:"text"`   // Текст сообщения
	IsUser bool   `json:"is_user"` // Флаг, указывающий, является ли сообщение пользовательским
}
