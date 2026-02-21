package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID           string      `json:"id" gorm:"primaryKey"`
	ChatID       string      `json:"chat_id"`
	SenderID     string      `json:"sender_id"`
	ReceiverID   string      `json:"receiver_id"`
	Content      string      `json:"content"`
	ResponseType string      `json:"response_type"`
	Attachments  StringArray `json:"attachments"`
	IsRead       bool        `json:"is_read"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.NewString()
	return
}

type CreateMessageRequest struct {
	SenderID     string      `json:"sender_id" binding:"required"`
	ReceiverID   string      `json:"receiver_id" binding:"required"`
	Content      string      `json:"content" binding:"required"`
	Attachments  StringArray `json:"attachments"`
	ResponseType string      `json:"response_type" binding:"required"`
}
