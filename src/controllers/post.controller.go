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
	if payload.Attachments != nil {
		post.Attachments = payload.Attachments
	} else {
		post.Attachments = []string{}
	}
	post.LikesCount = 0
	post.CommentsCount = 0

	result := pc.DB.Preload("User").Create(&post)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Post created successfully", post))
}

func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post
	result := pc.DB.Preload("User").Where("user_id = ?", c.MustGet("user_id").(string)).Order("created_at DESC").Find(&posts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Posts fetched successfully", posts))
}

func (pc *PostController) GetRandomPosts(c *gin.Context) {
	var posts []models.Post
	currentUserID := c.MustGet("user_id").(string)
	result := pc.DB.Preload("User").
		Where("user_id <> ?", currentUserID).
		Order("RANDOM()").
		Limit(10).
		Find(&posts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Random posts fetched successfully", posts))
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

func (pc *PostController) GetCommentsByPostId(c *gin.Context) {
	var comments []models.Comment
	result := pc.DB.Preload("User").Where("post_id = ?", c.Param("id")).Find(&comments)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Comments fetched successfully", comments))
}

func (pc *PostController) CreateComment(c *gin.Context) {
	var payload models.CreateCommentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var comment models.Comment
	comment.PostID = payload.PostID
	comment.UserID = c.MustGet("user_id").(string)
	comment.Comment = payload.Comment
	comment.LikesCount = 0

	result := pc.DB.Create(&comment).Preload("User").First(&comment)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	var post models.Post
	result = pc.DB.First(&post, "id = ?", payload.PostID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}
	post.CommentsCount++
	pc.DB.Save(&post)

	c.JSON(http.StatusOK, models.SuccessResponse("Comment created successfully", comment))
}

func (pc *PostController) LikePost(c *gin.Context) {
	var payload models.LikePostRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var post models.Post
	result := pc.DB.First(&post, "id = ?", payload.PostID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	post.LikesCount++

	res := pc.DB.Save(&post)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", res.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Post liked successfully", post))
}

func (pc *PostController) LikeComment(c *gin.Context) {
	var payload models.LikeCommentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request payload", err.Error()))
		return
	}

	var comment models.Comment
	result := pc.DB.First(&comment, "id = ?", payload.CommentID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", result.Error.Error()))
		return
	}

	comment.LikesCount++

	res := pc.DB.Save(&comment)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal server error", res.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Comment liked successfully", comment))
}
