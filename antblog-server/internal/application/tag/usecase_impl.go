package tag

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	domain "antblog/internal/domain/tag"
	apperrors "antblog/pkg/errors"
)

// ─── 依赖声明 ────────────────────────────────────────────────────────────────

// Deps fx 注入依赖
type Deps struct {
	fx.In
	Repo          domain.ITagRepository
	DomainService domain.IDomainService
	Logger        *zap.Logger
}

// tagUseCase 标签用例实现
type tagUseCase struct {
	repo   domain.ITagRepository
	svc    domain.IDomainService
	logger *zap.Logger
}

// NewTagUseCase 创建标签用例（fx provider）
func NewTagUseCase(deps Deps) ITagUseCase {
	return &tagUseCase{
		repo:   deps.Repo,
		svc:    deps.DomainService,
		logger: deps.Logger,
	}
}

// ─── 前台接口实现 ─────────────────────────────────────────────────────────────

// ListTags 获取所有标签列表（按文章数降序）
func (uc *tagUseCase) ListTags(ctx context.Context) ([]*TagResp, error) {
	list, err := uc.repo.FindAll(ctx)
	if err != nil {
		uc.logger.Error("list tags failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}
	return toTagRespList(list), nil
}

// GetTagBySlug 前台按 Slug 获取标签
func (uc *tagUseCase) GetTagBySlug(ctx context.Context, slug string) (*TagResp, error) {
	t, err := uc.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrTagNotFound()
	}
	return toTagResp(t), nil
}

// GetTagsByArticleID 获取某篇文章的所有标签（精简版）
func (uc *tagUseCase) GetTagsByArticleID(ctx context.Context, articleID uint64) ([]*TagSimpleResp, error) {
	list, err := uc.repo.FindByArticleID(ctx, articleID)
	if err != nil {
		uc.logger.Error("get tags by article failed",
			zap.Uint64("article_id", articleID), zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}
	return toTagSimpleRespList(list), nil
}

// ─── 后台管理接口实现 ──────────────────────────────────────────────────────────

// GetTagByID 后台按 ID 获取标签
func (uc *tagUseCase) GetTagByID(ctx context.Context, id uint64) (*TagResp, error) {
	t, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrTagNotFound()
	}
	return toTagResp(t), nil
}

// CreateTag 创建单个标签
func (uc *tagUseCase) CreateTag(ctx context.Context, req *CreateTagReq) (*TagResp, error) {
	// 1. 领域校验：名称/Slug 唯一性 + 颜色格式
	if err := uc.svc.ValidateCreate(ctx, req.Name, req.Slug, req.Color); err != nil {
		return nil, err
	}

	// 2. 自动生成/确保 Slug 唯一
	slug, err := uc.svc.EnsureSlug(ctx, req.Name, req.Slug, 0)
	if err != nil {
		return nil, err
	}

	// 3. 构建领域实体（含颜色规范化）
	t, err := uc.svc.BuildTag(req.Name, slug, req.Color)
	if err != nil {
		return nil, err
	}

	// 4. 持久化
	created, err := uc.repo.Create(ctx, t)
	if err != nil {
		uc.logger.Error("create tag failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	uc.logger.Info("tag created",
		zap.Uint64("id", created.ID),
		zap.String("name", created.Name),
		zap.String("slug", created.Slug),
	)
	return toTagResp(created), nil
}

// BatchCreateTags 批量创建标签（幂等：已存在的直接返回，不报错）
// 典型使用场景：文章编辑器提交新标签时，同时创建不存在的标签
func (uc *tagUseCase) BatchCreateTags(ctx context.Context, req *BatchCreateTagReq) ([]*TagResp, error) {
	results := make([]*TagResp, 0, len(req.Tags))

	for i := range req.Tags {
		item := &req.Tags[i]

		// 先尝试按名称查找，存在则直接复用
		existing, err := uc.repo.FindByName(ctx, item.Name)
		if err == nil && existing != nil {
			results = append(results, toTagResp(existing))
			continue
		}

		// 不存在则创建
		created, err := uc.CreateTag(ctx, item)
		if err != nil {
			// 并发场景下可能出现重复创建竞争，再尝试查一次
			if t, findErr := uc.repo.FindByName(ctx, item.Name); findErr == nil {
				results = append(results, toTagResp(t))
				continue
			}
			uc.logger.Warn("batch create tag failed, skipping",
				zap.String("name", item.Name), zap.Error(err))
			continue
		}
		results = append(results, created)
	}

	return results, nil
}

// UpdateTag 更新标签
func (uc *tagUseCase) UpdateTag(ctx context.Context, id uint64, req *UpdateTagReq) (*TagResp, error) {
	// 1. 查询是否存在
	t, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrTagNotFound()
	}

	// 2. 领域校验唯一性（排除自身）
	if err = uc.svc.ValidateUpdate(ctx, id, req.Name, req.Slug, req.Color); err != nil {
		return nil, err
	}

	// 3. 自动生成/确保 Slug 唯一（排除自身）
	slug, err := uc.svc.EnsureSlug(ctx, req.Name, req.Slug, id)
	if err != nil {
		return nil, err
	}

	// 4. 规范化颜色，调用领域方法更新
	color := uc.svc.NormalizeColor(req.Color)
	t.Update(req.Name, slug, color)

	// 5. 持久化
	if err = uc.repo.Update(ctx, t); err != nil {
		uc.logger.Error("update tag failed", zap.Uint64("id", id), zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	uc.logger.Info("tag updated",
		zap.Uint64("id", id),
		zap.String("name", t.Name),
		zap.String("color", t.Color),
	)
	return toTagResp(t), nil
}

// DeleteTag 删除标签（软删除，保留历史数据）
func (uc *tagUseCase) DeleteTag(ctx context.Context, id uint64) error {
	// 1. 确认存在
	t, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.ErrTagNotFound()
	}

	// 2. 软删除（article_tags 关联数据通过数据库 CASCADE 处理）
	if err = uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("delete tag failed", zap.Uint64("id", id), zap.Error(err))
		return apperrors.ErrInternalError(err)
	}

	uc.logger.Info("tag deleted", zap.Uint64("id", id), zap.String("name", t.Name))
	return nil
}

// ─── 映射函数 ────────────────────────────────────────────────────────────────

func toTagResp(t *domain.Tag) *TagResp {
	return &TagResp{
		ID:           t.ID,
		Name:         t.Name,
		Slug:         t.Slug,
		Color:        t.Color,
		ArticleCount: t.ArticleCount,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

func toTagSimpleResp(t *domain.Tag) *TagSimpleResp {
	return &TagSimpleResp{
		ID:    t.ID,
		Name:  t.Name,
		Slug:  t.Slug,
		Color: t.Color,
	}
}

func toTagRespList(list []*domain.Tag) []*TagResp {
	resp := make([]*TagResp, 0, len(list))
	for _, t := range list {
		resp = append(resp, toTagResp(t))
	}
	return resp
}

func toTagSimpleRespList(list []*domain.Tag) []*TagSimpleResp {
	resp := make([]*TagSimpleResp, 0, len(list))
	for _, t := range list {
		resp = append(resp, toTagSimpleResp(t))
	}
	return resp
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 标签应用模块
var Module = fx.Options(
	fx.Provide(
		NewTagUseCase,
		domain.NewDomainService,
	),
)
