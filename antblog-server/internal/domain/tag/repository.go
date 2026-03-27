package tag

import "context"

// ITagRepository 标签仓储接口（由 infrastructure 层实现）
type ITagRepository interface {
	// ── 单条操作 ─────────────────────────────────────────────────────────────

	// Create 创建标签
	Create(ctx context.Context, t *Tag) (*Tag, error)

	// FindByID 按 ID 查询
	FindByID(ctx context.Context, id uint64) (*Tag, error)

	// FindBySlug 按 Slug 查询
	FindBySlug(ctx context.Context, slug string) (*Tag, error)

	// FindByName 按名称精确查询
	FindByName(ctx context.Context, name string) (*Tag, error)

	// Update 更新标签
	Update(ctx context.Context, t *Tag) error

	// Delete 软删除标签
	Delete(ctx context.Context, id uint64) error

	// ── 批量查询 ─────────────────────────────────────────────────────────────

	// FindAll 查询所有标签（按文章数降序）
	FindAll(ctx context.Context) ([]*Tag, error)

	// FindByIDs 按 ID 批量查询（用于文章写入时关联标签）
	FindByIDs(ctx context.Context, ids []uint64) ([]*Tag, error)

	// FindByArticleID 查询某篇文章关联的所有标签
	FindByArticleID(ctx context.Context, articleID uint64) ([]*Tag, error)

	// ── 唯一性校验 ───────────────────────────────────────────────────────────

	// ExistsBySlug 检查 Slug 是否已存在（excludeID=0 表示不排除任何记录）
	ExistsBySlug(ctx context.Context, slug string, excludeID uint64) (bool, error)

	// ExistsByName 检查名称是否已存在（excludeID=0 表示不排除任何记录）
	ExistsByName(ctx context.Context, name string, excludeID uint64) (bool, error)

	// ── 计数管理 ─────────────────────────────────────────────────────────────

	// IncrArticleCount 原子性增减文章计数（delta 可为负数）
	IncrArticleCount(ctx context.Context, id uint64, delta int) error

	// UpdateArticleCount 直接设置文章计数（批量同步用）
	UpdateArticleCount(ctx context.Context, id uint64, count int) error
}
