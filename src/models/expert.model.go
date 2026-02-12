package models

import "time"

type ExpertProfile struct {
	ID                string         `json:"id" gorm:"primaryKey"`
	ProfessionalTitle string         `json:"professional_title"`
	Bio               string         `json:"bio"`
	Faq               StringArray    `json:"faq"`
	Rates             Rates          `json:"rates"`
	Availability      []Availability `json:"availability"`
}

type Rates struct {
	Text  float64 `json:"text"`
	Video float64 `json:"video"`
	Call  float64 `json:"call"`
}

type Availability struct {
	Day   string    `json:"day"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
