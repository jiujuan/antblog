package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment GORM 评论模型
type Comment struct {
	ID        uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	ArticleID uint64         `gorm:"column:article_id;not null;index:idx_comments_article"`
	UserID    *uint64        `gorm:"column:user_id;index:idx_comments_user"`
	ParentID  *uint64        `gorm:"column:parent_id;index:idx_comments_parent"`
	RootID    *uint64        `gorm:"column:root_id;index:idx_comments_root"`
	ReplyToID *uint64        `gorm:"column:reply_to_id"`
	Nickname  string         `gorm:"column:nickname;type:varchar(64);not null;default:''"`
	Email     string         `gorm:"column:email;type:varchar(128);not null;default:''"`
	Content   string         `gorm:"column:content;type:text;not null"`
	IP        string         `gorm:"column:ip;type:varchar(64);not null;default:''"`
	UserAgent string         `gorm:"column:user_agent;type:varchar(512);not null;default:''"`
	Status    int8           `gorm:"column:status;not null;default:1;index:idx_comments_status"`
	LikeCount int            `gorm:"column:like_count;not null;default:0"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_comments_deleted"`
}

func (Comment) TableName() string { return "comments" }
