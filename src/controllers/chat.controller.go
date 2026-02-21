package controllers

import (
	"net/http"
	"sort"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatController struct {
	DB *gorm.DB
}

func NewChatController(db *gorm.DB) *ChatController {
	return &ChatController{DB: db}
}

func (cc *ChatController) CreateChat(c *gin.Context) {
	var payload models.CreateChatRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var chat models.Chat
	chat.UserID = c.MustGet("user_id").(string)
	chat.ExpertID = payload.ExpertID

	result := cc.DB.Create(&chat)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var message models.Message
	message.ChatID = chat.ID
	message.SenderID = c.MustGet("user_id").(string)
	message.ReceiverID = payload.ExpertID
	message.Content = payload.Message
	message.ResponseType = payload.ResponseType
	message.Attachments = []string{}
	message.IsRead = false

	result = cc.DB.Create(&message)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Chat created successfully", chat))
}

func (cc *ChatController) GetUserChats(c *gin.Context) {
	var chats []models.Chat
	result := cc.DB.Preload("User").Preload("Expert").Preload("Expert.User").
		Where("user_id = ?", c.MustGet("user_id").(string)).
		Order("created_at DESC").
		Find(&chats)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	chatIDs := make([]string, 0, len(chats))
	for _, chat := range chats {
		chatIDs = append(chatIDs, chat.ID)
	}

	lastMessageByChatID := make(map[string]models.Message)
	if len(chatIDs) > 0 {
		var lastMessages []models.Message
		result = cc.DB.Raw(
			"SELECT DISTINCT ON (chat_id) * FROM messages WHERE chat_id IN ? ORDER BY chat_id, created_at DESC",
			chatIDs,
		).Scan(&lastMessages)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
			return
		}
		for _, message := range lastMessages {
			lastMessageByChatID[message.ChatID] = message
		}
	}

	sortChatsByLastMessage(chats, lastMessageByChatID)

	chatResponses := make([]models.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		lastMessage := lastMessageByChatID[chat.ID]
		chatResponses = append(chatResponses, chat.ToChatResponse(lastMessage))
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Chats fetched successfully", chatResponses))
}

func (cc *ChatController) GetExpertChats(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var expert models.ExpertProfile
	result := cc.DB.First(&expert, "user_id = ?", userID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var chats []models.Chat
	result = cc.DB.Preload("User").Preload("Expert").Where("expert_id = ?", expert.ID).Order("created_at DESC").Find(&chats)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	chatIDs := make([]string, 0, len(chats))
	for _, chat := range chats {
		chatIDs = append(chatIDs, chat.ID)
	}

	lastMessageByChatID := make(map[string]models.Message)
	if len(chatIDs) > 0 {
		var lastMessages []models.Message
		result = cc.DB.Raw(
			"SELECT DISTINCT ON (chat_id) * FROM messages WHERE chat_id IN ? ORDER BY chat_id, created_at DESC",
			chatIDs,
		).Scan(&lastMessages)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
			return
		}
		for _, message := range lastMessages {
			lastMessageByChatID[message.ChatID] = message
		}
	}

	sortChatsByLastMessage(chats, lastMessageByChatID)

	chatResponses := make([]models.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		lastMessage := lastMessageByChatID[chat.ID]
		chatResponses = append(chatResponses, chat.ToChatResponse(lastMessage))
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Chats fetched successfully", chatResponses))
}

func (cc *ChatController) GetChatMessages(c *gin.Context) {
	chatID := c.Param("id")

	var messages []models.Message
	result := cc.DB.Where("chat_id = ?", chatID).Order("created_at ASC").Find(&messages)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Messages fetched successfully", messages))
}

func (cc *ChatController) CreateMessage(c *gin.Context) {
	chatID := c.Param("id")

	var payload models.CreateMessageRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var message models.Message
	message.ChatID = chatID
	message.SenderID = payload.SenderID
	message.ReceiverID = payload.ReceiverID
	message.Content = payload.Content
	message.ResponseType = payload.ResponseType
	message.Attachments = []string{}
	message.IsRead = false

	result := cc.DB.Create(&message)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Message created successfully", message))
}

func sortChatsByLastMessage(chats []models.Chat, lastMessageByChatID map[string]models.Message) {
	sort.SliceStable(chats, func(i, j int) bool {
		leftMessage, leftOk := lastMessageByChatID[chats[i].ID]
		rightMessage, rightOk := lastMessageByChatID[chats[j].ID]

		if leftOk && rightOk {
			return leftMessage.CreatedAt.After(rightMessage.CreatedAt)
		}
		if leftOk != rightOk {
			return leftOk
		}
		return chats[i].CreatedAt.After(chats[j].CreatedAt)
	})
}
