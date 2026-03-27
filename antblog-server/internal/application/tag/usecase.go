package tag

import "context"

// ITagUseCase 标签用例接口
type ITagUseCase interface {
	// ── 前台接口 ─────────────────────────────────────────────────────────────

	// ListTags 获取所有标签（按文章数降序）
	ListTags(ctx context.Context) ([]*TagResp, error)

	// GetTagBySlug 按 Slug 获取标签详情
	GetTagBySlug(ctx context.Context, slug string) (*TagResp, error)

	// GetTagsByArticleID 获取某篇文章的所有标签（精简版，供文章详情使用）
	GetTagsByArticleID(ctx context.Context, articleID uint64) ([]*TagSimpleResp, error)

	// ── 后台管理接口 ──────────────────────────────────────────────────────────

	// GetTagByID 按 ID 获取标签（后台）
	GetTagByID(ctx context.Context, id uint64) (*TagResp, error)

	// CreateTag 创建单个标签
	CreateTag(ctx context.Context, req *CreateTagReq) (*TagResp, error)

	// BatchCreateTags 批量创建标签（幂等：已存在的名称直接返回，不报错）
	BatchCreateTags(ctx context.Context, req *BatchCreateTagReq) ([]*TagResp, error)

	// UpdateTag 更新标签
	UpdateTag(ctx context.Context, id uint64, req *UpdateTagReq) (*TagResp, error)

	// DeleteTag 删除标签
	DeleteTag(ctx context.Context, id uint64) error
}
