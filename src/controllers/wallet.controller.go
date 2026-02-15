package controllers

import (
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WalletController struct {
	DB *gorm.DB
}

func NewWalletController(db *gorm.DB) *WalletController {
	return &WalletController{DB: db}
}

func (wc *WalletController) CreateWallet(c *gin.Context, expertID string) {
	var wallet models.Wallet
	wallet.UserID = c.MustGet("user_id").(string)
	wallet.ExpertID = expertID
	result := wc.DB.Create(&wallet)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
}

func (wc *WalletController) GetWallet(c *gin.Context) {
	var wallet models.Wallet
	result := wc.DB.First(&wallet, "user_id = ?", c.MustGet("user_id").(string))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Wallet fetched successfully", wallet))
}
