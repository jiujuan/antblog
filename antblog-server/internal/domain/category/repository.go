package category

import "context"

// ICategoryRepository 分类仓储接口（由 infrastructure 层实现）
type ICategoryRepository interface {
	// Create 创建分类
	Create(ctx context.Context, c *Category) (*Category, error)

	// FindByID 按 ID 查询分类
	FindByID(ctx context.Context, id uint64) (*Category, error)

	// FindBySlug 按 Slug 查询分类
	FindBySlug(ctx context.Context, slug string) (*Category, error)

	// FindAll 查询所有分类（按 sort_order 降序）
	FindAll(ctx context.Context) ([]*Category, error)

	// Update 更新分类
	Update(ctx context.Context, c *Category) error

	// Delete 软删除分类
	Delete(ctx context.Context, id uint64) error

	// ExistsBySlug 检查 Slug 是否已存在（排除指定 ID）
	ExistsBySlug(ctx context.Context, slug string, excludeID uint64) (bool, error)

	// ExistsByName 检查名称是否已存在（排除指定 ID）
	ExistsByName(ctx context.Context, name string, excludeID uint64) (bool, error)

	// UpdateArticleCount 直接设置文章计数（批量同步用）
	UpdateArticleCount(ctx context.Context, id uint64, count int) error

	// IncrArticleCount 原子性增加文章计数
	IncrArticleCount(ctx context.Context, id uint64, delta int) error
}
