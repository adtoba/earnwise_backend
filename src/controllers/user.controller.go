package controllers

import (
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) GetUserById(c *gin.Context) {
	var user models.User

	result := uc.DB.First(&user, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User fetched successfully", user))
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var payload models.UpdateUserRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var user models.User
	result := uc.DB.First(&user, "id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	user.Gender = payload.Gender
	user.Phone = payload.PhoneNumber
	user.Country = payload.Country
	user.State = payload.State
	user.City = payload.City
	user.Address = payload.Address
	user.Zip = payload.Zip
	user.ProfilePicture = payload.ProfilePicture
	user.DOB = payload.DOB

	res := uc.DB.Save(&user)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User updated successfully", user.ToUserResponse()))
}

func (uc *UserController) SaveExpert(c *gin.Context) {
	var payload models.SaveExpertRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var savedExpert models.SavedExpert
	result := uc.DB.First(&savedExpert, "user_id = ? AND expert_id = ?", c.MustGet("user_id").(string), payload.ExpertID)

	if savedExpert.ID != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Expert already saved", nil))
		return
	}

	savedExpert.UserID = c.MustGet("user_id").(string)
	savedExpert.ExpertID = payload.ExpertID

	res := uc.DB.Create(&savedExpert)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Expert saved successfully", savedExpert))
}

func (uc *UserController) GetSavedExperts(c *gin.Context) {
	var savedExperts []models.SavedExpert
	result := uc.DB.Where("user_id = ?", c.MustGet("user_id").(string)).Find(&savedExperts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var experts []models.ExpertProfile
	for _, savedExpert := range savedExperts {
		var expert models.ExpertProfile
		result := uc.DB.Preload("User").First(&expert, "id = ?", savedExpert.ExpertID)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
			return
		}
		experts = append(experts, expert)
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Saved experts fetched successfully", experts))
}

func (uc *UserController) GetUserProfile(c *gin.Context) {
	var user models.User
	result := uc.DB.First(&user, "id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var expertProfile models.ExpertProfile
	res := uc.DB.First(&expertProfile, "user_id = ?", user.ID)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var expertProfileSummary *models.ExpertProfileSummaryResponse
	if res.Error == nil {
		summary := expertProfile.ToExpertProfileSummaryResponse()
		expertProfileSummary = &summary
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User profile fetched successfully", models.UserProfileResponse{
		User:          user.ToUserResponse(),
		ExpertProfile: expertProfileSummary,
	}))
}
