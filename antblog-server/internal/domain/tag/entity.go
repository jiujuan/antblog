// Package tag 标签领域层 —— 纯业务实体，零外部依赖。
package tag

import "time"

// Tag 文章标签聚合根
type Tag struct {
	ID           uint64
	Name         string
	Slug         string    // URL 友好标识，全局唯一
	Color        string    // 展示颜色（Hex 色值，如 #00ADD8）
	ArticleCount int       // 关联已发布文章数（冗余字段）
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Update 更新标签基本信息
func (t *Tag) Update(name, slug, color string) {
	if name != "" {
		t.Name = name
	}
	if slug != "" {
		t.Slug = slug
	}
	if color != "" {
		t.Color = color
	}
	t.UpdatedAt = time.Now()
}

// IncrArticleCount 增加文章计数
func (t *Tag) IncrArticleCount() {
	t.ArticleCount++
}

// DecrArticleCount 减少文章计数（最小为 0）
func (t *Tag) DecrArticleCount() {
	if t.ArticleCount > 0 {
		t.ArticleCount--
	}
}
