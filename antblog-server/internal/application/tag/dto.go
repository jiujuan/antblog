// Package tag 标签应用层。
package tag

import "time"

// ─── 请求 DTO ────────────────────────────────────────────────────────────────

// CreateTagReq 创建标签请求
type CreateTagReq struct {
	Name  string `json:"name"  validate:"required,max=64"`
	Slug  string `json:"slug"  validate:"omitempty,max=128"`
	Color string `json:"color" validate:"omitempty,max=16"`
}

// UpdateTagReq 更新标签请求
type UpdateTagReq struct {
	Name  string `json:"name"  validate:"required,max=64"`
	Slug  string `json:"slug"  validate:"omitempty,max=128"`
	Color string `json:"color" validate:"omitempty,max=16"`
}

// BatchCreateTagReq 批量创建标签请求（文章编辑场景）
type BatchCreateTagReq struct {
	Tags []CreateTagReq `json:"tags" validate:"required,min=1,max=20,dive"`
}

// ─── 响应 DTO ────────────────────────────────────────────────────────────────

// TagResp 标签信息响应
type TagResp struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Color        string    `json:"color"`
	ArticleCount int       `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TagSimpleResp 标签精简响应（用于文章详情嵌入展示）
type TagSimpleResp struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Color string `json:"color"`
}
