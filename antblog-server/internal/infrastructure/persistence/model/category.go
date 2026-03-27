package model

import (
	"time"

	"gorm.io/gorm"
)

// Category GORM 分类模型
type Category struct {
	ID           uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string         `gorm:"column:name;type:varchar(64);not null;uniqueIndex"`
	Slug         string         `gorm:"column:slug;type:varchar(128);not null;uniqueIndex"`
	Description  string         `gorm:"column:description;type:varchar(512);not null;default:''"`
	Cover        string         `gorm:"column:cover;type:varchar(512);not null;default:''"`
	SortOrder    int            `gorm:"column:sort_order;not null;default:0;index:idx_sort,sort:desc"`
	ArticleCount int            `gorm:"column:article_count;not null;default:0"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Category) TableName() string { return "categories" }
