package controllers

import (
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExpertController struct {
	DB *gorm.DB
}

func NewExpertController(db *gorm.DB) *ExpertController {
	return &ExpertController{DB: db}
}

func (ec *ExpertController) CreateExpertProfile(c *gin.Context) {
	var payload models.CreateExpertProfileRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var profileExists models.ExpertProfile
	result := ec.DB.First(&profileExists, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	if profileExists.ID != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Expert profile already exists", nil))
		return
	}

	expertProfile := models.ExpertProfile{
		UserID:            c.MustGet("user_id").(string),
		ProfessionalTitle: payload.ProfessionalTitle,
		Categories:        payload.Categories,
		Bio:               payload.Bio,
		Faq:               payload.Faq,
		Rates:             payload.Rates,
		Availability:      payload.Availability,
	}

	res := ec.DB.Create(&expertProfile)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert profile created successfully", expertProfile))
}

func (ec *ExpertController) GetExpertProfileById(c *gin.Context) {
	var expertProfile models.ExpertProfile
	result := ec.DB.First(&expertProfile, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var user models.User
	result = ec.DB.First(&user, "id = ?", expertProfile.UserID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert profile fetched successfully", expertProfile.ToExpertProfileResponse(user.ToUserResponse())))
}

func (ec *ExpertController) GetExpertProfile(c *gin.Context) {
	var expertProfile models.ExpertProfile
	result := ec.DB.First(&expertProfile, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var user models.User
	result = ec.DB.First(&user, "id = ?", expertProfile.UserID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert profile fetched successfully", expertProfile.ToExpertProfileResponse(user.ToUserResponse())))
}
