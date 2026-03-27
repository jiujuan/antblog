// Package comment 评论领域层 —— 纯业务实体，零外部依赖。
package comment

import "time"

// ─── 枚举 ────────────────────────────────────────────────────────────────────

// Status 评论审核状态
type Status int8

const (
	StatusPending  Status = 1 // 待审核
	StatusApproved Status = 2 // 已通过
	StatusRejected Status = 3 // 已拒绝
	StatusSpam     Status = 4 // 垃圾评论
)

func (s Status) IsValid() bool {
	return s >= StatusPending && s <= StatusSpam
}

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusApproved:
		return "approved"
	case StatusRejected:
		return "rejected"
	case StatusSpam:
		return "spam"
	default:
		return "unknown"
	}
}

// IsVisible 是否对前台用户可见
func (s Status) IsVisible() bool { return s == StatusApproved }

// ─── 聚合根 ──────────────────────────────────────────────────────────────────

// Comment 评论聚合根（支持二级嵌套）
//
// 层级规则：
//   - 顶级评论：ParentID == nil, RootID == nil
//   - 子评论：  ParentID != nil（指向顶级评论），RootID == ParentID
//   - 楼层内回复：ParentID != nil, RootID != nil（RootID 指向顶级，ParentID 指向被回复评论）
type Comment struct {
	ID         uint64
	ArticleID  uint64
	UserID     *uint64    // nil = 游客
	ParentID   *uint64    // nil = 顶级评论
	RootID     *uint64    // nil = 顶级评论；否则指向楼层根评论
	ReplyToID  *uint64    // 楼层内回复目标（区分直接回复和楼层内回复）
	Nickname   string     // 游客昵称；登录用户置空（展示时从 User 信息取）
	Email      string     // 游客邮箱（Gravatar 头像）；登录用户置空
	Content    string     // 评论正文（纯文本）
	IP         string     // 评论者 IP（反垃圾）
	UserAgent  string     // 客户端 UA
	Status     Status
	LikeCount  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ─── 业务方法 ────────────────────────────────────────────────────────────────

// Approve 通过审核
func (c *Comment) Approve() {
	c.Status = StatusApproved
	c.UpdatedAt = time.Now()
}

// Reject 拒绝
func (c *Comment) Reject() {
	c.Status = StatusRejected
	c.UpdatedAt = time.Now()
}

// MarkSpam 标记为垃圾评论
func (c *Comment) MarkSpam() {
	c.Status = StatusSpam
	c.UpdatedAt = time.Now()
}

// IsTopLevel 是否为顶级评论
func (c *Comment) IsTopLevel() bool { return c.ParentID == nil }

// IsGuest 是否为游客评论
func (c *Comment) IsGuest() bool { return c.UserID == nil }

// IsVisible 是否对前台可见
func (c *Comment) IsVisible() bool { return c.Status.IsVisible() }

// SetRootID 根据父评论设置 RootID（构建时调用）
func (c *Comment) SetRootID(parent *Comment) {
	if parent == nil {
		return
	}
	if parent.IsTopLevel() {
		// 父评论是顶级，根就是父评论
		c.RootID = &parent.ID
	} else {
		// 父评论已有 root，继承
		c.RootID = parent.RootID
	}
}
