package controllers

import (
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{DB: db}
}

func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var payload models.CreateCategoryRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	category := models.Category{
		Name:        payload.Name,
		Description: payload.Description,
	}

	result := cc.DB.Create(&category)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Category created successfully", category))
}

func (cc *CategoryController) CreateCategories(c *gin.Context) {
	var payload []models.CreateCategoryRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	if len(payload) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", "categories list is empty"))
		return
	}

	categories := make([]models.Category, 0, len(payload))
	for _, item := range payload {
		categories = append(categories, models.Category{
			Name:        item.Name,
			Description: item.Description,
		})
	}

	result := cc.DB.Create(&categories)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Categories created successfully", categories))
}

func (cc *CategoryController) GetAllCategories(c *gin.Context) {
	var categories []models.Category
	result := cc.DB.Find(&categories)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Categories fetched successfully", categories))
}

func (cc *CategoryController) GetCategoryById(c *gin.Context) {
	var category models.Category
	result := cc.DB.First(&category, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Category fetched successfully", category))
}
