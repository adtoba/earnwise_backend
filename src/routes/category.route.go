package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/gin-gonic/gin"
)

type CategoryRouteController struct {
	categoryController controllers.CategoryController
}

func NewCategoryRouteController(categoryController controllers.CategoryController) *CategoryRouteController {
	return &CategoryRouteController{categoryController: categoryController}
}

func (rc *CategoryRouteController) RegisterCategoryRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/categories")
	router.GET("/", rc.categoryController.GetAllCategories)
	router.GET("/:id", rc.categoryController.GetCategoryById)
	router.POST("/", rc.categoryController.CreateCategory)
	router.POST("/bulk", rc.categoryController.CreateCategories)
}
