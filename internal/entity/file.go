package entity

// File представляет структуру файла, содержащего путь и его содержимое.
type File struct {
	Path    string `json:"path"`    // Путь к файлу
	Content string `json:"content"` // Содержимое файла
}
