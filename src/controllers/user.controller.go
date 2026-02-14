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

	res := uc.DB.Save(&user)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User updated successfully", user.ToUserResponse()))
}
