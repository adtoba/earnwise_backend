package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
	ID        string        `json:"id" gorm:"primaryKey"`
	UserID    string        `json:"user_id"`
	ExpertID  string        `json:"expert_id"`
	User      User          `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Expert    ExpertProfile `json:"expert" gorm:"foreignKey:ExpertID;references:ID"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (c *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return
}

func (c *Chat) ToChatResponse(message Message) ChatResponse {
	return ChatResponse{
		ID:          c.ID,
		UserID:      c.UserID,
		ExpertID:    c.ExpertID,
		User:        c.User,
		Expert:      c.Expert,
		LastMessage: message,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type ChatResponse struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	ExpertID    string        `json:"expert_id"`
	User        User          `json:"user"`
	LastMessage Message       `json:"last_message"`
	Expert      ExpertProfile `json:"expert"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type CreateChatRequest struct {
	ExpertID     string `json:"expert_id" binding:"required"`
	Message      string `json:"message" binding:"required"`
	ResponseType string `json:"response_type" binding:"required"`
}
