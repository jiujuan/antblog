package comment

import "context"

// ICommentUseCase 评论用例接口
type ICommentUseCase interface {
	// ── 前台接口 ─────────────────────────────────────────────────────────────

	// ListComments 获取文章评论树（顶级 + 子评论，仅已通过，分页）
	ListComments(ctx context.Context, req *ListCommentReq) (*CommentTreeResp, error)

	// CreateComment 发表评论（支持登录用户和游客）
	CreateComment(ctx context.Context, req *CreateCommentReq, userID *uint64, ip, ua string) (*CommentResp, error)

	// LikeComment 点赞评论（幂等，不强制登录）
	LikeComment(ctx context.Context, id uint64) error

	// ── 后台管理接口 ──────────────────────────────────────────────────────────

	// AdminListComments 后台评论列表（全状态，支持多条件过滤）
	AdminListComments(ctx context.Context, req *AdminListCommentReq) ([]*CommentResp, int64, error)

	// AdminGetComment 后台获取评论详情
	AdminGetComment(ctx context.Context, id uint64) (*CommentResp, error)

	// AdminUpdateStatus 审核评论（通过/拒绝/标记垃圾）
	AdminUpdateStatus(ctx context.Context, id uint64, req *AdminUpdateStatusReq) (*CommentResp, error)

	// AdminDeleteComment 后台删除评论
	AdminDeleteComment(ctx context.Context, id uint64) error
}
