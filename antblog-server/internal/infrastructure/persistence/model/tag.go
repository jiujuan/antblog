package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag GORM 标签模型
type Tag struct {
	ID           uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string         `gorm:"column:name;type:varchar(64);not null;uniqueIndex"`
	Slug         string         `gorm:"column:slug;type:varchar(128);not null;uniqueIndex"`
	Color        string         `gorm:"column:color;type:varchar(16);not null;default:'#6B7280'"`
	ArticleCount int            `gorm:"column:article_count;not null;default:0"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Tag) TableName() string { return "tags" }

// ArticleTag GORM 文章-标签关联模型
type ArticleTag struct {
	ArticleID uint64    `gorm:"column:article_id;primaryKey"`
	TagID     uint64    `gorm:"column:tag_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ArticleTag) TableName() string { return "article_tags" }
