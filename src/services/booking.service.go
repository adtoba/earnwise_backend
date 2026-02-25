package services

import (
	"errors"
	"time"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IsSlotAvailable(
	tx *gorm.DB,
	expertID string,
	start time.Time,
	duration int,
) (bool, error) {

	end := start.Add(time.Duration(duration) * time.Minute)

	var count int64

	err := tx.Model(&models.Call{}).
		Where("expert_id = ?", expertID).
		Where("status IN ?", []string{models.CallStatusPending, models.CallStatusAccepted}).
		Where("scheduled_at < ?", end).
		Where("scheduled_at + (duration_mins || ' minutes')::interval > ?", start).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func CreateCall(
	db *gorm.DB,
	ctx *gin.Context,
	payload models.CreateCallRequest,
) (*models.Call, error) {

	var call models.Call

	err := db.Transaction(func(tx *gorm.DB) error {

		ok, err := IsSlotAvailable(tx, payload.ExpertID, payload.ScheduledAt, payload.DurationMins)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("slot already booked")
		}

		userID := ctx.MustGet("user_id").(string)

		call = models.Call{
			UserID:        userID,
			ExpertID:      payload.ExpertID,
			ScheduledAt:   payload.ScheduledAt,
			DurationMins:  payload.DurationMins,
			Status:        "pending",
			Price:         payload.Price,
			PaymentStatus: "pending",
			PaymentRef:    payload.PaymentRef,
			Subject:       payload.Subject,
			Description:   payload.Description,
		}

		if err := tx.Create(&call).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &call, nil
}
