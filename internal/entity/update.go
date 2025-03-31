package entity

// Update представляет структуру обновления кода, содержащего исправленный код и описание изменений.
type Update struct {
	Code        string `json:"code"`        // Исправленный код
	Description string `json:"description"` // Описание изменений
	Path        string `json:"path,omitempty"` // Путь к файлу, если применимо
}