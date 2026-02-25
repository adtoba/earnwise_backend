package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	CallStatusPending   = "pending"
	CallStatusAccepted  = "accepted"
	CallStatusRejected  = "rejected"
	CallStatusCancelled = "cancelled"
	CallStatusExpired   = "expired"
	CallStatusCompleted = "completed"

	PaymentStatusHeld     = "held"
	PaymentStatusCaptured = "captured"
	PaymentStatusReleased = "released"
)

type Call struct {
	ID            string        `json:"id" gorm:"primaryKey"`
	UserID        string        `json:"user_id"`
	ExpertID      string        `json:"expert_id"`
	User          User          `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Expert        ExpertProfile `json:"expert" gorm:"foreignKey:ExpertID;references:ID"`
	ScheduledAt   time.Time     `json:"scheduled_at"`
	Subject       string        `json:"subject"`
	Description   string        `json:"description"`
	DurationMins  int           `json:"duration_mins"`
	Status        string        `json:"status" gorm:"default:pending"`
	Price         float64       `json:"price"`
	PaymentStatus string        `json:"payment_status"`
	PaymentRef    string        `json:"payment_ref"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func (c *Call) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return
}

type CreateCallRequest struct {
	ExpertID     string    `json:"expert_id" binding:"required"`
	ScheduledAt  time.Time `json:"scheduled_at" binding:"required"`
	Subject      string    `json:"subject" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	DurationMins int       `json:"duration_mins" binding:"required"`
	Price        float64   `json:"price"`
	PaymentRef   string    `json:"payment_ref"`
	Timezone     string    `json:"timezone"`
}

type CallResponse struct {
	ID            string                `json:"id"`
	UserID        string                `json:"user_id"`
	ExpertID      string                `json:"expert_id"`
	Subject       string                `json:"subject"`
	Description   string                `json:"description"`
	Timezone      string                `json:"timezone"`
	User          UserResponse          `json:"user"`
	Expert        ExpertProfileResponse `json:"expert"`
	ScheduledAt   time.Time             `json:"scheduled_at"`
	DurationMins  int                   `json:"duration_mins"`
	Status        string                `json:"status"`
	Price         float64               `json:"price"`
	PaymentStatus string                `json:"payment_status"`
	PaymentRef    string                `json:"payment_ref"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

func (c *Call) ToCallResponse() CallResponse {
	return CallResponse{
		ID:            c.ID,
		UserID:        c.UserID,
		ExpertID:      c.ExpertID,
		Subject:       c.Subject,
		Description:   c.Description,
		User:          c.User.ToUserResponse(),
		Expert:        c.Expert.ToExpertProfileResponse(),
		ScheduledAt:   c.ScheduledAt,
		DurationMins:  c.DurationMins,
		Status:        c.Status,
		Price:         c.Price,
		PaymentStatus: c.PaymentStatus,
		PaymentRef:    c.PaymentRef,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}
