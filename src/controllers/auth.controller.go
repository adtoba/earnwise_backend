package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/adtoba/earnwise_backend/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthController struct {
	DB          *gorm.DB
	TokenMaker  *utils.JWTMaker
	RedisClient *redis.Client
}

func NewAuthController(db *gorm.DB, tokenMaker *utils.JWTMaker, redisClient *redis.Client) *AuthController {
	return &AuthController{DB: db, TokenMaker: tokenMaker, RedisClient: redisClient}
}

func (ac *AuthController) Login(c *gin.Context) {
	var payload models.LoginUserRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorResponse("User not found", nil))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	if err := utils.CompareHashAndPassword(payload.Password, user.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("Invalid email or password", nil))
		return
	}

	accessToken, _, err := ac.TokenMaker.CreateToken(user.ID, user.Email, user.Role, time.Hour*24, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}

	refreshToken, _, err := ac.TokenMaker.CreateToken(user.ID, user.Email, user.Role, time.Hour*24*30, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}

	var expertProfile models.ExpertProfile
	res := ac.DB.First(&expertProfile, "user_id = ?", user.ID)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	var expertProfileSummary *models.ExpertProfileSummaryResponse
	if res.Error == nil {
		summary := expertProfile.ToExpertProfileSummaryResponse()
		expertProfileSummary = &summary
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Login successful", models.LoginUserResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		User:          user.ToUserResponse(),
		ExpertProfile: expertProfileSummary,
	}))
}

func (ac *AuthController) Register(c *gin.Context) {
	var payload models.CreateUserRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}

	user := models.User{
		FirstName:       payload.FirstName,
		LastName:        payload.LastName,
		Email:           strings.ToLower(payload.Email),
		Password:        hashedPassword,
		Role:            "user",
		IsEmailVerified: false,
		IsPhoneVerified: false,
		IsBlocked:       false,
	}

	result := ac.DB.Create(&user)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User registered successfully", user.ToUserResponse()))
}
