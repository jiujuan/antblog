package category

import "context"

// ICategoryUseCase 分类用例接口
type ICategoryUseCase interface {
	// ── 前台接口 ─────────────────────────────────────────────────────────

	// ListCategories 获取所有分类列表（按 sort_order 降序）
	ListCategories(ctx context.Context) ([]*CategoryResp, error)

	// GetCategoryBySlug 按 Slug 获取分类详情
	GetCategoryBySlug(ctx context.Context, slug string) (*CategoryResp, error)

	// ── 后台管理接口 ──────────────────────────────────────────────────────

	// GetCategoryByID 按 ID 获取分类（管理员）
	GetCategoryByID(ctx context.Context, id uint64) (*CategoryResp, error)

	// CreateCategory 创建分类
	CreateCategory(ctx context.Context, req *CreateCategoryReq) (*CategoryResp, error)

	// UpdateCategory 更新分类
	UpdateCategory(ctx context.Context, id uint64, req *UpdateCategoryReq) (*CategoryResp, error)

	// DeleteCategory 删除分类（需检查是否有关联文章）
	DeleteCategory(ctx context.Context, id uint64) error
}
