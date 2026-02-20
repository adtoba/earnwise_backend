package controllers

import (
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

func (ec *ExpertController) GetExpertDashboard(c *gin.Context) {
	var expertProfile models.ExpertProfile
	result := ec.DB.Preload("User").First(&expertProfile, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var wallet models.Wallet
	result = ec.DB.First(&wallet, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	expertDashboardResponse := models.ExpertDashboardResponse{
		ExpertProfile: expertProfile.ToExpertProfileResponse(),
		Wallet:        wallet,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Expert dashboard fetched successfully", expertDashboardResponse))
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
	currentUserID := c.MustGet("user_id").(string)
	savedMap, err := ec.getSavedExpertMap(currentUserID, []string{expertProfile.ID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}

	expertResponse := expertProfile.ToExpertProfileResponse()
	expertResponse.IsSaved = savedMap[expertProfile.ID]
	c.JSON(http.StatusOK, models.SuccessResponse("Expert profile fetched successfully", expertResponse))
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
	category := c.Param("category")
	currentUserID := c.MustGet("user_id").(string)
	result := ec.DB.Preload("User").
		Where("categories @> ?", pq.Array([]string{category})).
		Where("user_id <> ?", currentUserID).
		Find(&experts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var expertResponses []models.ExpertProfileResponse
	expertIDs := make([]string, 0, len(experts))
	for _, expert := range experts {
		expertIDs = append(expertIDs, expert.ID)
	}
	savedMap, err := ec.getSavedExpertMap(currentUserID, expertIDs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}
	for _, expert := range experts {
		expertResponse := expert.ToExpertProfileResponse()
		expertResponse.IsSaved = savedMap[expert.ID]
		expertResponses = append(expertResponses, expertResponse)
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Experts fetched successfully", expertResponses))
}

func (ec *ExpertController) GetRecommendedTopExperts(c *gin.Context) {
	var experts []models.ExpertProfile
	currentUserID := c.MustGet("user_id").(string)
	result := ec.DB.Preload("User").Order("rating DESC").Limit(10).Where("user_id <> ?", currentUserID).Find(&experts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var expertResponses []models.ExpertProfileResponse
	expertIDs := make([]string, 0, len(experts))
	for _, expert := range experts {
		expertIDs = append(expertIDs, expert.ID)
	}
	savedMap, err := ec.getSavedExpertMap(currentUserID, expertIDs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}
	for _, expert := range experts {
		expertResponse := expert.ToExpertProfileResponse()
		expertResponse.IsSaved = savedMap[expert.ID]
		expertResponses = append(expertResponses, expertResponse)
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Recommended top experts fetched successfully", expertResponses))
}

func (ec *ExpertController) getSavedExpertMap(userID string, expertIDs []string) (map[string]bool, error) {
	savedMap := make(map[string]bool)
	if len(expertIDs) == 0 {
		return savedMap, nil
	}

	var savedExperts []models.SavedExpert
	result := ec.DB.Where("user_id = ? AND expert_id IN ?", userID, expertIDs).Find(&savedExperts)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, savedExpert := range savedExperts {
		savedMap[savedExpert.ExpertID] = true
	}

	return savedMap, nil
}
