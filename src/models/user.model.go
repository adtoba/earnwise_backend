package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email" gorm:"unique"`
	Password        string    `json:"password"`
	Gender          string    `json:"gender"`
	DOB             time.Time `json:"dob"`
	Phone           string    `json:"phone"`
	Address         string    `json:"address"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	Zip             string    `json:"zip"`
	Country         string    `json:"country"`
	Role            string    `json:"role"`
	ProfilePicture  string    `json:"profile_picture"`
	IsBlocked       bool      `json:"is_blocked"`
	IsEmailVerified bool      `json:"is_email_verified"`
	IsPhoneVerified bool      `json:"is_phone_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Gender          string    `json:"gender"`
	DOB             time.Time `json:"dob"`
	Phone           string    `json:"phone"`
	Address         string    `json:"address"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	Zip             string    `json:"zip"`
	Country         string    `json:"country"`
	Role            string    `json:"role"`
	IsBlocked       bool      `json:"is_blocked"`
	IsEmailVerified bool      `json:"is_email_verified"`
	IsPhoneVerified bool      `json:"is_phone_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return
}

func (user *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Gender:          user.Gender,
		DOB:             user.DOB,
		Phone:           user.Phone,
		Address:         user.Address,
		City:            user.City,
		State:           user.State,
		Zip:             user.Zip,
		Country:         user.Country,
		Role:            user.Role,
		IsBlocked:       user.IsBlocked,
		IsEmailVerified: user.IsEmailVerified,
		IsPhoneVerified: user.IsPhoneVerified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Gender         string    `json:"gender" binding:"required"`
	DOB            time.Time `json:"dob"`
	PhoneNumber    string    `json:"phone_number" binding:"required"`
	ProfilePicture string    `json:"profile_picture"`
	Country        string    `json:"country" binding:"required"`
	State          string    `json:"state" binding:"required"`
	City           string    `json:"city" binding:"required"`
	Address        string    `json:"address" binding:"required"`
	Zip            string    `json:"zip" binding:"required"`
}

type LoginUserResponse struct {
	AccessToken   string                        `json:"access_token"`
	RefreshToken  string                        `json:"refresh_token"`
	User          UserResponse                  `json:"user"`
	ExpertProfile *ExpertProfileSummaryResponse `json:"expert_profile"`
}

type UserProfileResponse struct {
	User          UserResponse                  `json:"user"`
	ExpertProfile *ExpertProfileSummaryResponse `json:"expert_profile"`
}
