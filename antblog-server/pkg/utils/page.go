// Package utils 提供通用工具函数集合。
package utils

import "math"

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// PageQuery 分页查询参数
type PageQuery struct {
	Page     int `form:"page"      json:"page"      validate:"min=1"`
	PageSize int `form:"page_size" json:"page_size" validate:"min=1,max=100"`
}

// Normalize 修正分页参数为合法值
func (p *PageQuery) Normalize() {
	if p.Page <= 0 {
		p.Page = DefaultPage
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
}

// Offset 计算 SQL OFFSET
func (p *PageQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// TotalPages 根据总数计算总页数
func (p *PageQuery) TotalPages(total int64) int {
	if p.PageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(p.PageSize)))
}

// HasNextPage 是否有下一页
func (p *PageQuery) HasNextPage(total int64) bool {
	return int64(p.Page*p.PageSize) < total
}

// ─── 独立工具函数 ────────────────────────────────────────────────────────────

// CalcOffset 直接计算偏移量
func CalcOffset(page, pageSize int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * pageSize
}

// NormalizePage 修正 page 值
func NormalizePage(page int) int {
	if page <= 0 {
		return DefaultPage
	}
	return page
}

// NormalizePageSize 修正 pageSize 值
func NormalizePageSize(pageSize int) int {
	if pageSize <= 0 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}
