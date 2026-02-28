package services

import (
	"net/http"
	"time"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GenerateCallToken(db *gorm.DB, callID string, appID string, appCertificate string, isUser bool, expertID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		var call models.Call
		result := db.First(&call, "id = ?", callID)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
			return
		}

		if isUser == true && call.UserID != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse("Forbidden", nil))
			return
		} else {
			if call.ExpertID != expertID {
				c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse("Forbidden", nil))
				return
			}
		}

		if call.Status != models.CallStatusAccepted {
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse("Forbidden", nil))
			return
		}

		now := time.Now().UTC()

		if now.Before(call.ScheduledAt.Add(-10 * time.Minute)) {
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse("Too early to join call", nil))
			return
		}

		expireTimeInSeconds := uint32(3600)
		currentTimestamp := uint32(time.Now().Unix())
		expireTimestamp := currentTimestamp + expireTimeInSeconds

		token, err := rtctokenbuilder.BuildTokenWithUid(
			appID,
			appCertificate,
			call.ChannelName,
			0,
			rtctokenbuilder.RolePublisher,
			expireTimestamp,
			expireTimestamp,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse("Call token generated successfully", map[string]string{
			"channel": call.ChannelName,
			"token":   token,
			"appId":   appID,
		}))

	}
}
