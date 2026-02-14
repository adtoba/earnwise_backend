package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpertProfile struct {
	ID                string           `json:"id" gorm:"primaryKey"`
	UserID            string           `json:"user_id"`
	ProfessionalTitle string           `json:"professional_title"`
	Categories        StringArray      `json:"categories"`
	Bio               string           `json:"bio"`
	Faq               StringArray      `json:"faq"`
	Rates             Rates            `json:"rates" gorm:"type:jsonb"`
	Availability      AvailabilityList `json:"availability" gorm:"type:jsonb"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

func (e *ExpertProfile) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.NewString()
	return
}

func (e *ExpertProfile) ToExpertProfileResponse(user UserResponse) ExpertProfileResponse {
	return ExpertProfileResponse{
		ID:                e.ID,
		User:              user,
		ProfessionalTitle: e.ProfessionalTitle,
		Categories:        e.Categories,
		Bio:               e.Bio,
		Faq:               e.Faq,
		Rates:             e.Rates,
		Availability:      e.Availability,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

type Rates struct {
	Text  float64 `json:"text"`
	Video float64 `json:"video"`
	Call  float64 `json:"call"`
}

func (r Rates) Value() (driver.Value, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (r *Rates) Scan(value interface{}) error {
	if value == nil {
		*r = Rates{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, r)
	case string:
		return json.Unmarshal([]byte(v), r)
	default:
		return fmt.Errorf("unsupported type for Rates: %T", value)
	}
}

type Availability struct {
	Day   string    `json:"day"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type AvailabilityList []Availability

func (a AvailabilityList) Value() (driver.Value, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (a *AvailabilityList) Scan(value interface{}) error {
	if value == nil {
		*a = AvailabilityList{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return fmt.Errorf("unsupported type for AvailabilityList: %T", value)
	}
}

type CreateExpertProfileRequest struct {
	ProfessionalTitle string           `json:"professional_title" binding:"required"`
	Categories        StringArray      `json:"categories" binding:"required"`
	Bio               string           `json:"bio" binding:"required"`
	Faq               StringArray      `json:"faq" binding:"required"`
	Rates             Rates            `json:"rates" binding:"required"`
	Availability      AvailabilityList `json:"availability" binding:"required"`
}

type ExpertProfileResponse struct {
	ID                string           `json:"id"`
	User              UserResponse     `json:"user"`
	ProfessionalTitle string           `json:"professional_title"`
	Categories        StringArray      `json:"categories"`
	Bio               string           `json:"bio"`
	Faq               StringArray      `json:"faq"`
	Rates             Rates            `json:"rates"`
	Availability      AvailabilityList `json:"availability"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}
