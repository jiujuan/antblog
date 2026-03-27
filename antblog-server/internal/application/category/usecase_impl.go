package category

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	domain "antblog/internal/domain/category"
	apperrors "antblog/pkg/errors"
)

// ─── 依赖声明 ────────────────────────────────────────────────────────────────

// Deps fx 注入依赖
type Deps struct {
	fx.In
	Repo          domain.ICategoryRepository
	DomainService domain.IDomainService
	Logger        *zap.Logger
}

// categoryUseCase 分类用例实现
type categoryUseCase struct {
	repo   domain.ICategoryRepository
	svc    domain.IDomainService
	logger *zap.Logger
}

// NewCategoryUseCase 创建分类用例（fx provider）
func NewCategoryUseCase(deps Deps) ICategoryUseCase {
	return &categoryUseCase{
		repo:   deps.Repo,
		svc:    deps.DomainService,
		logger: deps.Logger,
	}
}

// ─── 前台接口实现 ─────────────────────────────────────────────────────────────

// ListCategories 获取所有分类（前台展示用，按排序权重降序）
func (uc *categoryUseCase) ListCategories(ctx context.Context) ([]*CategoryResp, error) {
	list, err := uc.repo.FindAll(ctx)
	if err != nil {
		uc.logger.Error("list categories failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}
	return toCategoryRespList(list), nil
}

// GetCategoryBySlug 前台按 slug 获取分类
func (uc *categoryUseCase) GetCategoryBySlug(ctx context.Context, slug string) (*CategoryResp, error) {
	c, err := uc.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrCategoryNotFound()
	}
	return toCategoryResp(c), nil
}

// ─── 后台管理接口实现 ──────────────────────────────────────────────────────────

// GetCategoryByID 后台按 ID 获取分类
func (uc *categoryUseCase) GetCategoryByID(ctx context.Context, id uint64) (*CategoryResp, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrCategoryNotFound()
	}
	return toCategoryResp(c), nil
}

// CreateCategory 创建分类
func (uc *categoryUseCase) CreateCategory(ctx context.Context, req *CreateCategoryReq) (*CategoryResp, error) {
	// 1. 领域服务校验名称唯一性
	if err := uc.svc.ValidateCreate(ctx, req.Name, req.Slug); err != nil {
		return nil, err
	}

	// 2. 自动生成/确保 Slug 唯一
	slug, err := uc.svc.EnsureSlug(ctx, req.Name, req.Slug, 0)
	if err != nil {
		return nil, err
	}

	// 3. 构建领域实体
	c, err := uc.svc.BuildCategory(req.Name, slug, req.Description, req.Cover, req.SortOrder)
	if err != nil {
		return nil, err
	}

	// 4. 持久化
	created, err := uc.repo.Create(ctx, c)
	if err != nil {
		uc.logger.Error("create category failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	uc.logger.Info("category created",
		zap.Uint64("id", created.ID),
		zap.String("name", created.Name),
		zap.String("slug", created.Slug),
	)
	return toCategoryResp(created), nil
}

// UpdateCategory 更新分类
func (uc *categoryUseCase) UpdateCategory(ctx context.Context, id uint64, req *UpdateCategoryReq) (*CategoryResp, error) {
	// 1. 查询是否存在
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrCategoryNotFound()
	}

	// 2. 领域校验唯一性（排除自身）
	if err = uc.svc.ValidateUpdate(ctx, id, req.Name, req.Slug); err != nil {
		return nil, err
	}

	// 3. 自动生成/确保 Slug 唯一（排除自身）
	slug, err := uc.svc.EnsureSlug(ctx, req.Name, req.Slug, id)
	if err != nil {
		return nil, err
	}

	// 4. 调用领域方法更新
	c.Update(req.Name, slug, req.Description, req.Cover, req.SortOrder)

	// 5. 持久化
	if err = uc.repo.Update(ctx, c); err != nil {
		uc.logger.Error("update category failed", zap.Uint64("id", id), zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	uc.logger.Info("category updated", zap.Uint64("id", id), zap.String("name", c.Name))
	return toCategoryResp(c), nil
}

// DeleteCategory 删除分类（检查是否有文章绑定）
func (uc *categoryUseCase) DeleteCategory(ctx context.Context, id uint64) error {
	// 1. 查询是否存在
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.ErrCategoryNotFound()
	}

	// 2. 若存在关联文章，拒绝删除
	if c.ArticleCount > 0 {
		return apperrors.FromCode(apperrors.CodeCategoryHasArticle)
	}

	// 3. 软删除
	if err = uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("delete category failed", zap.Uint64("id", id), zap.Error(err))
		return apperrors.ErrInternalError(err)
	}

	uc.logger.Info("category deleted", zap.Uint64("id", id), zap.String("name", c.Name))
	return nil
}

// ─── 映射函数 ────────────────────────────────────────────────────────────────

func toCategoryResp(c *domain.Category) *CategoryResp {
	return &CategoryResp{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		Cover:        c.Cover,
		SortOrder:    c.SortOrder,
		ArticleCount: c.ArticleCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func toCategoryRespList(list []*domain.Category) []*CategoryResp {
	resp := make([]*CategoryResp, 0, len(list))
	for _, c := range list {
		resp = append(resp, toCategoryResp(c))
	}
	return resp
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 分类应用模块
var Module = fx.Options(
	fx.Provide(
		NewCategoryUseCase,
		domain.NewDomainService,
	),
)
