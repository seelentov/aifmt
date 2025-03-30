package entity

type Message struct {
	Text   string `gorm:"not null" json:"text"`
	IsUser bool   `gorm:"not null" json:"is_user"`
}
