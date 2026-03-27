// Package comment 评论应用层。
package comment

import "time"

// ─── 请求 DTO ────────────────────────────────────────────────────────────────

// CreateCommentReq 发表评论请求
type CreateCommentReq struct {
	ArticleID uint64  `json:"article_id" validate:"required"`
	ParentID  *uint64 `json:"parent_id"`   // 为 nil 时发布顶级评论
	ReplyToID *uint64 `json:"reply_to_id"` // 楼层内回复目标
	// 游客信息（登录用户无需填写）
	Nickname string `json:"nickname" validate:"max=64"`
	Email    string `json:"email"    validate:"omitempty,max=128,email"`
	Content  string `json:"content"  validate:"required,max=2000"`
}

// AdminUpdateStatusReq 后台审核请求
type AdminUpdateStatusReq struct {
	Status int8 `json:"status" validate:"required,oneof=1 2 3 4"`
}

// AdminListCommentReq 后台评论列表查询请求
type AdminListCommentReq struct {
	Page      int    `form:"page"       validate:"min=1"`
	PageSize  int    `form:"page_size"  validate:"min=1,max=100"`
	ArticleID uint64 `form:"article_id"`
	UserID    uint64 `form:"user_id"`
	Status    int8   `form:"status"     validate:"omitempty,oneof=1 2 3 4"`
	Keyword   string `form:"keyword"    validate:"max=100"`
}

// ListCommentReq 前台评论列表请求
type ListCommentReq struct {
	ArticleID uint64 `form:"article_id" validate:"required"`
	Page      int    `form:"page"       validate:"min=1"`
	PageSize  int    `form:"page_size"  validate:"min=1,max=50"`
}

// ─── 响应 DTO ────────────────────────────────────────────────────────────────

// CommentResp 评论响应（通用）
type CommentResp struct {
	ID         uint64         `json:"id"`
	ArticleID  uint64         `json:"article_id"`
	UserID     *uint64        `json:"user_id"`
	ParentID   *uint64        `json:"parent_id"`
	RootID     *uint64        `json:"root_id"`
	ReplyToID  *uint64        `json:"reply_to_id"`
	Nickname   string         `json:"nickname"`   // 游客昵称或用户昵称
	Avatar     string         `json:"avatar"`     // Gravatar URL 或用户头像
	Content    string         `json:"content"`
	Status     int8           `json:"status,omitempty"` // 前台不暴露
	StatusText string         `json:"status_text,omitempty"`
	LikeCount  int            `json:"like_count"`
	CreatedAt  time.Time      `json:"created_at"`
	Children   []*CommentResp `json:"children,omitempty"` // 子评论（前台树形返回时使用）
}

// CommentTreeResp 前台楼层树形响应（顶级评论 + 子评论列表）
type CommentTreeResp struct {
	List     []*CommentResp `json:"list"`
	Total    int64          `json:"total"`     // 顶级评论总数
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}
