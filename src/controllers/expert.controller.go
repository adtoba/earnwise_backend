package controllers

import (
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExpertController struct {
	DB               *gorm.DB
	WalletController *WalletController
}

func NewExpertController(db *gorm.DB, walletController *WalletController) *ExpertController {
	return &ExpertController{DB: db, WalletController: walletController}
}

func (ec *ExpertController) CreateExpertProfile(c *gin.Context) {
	var payload models.CreateExpertProfileRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var profileExists models.ExpertProfile
	result := ec.DB.First(&profileExists, "user_id = ?", c.MustGet("user_id").(string))

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
		Socials:           payload.Socials,
		Availability:      payload.Availability,
	}

	res := ec.DB.Create(&expertProfile)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	ec.WalletController.CreateWallet(c, expertProfile.ID)

	c.JSON(http.StatusOK, models.SuccessResponse("Expert profile created successfully", expertProfile))
}

func (ec *ExpertController) UpdateExpertRate(c *gin.Context) {
	var payload models.UpdateExpertRateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var expertProfile models.ExpertProfile
	result := ec.DB.First(&expertProfile, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	expertProfile.Rates = models.Rates{
		Text:  payload.Text,
		Video: payload.Video,
		Call:  payload.Call,
	}

	res := ec.DB.Save(&expertProfile)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", res.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert rate updated successfully", expertProfile.ToExpertProfileResponse()))
}

func (ec *ExpertController) UpdateExpertSocials(c *gin.Context) {
	var payload models.UpdateExpertSocialsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var expertProfile models.ExpertProfile
	result := ec.DB.First(&expertProfile, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	expertProfile.Socials = models.Socials{
		Instagram: payload.Instagram,
		X:         payload.X,
		Linkedin:  payload.Linkedin,
		Website:   payload.Website,
	}

	res := ec.DB.Save(&expertProfile)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", res.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert socials updated successfully", expertProfile.ToExpertProfileResponse()))
}

func (ec *ExpertController) UpdateExpertDetails(c *gin.Context) {
	var payload models.UpdateExpertDetailsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var expertProfile models.ExpertProfile
	result := ec.DB.First(&expertProfile, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	expertProfile.ProfessionalTitle = payload.ProfessionalTitle
	expertProfile.Categories = payload.Categories
	expertProfile.Bio = payload.Bio
	expertProfile.Faq = payload.Faq

	res := ec.DB.Save(&expertProfile)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", res.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert details updated successfully", expertProfile.ToExpertProfileResponse()))
}

func (ec *ExpertController) UpdateExpertAvailability(c *gin.Context) {
	var payload models.UpdateExpertAvailabilityRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var expertProfile models.ExpertProfile
	result := ec.DB.First(&expertProfile, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	expertProfile.Availability = payload.Availability
	res := ec.DB.Save(&expertProfile)

	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", res.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert availability updated successfully", expertProfile.ToExpertProfileResponse()))
}

func (ec *ExpertController) GetExpertProfileById(c *gin.Context) {
	var expertProfile models.ExpertProfile
	result := ec.DB.Preload("User").Where("id = ?", c.Param("id")).First(&expertProfile)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert profile fetched successfully", expertProfile.ToExpertProfileResponse()))
}

func (ec *ExpertController) GetExpertProfile(c *gin.Context) {
	var expertProfile models.ExpertProfile
	result := ec.DB.Preload("User").Where("user_id = ?", c.MustGet("user_id").(string)).First(&expertProfile)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Expert profile fetched successfully", expertProfile.ToExpertProfileResponse()))
}

func (ec *ExpertController) GetExpertsByCategory(c *gin.Context) {
	var experts []models.ExpertProfile
	result := ec.DB.Preload("User").Where("categories @> ?", c.Param("category")).Find(&experts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Experts fetched successfully", experts))
}
