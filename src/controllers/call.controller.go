package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/adtoba/earnwise_backend/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CallController struct {
	DB                  *gorm.DB
	NotificationService *services.NotificationService
	AgoraAppID          string
	AgoraAppCertificate string
}

const callAcceptCutoffMinutes = 10
const minCallLeadMinutes = 10

func NewCallController(db *gorm.DB, notificationService *services.NotificationService, agoraAppID string, agoraAppCertificate string) *CallController {
	return &CallController{DB: db, NotificationService: notificationService, AgoraAppID: agoraAppID, AgoraAppCertificate: agoraAppCertificate}
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
	nowLocal := time.Now().In(loc)
	minAllowed := nowLocal.Add(time.Duration(minCallLeadMinutes) * time.Minute)
	if scheduledLocal.Before(minAllowed) {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Call scheduled too soon", fmt.Sprintf("call must be scheduled at least %d minutes in advance", minCallLeadMinutes)))
		return
	}
	payload.ScheduledAt = scheduledLocal.UTC()
	payload.Price = float64(payload.DurationMins) * expert.Rates.Video

	call, err := services.CreateCall(cc.DB, c, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", err.Error()))
		return
	}

	fmt.Println(expert.UserID)

	message := "You have a new call request scheduled for " + scheduledLocal.Format("2006-01-02 15:04:05") + ". Go to your dashboard to accept or reject now!"

	cc.NotificationService.SendNotification(message, expert.UserID, "New Call Alert")

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

	if status != "past" {
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
	}

	var calls []models.Call
	query := cc.DB.Preload("User").Preload("Expert").Preload("Expert.User").
		Where("user_id = ?", userID)
	if status == "past" {
		query = query.
			Where("scheduled_at + (duration_mins * interval '1 minute') < ?", time.Now().UTC()).
			Order("scheduled_at DESC")
	} else {
		query = query.
			Where("status = ?", status).
			Order("scheduled_at ASC")
		if status == models.CallStatusAccepted {
			query = query.Where("scheduled_at + (duration_mins * interval '1 minute') >= ?", time.Now().UTC())
		}
	}
	result := query.Find(&calls)
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

	if status != "past" {
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
	query := cc.DB.Preload("User").Preload("Expert").Preload("Expert.User").
		Where("expert_id = ?", expert.ID)
	if status == "past" {
		query = query.
			Where("scheduled_at + (duration_mins * interval '1 minute') < ?", time.Now().UTC()).
			Order("scheduled_at DESC")
	} else {
		query = query.
			Where("status = ?", status).
			Order("scheduled_at ASC")
		if status == models.CallStatusAccepted {
			query = query.Where("scheduled_at + (duration_mins * interval '1 minute') >= ?", time.Now().UTC())
		}
	}
	result = query.Find(&calls)
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

	res := cc.DB.Model(&models.Call{}).Where("id = ?", call.ID).Updates(map[string]interface{}{
		"status": models.CallStatusAccepted,
	})
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", res.Error.Error()))
		return
	}

	result = cc.DB.Preload("Expert.User").First(&call, "id = ?", call.ID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	message := "Your call request with " +
		call.Expert.User.FirstName + " " + call.Expert.User.LastName + " " +
		" has been accepted. Please prepare for the call at " + call.ScheduledAt.Format("2006-01-02 15:04:05") + ". Go to your dashboard to view the call details."

	cc.NotificationService.SendNotification(message, call.UserID, "Call Request Accepted")

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

func (cc *CallController) GenerateCallToken(c *gin.Context) {
	callID := c.Param("id")
	var payload models.GenerateTokenRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}
	services.GenerateCallToken(cc.DB, callID, cc.AgoraAppID, cc.AgoraAppCertificate, payload.IsUser, payload.ExpertID)(c)
}
