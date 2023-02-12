package models

import (
	"github.com/oxanahr/discord-bot/cmd/database"
	"time"
)

type Comment struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	TaskID    uint64    `json:"taskID" gorm:"not null"`
	AuthorID  string    `json:"authorID" gorm:"not null"`
	Text      string    `json:"text" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
}

func (t *Comment) Create() error {
	t.CreatedAt = time.Now()
	return database.DB.Create(t).Error
}
