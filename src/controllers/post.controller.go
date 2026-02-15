package controllers

import (
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(db *gorm.DB) *PostController {
	return &PostController{DB: db}
}

func (pc *PostController) CreatePost(c *gin.Context) {
	var payload models.CreatePostRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var post models.Post
	post.ExpertID = payload.ExpertID
	post.UserID = c.MustGet("user_id").(string)
	post.Content = payload.Content
	post.Attachments = payload.Attachments
	post.LikesCount = 0
	post.CommentsCount = 0

	result := pc.DB.Create(&post)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Post created successfully", post))
}

func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post
	result := pc.DB.Preload("User").Where("user_id = ?", c.MustGet("user_id").(string)).Find(&posts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Posts fetched successfully", posts))
}

func (pc *PostController) GetPostById(c *gin.Context) {
	var post models.Post
	result := pc.DB.Preload("User").Where("id = ?", c.Param("id")).First(&post)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Post fetched successfully", post))
}

func (pc *PostController) GetPostsByExpertId(c *gin.Context) {
	var posts []models.Post
	result := pc.DB.Preload("User").Where("expert_id = ?", c.Param("id")).Find(&posts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Posts fetched successfully", posts))
}
