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
	ID                 string           `json:"id" gorm:"primaryKey"`
	UserID             string           `json:"user_id"`
	User               User             `json:"user" gorm:"foreignKey:UserID;references:ID"`
	ProfessionalTitle  string           `json:"professional_title"`
	Categories         StringArray      `json:"categories"`
	Bio                string           `json:"bio"`
	Faq                StringArray      `json:"faq"`
	Rates              Rates            `json:"rates" gorm:"type:jsonb"`
	Availability       AvailabilityList `json:"availability" gorm:"type:jsonb"`
	Socials            Socials          `json:"socials" gorm:"type:jsonb"`
	VerificationStatus string           `json:"verification_status" gorm:"default:pending"`
	Rating             float64          `json:"rating"`
	ReviewsCount       int              `json:"reviews_count"`
	TotalConsultations int              `json:"total_consultations"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

func (e *ExpertProfile) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.NewString()
	return
}

func (e *ExpertProfile) ToExpertProfileResponse() ExpertProfileResponse {
	return ExpertProfileResponse{
		ID:                 e.ID,
		User:               e.User.ToUserResponse(),
		ProfessionalTitle:  e.ProfessionalTitle,
		Categories:         e.Categories,
		Bio:                e.Bio,
		Faq:                e.Faq,
		Rates:              e.Rates,
		Availability:       e.Availability,
		VerificationStatus: e.VerificationStatus,
		ReviewsCount:       e.ReviewsCount,
		Rating:             e.Rating,
		Socials:            e.Socials,
		TotalConsultations: e.TotalConsultations,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
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
	Day    string `json:"day"`
	Status string `json:"status" gorm:"default:available"`
	Start  string `json:"start"`
	End    string `json:"end"`
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

type Socials struct {
	Instagram string `json:"instagram"`
	X         string `json:"x"`
	Linkedin  string `json:"linkedin"`
	Website   string `json:"website"`
}

func (s Socials) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (s *Socials) Scan(value interface{}) error {
	if value == nil {
		*s = Socials{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("unsupported type for Socials: %T", value)
	}
}

type SavedExpert struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	ExpertID  string    `json:"expert_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *SavedExpert) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewString()
	return
}

type SaveExpertRequest struct {
	ExpertID string `json:"expert_id" binding:"required"`
}

type CreateExpertProfileRequest struct {
	ProfessionalTitle string           `json:"professional_title" binding:"required"`
	Categories        StringArray      `json:"categories" binding:"required"`
	Bio               string           `json:"bio" binding:"required"`
	Faq               StringArray      `json:"faq" binding:"required"`
	Rates             Rates            `json:"rates" binding:"required"`
	Availability      AvailabilityList `json:"availability" binding:"required"`
	Socials           Socials          `json:"socials" gorm:"type:jsonb"`
}

type UpdateExpertDetailsRequest struct {
	ProfessionalTitle string      `json:"professional_title"`
	Categories        StringArray `json:"categories"`
	Bio               string      `json:"bio"`
	Faq               StringArray `json:"faq"`
}

type UpdateExpertRateRequest struct {
	Text  float64 `json:"text"`
	Video float64 `json:"video"`
	Call  float64 `json:"call"`
}

type UpdateExpertSocialsRequest struct {
	Instagram string `json:"instagram"`
	X         string `json:"x"`
	Linkedin  string `json:"linkedin"`
	Website   string `json:"website"`
}

type UpdateExpertAvailabilityRequest struct {
	Availability AvailabilityList `json:"availability"`
}

type ExpertProfileResponse struct {
	ID                 string           `json:"id"`
	User               UserResponse     `json:"user"`
	ProfessionalTitle  string           `json:"professional_title"`
	Categories         StringArray      `json:"categories"`
	Bio                string           `json:"bio"`
	Faq                StringArray      `json:"faq"`
	Rates              Rates            `json:"rates"`
	Availability       AvailabilityList `json:"availability"`
	Socials            Socials          `json:"socials"`
	VerificationStatus string           `json:"verification_status"`
	Rating             float64          `json:"rating"`
	ReviewsCount       int              `json:"reviews_count"`
	TotalConsultations int              `json:"total_consultations"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}
