package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	UserID           string    `json:"user_id"`
	ExpertID         string    `json:"expert_id"`
	AvailableBalance float64   `json:"available_balance"`
	PendingBalance   float64   `json:"pending_balance"`
	TotalWithdrawals float64   `json:"total_withdrawals"`
	TotalEarnings    float64   `json:"total_earnings"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.NewString()
	return
}
