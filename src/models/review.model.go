package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	ExpertID  string    `json:"expert_id"`
	FullName  string    `json:"full_name"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	return
}

type CreateReviewRequest struct {
	UserID   string  `json:"user_id" binding:"required"`
	ExpertID string  `json:"expert_id" binding:"required"`
	FullName string  `json:"full_name" binding:"required"`
	Rating   float64 `json:"rating" binding:"required"`
	Comment  string  `json:"comment"`
}
