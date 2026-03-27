// Package category 分类应用层。
package category

import "time"

// ─── 请求 DTO ────────────────────────────────────────────────────────────────

// CreateCategoryReq 创建分类请求
type CreateCategoryReq struct {
	Name        string `json:"name"        validate:"required,max=64"`
	Slug        string `json:"slug"        validate:"omitempty,max=128"`
	Description string `json:"description" validate:"max=512"`
	Cover       string `json:"cover"       validate:"omitempty,max=512,url"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryReq 更新分类请求
type UpdateCategoryReq struct {
	Name        string `json:"name"        validate:"required,max=64"`
	Slug        string `json:"slug"        validate:"omitempty,max=128"`
	Description string `json:"description" validate:"max=512"`
	Cover       string `json:"cover"       validate:"omitempty,max=512,url"`
	SortOrder   int    `json:"sort_order"`
}

// ─── 响应 DTO ────────────────────────────────────────────────────────────────

// CategoryResp 分类信息响应
type CategoryResp struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	Cover        string    `json:"cover"`
	SortOrder    int       `json:"sort_order"`
	ArticleCount int       `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
