package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (notification *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	notification.ID = uuid.NewString()
	return
}
