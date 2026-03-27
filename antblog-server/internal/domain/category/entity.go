// Package category 分类领域层 —— 纯业务实体，零外部依赖。
package category

import "time"

// Category 文章分类聚合根
type Category struct {
	ID           uint64
	Name         string
	Slug         string // URL 友好标识，全局唯一
	Description  string
	Cover        string // 封面图 URL
	SortOrder    int    // 排序权重，越大越靠前
	ArticleCount int    // 文章数量（冗余字段，加速列表查询）
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Update 更新分类基本信息
func (c *Category) Update(name, slug, description, cover string, sortOrder int) {
	if name != "" {
		c.Name = name
	}
	if slug != "" {
		c.Slug = slug
	}
	c.Description = description
	c.Cover = cover
	c.SortOrder = sortOrder
	c.UpdatedAt = time.Now()
}

// IncrArticleCount 增加文章计数
func (c *Category) IncrArticleCount() {
	c.ArticleCount++
}

// DecrArticleCount 减少文章计数（最小为 0）
func (c *Category) DecrArticleCount() {
	if c.ArticleCount > 0 {
		c.ArticleCount--
	}
}
