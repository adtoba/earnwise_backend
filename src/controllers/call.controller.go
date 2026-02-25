package controllers

import (
	"net/http"
	"time"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/adtoba/earnwise_backend/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CallController struct {
	DB *gorm.DB
}

const callAcceptCutoffMinutes = 10

func NewCallController(db *gorm.DB) *CallController {
	return &CallController{DB: db}
}

func (cc *CallController) CreateCall(c *gin.Context) {
	var payload models.CreateCallRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var expert models.ExpertProfile
	result := cc.DB.First(&expert, "id = ?", payload.ExpertID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	timezone := payload.Timezone
	if timezone == "" {
		timezone = "Africa/Lagos"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", "timezone is invalid"))
		return
	}

	scheduledLocal := time.Date(
		payload.ScheduledAt.Year(),
		payload.ScheduledAt.Month(),
		payload.ScheduledAt.Day(),
		payload.ScheduledAt.Hour(),
		payload.ScheduledAt.Minute(),
		payload.ScheduledAt.Second(),
		0,
		loc,
	)
	payload.ScheduledAt = scheduledLocal.UTC()
	payload.Price = float64(payload.DurationMins) * expert.Rates.Video

	call, err := services.CreateCall(cc.DB, c, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Call created successfully", call))
}

func (cc *CallController) GetUserCalls(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	status := c.Query("status")
	if status == "" {
		status = models.CallStatusPending
	}

	if err := cc.expirePendingCalls("user_id = ?", userID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}

	switch status {
	case models.CallStatusPending:
		status = models.CallStatusPending
	case models.CallStatusAccepted:
		status = models.CallStatusAccepted
	case models.CallStatusRejected:
		status = models.CallStatusRejected
	case models.CallStatusCancelled:
		status = models.CallStatusCancelled
	case models.CallStatusExpired:
		status = models.CallStatusExpired
	case models.CallStatusCompleted:
		status = models.CallStatusCompleted
	default:
		status = models.CallStatusPending
	}

	var calls []models.Call
	result := cc.DB.Preload("User").Preload("Expert").Preload("Expert.User").
		Where("user_id = ?", userID).
		Where("status = ?", status).
		Order("scheduled_at ASC").
		Find(&calls)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	callResponses := make([]models.CallResponse, 0, len(calls))
	for _, call := range calls {
		callResponses = append(callResponses, call.ToCallResponse())
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Calls fetched successfully", callResponses))
}

func (cc *CallController) GetExpertCalls(c *gin.Context) {
	expertID := c.MustGet("user_id").(string)
	status := c.Query("status")
	if status == "" {
		status = models.CallStatusPending
	}

	switch status {
	case models.CallStatusPending:
		status = models.CallStatusPending
	case models.CallStatusAccepted:
		status = models.CallStatusAccepted
	case models.CallStatusRejected:
		status = models.CallStatusRejected
	case models.CallStatusCancelled:
		status = models.CallStatusCancelled
	case models.CallStatusExpired:
		status = models.CallStatusExpired
	case models.CallStatusCompleted:
		status = models.CallStatusCompleted
	default:
		status = models.CallStatusPending
	}

	var expert models.ExpertProfile
	result := cc.DB.First(&expert, "user_id = ?", expertID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	if err := cc.expirePendingCalls("expert_id = ?", expert.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}

	var calls []models.Call
	result = cc.DB.Preload("User").Preload("Expert").Preload("Expert.User").
		Where("expert_id = ?", expert.ID).
		Where("status = ?", status).
		Order("scheduled_at ASC").
		Find(&calls)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	callResponses := make([]models.CallResponse, 0, len(calls))
	for _, call := range calls {
		callResponses = append(callResponses, call.ToCallResponse())
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Calls fetched successfully", callResponses))

}

func (cc *CallController) AcceptCall(c *gin.Context) {
	callID := c.Param("id")
	var call models.Call
	result := cc.DB.First(&call, "id = ?", callID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	if call.Status != models.CallStatusPending {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Call not pending", nil))
		return
	}
	if call.ScheduledAt.Before(time.Now().UTC()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Call is in the past", nil))
		return
	}

	res := cc.DB.Model(&call).Updates(map[string]interface{}{
		"status": models.CallStatusAccepted,
	})
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", res.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Call accepted successfully", call.ToCallResponse()))
}

func (cc *CallController) expirePendingCalls(where string, args ...interface{}) error {
	cutoffTime := time.Now().UTC().Add(time.Duration(callAcceptCutoffMinutes) * time.Minute)
	result := cc.DB.Model(&models.Call{}).
		Where("status = ?", models.CallStatusPending).
		Where("scheduled_at <= ?", cutoffTime).
		Where(where, args...).
		Updates(map[string]interface{}{
			"status":         models.CallStatusExpired,
			"payment_status": models.PaymentStatusReleased,
		})
	return result.Error
}
