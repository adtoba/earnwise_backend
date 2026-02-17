package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID            string      `json:"id" gorm:"primaryKey"`
	ExpertID      string      `json:"expert_id"`
	UserID        string      `json:"user_id"`
	User          User        `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Content       string      `json:"content"`
	Attachments   StringArray `json:"attachments"`
	LikesCount    int         `json:"likes_count"`
	CommentsCount int         `json:"comments_count"`
	IsPostLiked   bool        `json:"is_post_liked"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}

type CreatePostRequest struct {
	ExpertID    string      `json:"expert_id" binding:"required"`
	Content     string      `json:"content" binding:"required"`
	Attachments StringArray `json:"attachments"`
}

type Comment struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	PostID     string    `json:"post_id"`
	UserID     string    `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Comment    string    `json:"comment"`
	LikesCount int       `json:"likes_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return
}

type CreateCommentRequest struct {
	PostID  string `json:"post_id" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

type LikePostRequest struct {
	PostID string `json:"post_id" binding:"required"`
}

type LikeCommentRequest struct {
	CommentID string `json:"comment_id" binding:"required"`
}
