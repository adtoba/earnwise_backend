package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"gorm.io/gorm"
)

type NotificationService struct {
	DB     *gorm.DB
	AppID  string
	APIKey string
}

func NewNotificationService(db *gorm.DB, appID string, apiKey string) *NotificationService {
	return &NotificationService{DB: db, AppID: appID, APIKey: apiKey}
}

func (ns *NotificationService) SendNotification(message string, userID string, title string) error {
	payload := map[string]interface{}{
		"headings": map[string]string{
			"en": title,
		},
		"contents": map[string]string{
			"en": message,
		},
		// High priority for Android and iOS (OneSignal supports these fields)
		"priority":       10,
		"apns_priority":  "10",
		"apns_push_type": "alert",
		"target_channel": "push",
		"app_id":         ns.AppID,
		"include_aliases": map[string]interface{}{
			"external_id": []string{userID},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := "https://api.onesignal.com/notifications?c=push"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+ns.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

func (ns *NotificationService) CreateNotification(notification models.Notification) (models.Notification, error) {
	result := ns.DB.Create(&notification)
	if result.Error != nil {
		return models.Notification{}, result.Error
	}
	return notification, nil
}
