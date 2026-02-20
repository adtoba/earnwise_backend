package controllers

import (
	"errors"
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReviewController struct {
	DB *gorm.DB
}

func NewReviewController(db *gorm.DB) *ReviewController {
	return &ReviewController{DB: db}
}

func (rc *ReviewController) CreateReview(c *gin.Context) {
	var payload models.CreateReviewRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var review models.Review

	result := rc.DB.Transaction(func(tx *gorm.DB) error {
		var expert models.ExpertProfile
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&expert, "id = ?", payload.ExpertID).Error; err != nil {
			return err
		}

		review = models.Review{
			UserID:   payload.UserID,
			ExpertID: payload.ExpertID,
			FullName: payload.FullName,
			Rating:   payload.Rating,
			Comment:  payload.Comment,
		}

		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		newCount := expert.ReviewsCount + 1
		newRating := (expert.Rating*float64(expert.ReviewsCount) + payload.Rating) / float64(newCount)

		if err := tx.Model(&expert).Updates(map[string]interface{}{
			"rating":        newRating,
			"reviews_count": newCount,
		}).Error; err != nil {
			return err
		}

		return nil
	})
	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorResponse("Expert not found", nil))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Review created successfully", review))
}

func (rc *ReviewController) GetReviews(c *gin.Context) {
	var reviews []models.Review
	result := rc.DB.Where("comment <> ''").Where("user_id = ?", c.MustGet("user_id").(string)).Order("created_at DESC").Find(&reviews)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Reviews fetched successfully", reviews))
}

func (rc *ReviewController) GetReviewsByUserId(c *gin.Context) {
	var reviews []models.Review
	result := rc.DB.Where("user_id = ?", c.Param("id")).Where("comment <> ''").Order("created_at DESC").Find(&reviews)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Reviews fetched successfully", reviews))
}

func (rc *ReviewController) GetReviewsByExpertId(c *gin.Context) {
	var reviews []models.Review
	result := rc.DB.Where("expert_id = ?", c.Param("id")).Where("comment <> ''").Order("created_at DESC").Find(&reviews)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Reviews fetched successfully", reviews))
}
