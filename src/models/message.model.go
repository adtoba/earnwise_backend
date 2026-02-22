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
	ContentType  string      `json:"content_type"`
	IsResponseTo string      `json:"is_response_to"`
	ResponseToID string      `json:"response_to_id"`
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
	ContentType  string      `json:"content_type" binding:"required"`
	Attachments  StringArray `json:"attachments"`
	IsResponseTo string      `json:"is_response_to"`
	ResponseToID string      `json:"response_to_id"`
	ResponseType string      `json:"response_type" binding:"required"`
}

type EditMessageRequest struct {
	Content string `json:"content" binding:"required"`
}
