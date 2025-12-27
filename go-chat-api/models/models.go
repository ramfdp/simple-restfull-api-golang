package models

import "time"

type Message struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
