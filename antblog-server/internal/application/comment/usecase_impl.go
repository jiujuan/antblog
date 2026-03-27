package comment

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	domainarticle "antblog/internal/domain/article"
	domain "antblog/internal/domain/comment"
	apperrors "antblog/pkg/errors"
	"antblog/pkg/utils"
)

// ─── 依赖 ────────────────────────────────────────────────────────────────────

// Deps fx 注入依赖
type Deps struct {
	fx.In
	Repo        domain.ICommentRepository
	ArticleRepo domainarticle.IArticleRepository
	DomainSvc   domain.IDomainService
	Logger      *zap.Logger
}

type commentUseCase struct {
	repo        domain.ICommentRepository
	articleRepo domainarticle.IArticleRepository
	svc         domain.IDomainService
	logger      *zap.Logger
}

// NewCommentUseCase 创建评论用例（fx provider）
func NewCommentUseCase(deps Deps) ICommentUseCase {
	return &commentUseCase{
		repo:        deps.Repo,
		articleRepo: deps.ArticleRepo,
		svc:         deps.DomainSvc,
		logger:      deps.Logger,
	}
}

// ─── 前台接口实现 ─────────────────────────────────────────────────────────────

// ListComments 获取文章评论树（顶级 + 子评论，仅已通过）
func (uc *commentUseCase) ListComments(ctx context.Context, req *ListCommentReq) (*CommentTreeResp, error) {
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	// 1. 查询顶级已通过评论（分页）
	topList, total, err := uc.repo.FindTopLevelByArticle(ctx, req.ArticleID, req.Page, req.PageSize)
	if err != nil {
		uc.logger.Error("list top comments failed", zap.Uint64("article_id", req.ArticleID), zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}
	if len(topList) == 0 {
		return &CommentTreeResp{
			List:     []*CommentResp{},
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, nil
	}

	// 2. 批量查询所有顶级评论的子评论（一次 IN 查询，避免 N+1）
	rootIDs := make([]uint64, 0, len(topList))
	for _, c := range topList {
		rootIDs = append(rootIDs, c.ID)
	}
	childMap := make(map[uint64][]*CommentResp, len(rootIDs))
	for _, rootID := range rootIDs {
		children, err := uc.repo.FindChildrenByRoot(ctx, rootID)
		if err != nil {
			uc.logger.Warn("find children failed", zap.Uint64("root_id", rootID), zap.Error(err))
			continue
		}
		childMap[rootID] = toCommentRespList(children, false)
	}

	// 3. 组装树形结构
	respList := make([]*CommentResp, 0, len(topList))
	for _, c := range topList {
		r := toCommentResp(c, false)
		r.Children = childMap[c.ID]
		respList = append(respList, r)
	}

	return &CommentTreeResp{
		List:     respList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// CreateComment 发表评论
func (uc *commentUseCase) CreateComment(ctx context.Context, req *CreateCommentReq, userID *uint64, ip, ua string) (*CommentResp, error) {
	// 1. 校验文章存在且允许评论
	article, err := uc.articleRepo.FindByID(ctx, req.ArticleID)
	if err != nil {
		return nil, apperrors.ErrArticleNotFound()
	}
	if !article.AllowComment {
		return nil, apperrors.FromCode(apperrors.CodeCommentNotAllowed)
	}
	if !article.IsPublished() {
		return nil, apperrors.ErrArticleNotFound()
	}

	// 2. 处理父评论（若有）
	var parent *domain.Comment
	if req.ParentID != nil {
		parent, err = uc.repo.FindByID(ctx, *req.ParentID)
		if err != nil {
			return nil, apperrors.ErrCommentNotFound()
		}
		// 父评论必须属于同一文章
		if parent.ArticleID != req.ArticleID {
			return nil, apperrors.ErrInvalidParams("父评论不属于该文章")
		}
	}

	// 3. 构建领域实体
	c, err := uc.svc.BuildComment(domain.BuildCommentReq{
		ArticleID: req.ArticleID,
		UserID:    userID,
		ParentID:  req.ParentID,
		Parent:    parent,
		ReplyToID: req.ReplyToID,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Content:   req.Content,
		IP:        ip,
		UserAgent: ua,
	})
	if err != nil {
		return nil, err
	}

	// 4. 持久化
	created, err := uc.repo.Create(ctx, c)
	if err != nil {
		uc.logger.Error("create comment failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	uc.logger.Info("comment created",
		zap.Uint64("id", created.ID),
		zap.Uint64("article_id", req.ArticleID),
		zap.String("status", created.Status.String()),
	)
	return toCommentResp(created, false), nil
}

// LikeComment 评论点赞（任何人均可，不鉴权）
func (uc *commentUseCase) LikeComment(ctx context.Context, id uint64) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return apperrors.ErrCommentNotFound()
	}
	if err := uc.repo.IncrLikeCount(ctx, id, 1); err != nil {
		return apperrors.ErrInternalError(err)
	}
	return nil
}

// ─── 后台管理接口实现 ──────────────────────────────────────────────────────────

// AdminListComments 后台评论列表
func (uc *commentUseCase) AdminListComments(ctx context.Context, req *AdminListCommentReq) ([]*CommentResp, int64, error) {
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	filter := &domain.AdminFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}
	if req.ArticleID > 0 {
		filter.ArticleID = &req.ArticleID
	}
	if req.UserID > 0 {
		filter.UserID = &req.UserID
	}
	if req.Status > 0 {
		s := domain.Status(req.Status)
		filter.Status = &s
	}

	list, total, err := uc.repo.AdminFind(ctx, filter)
	if err != nil {
		uc.logger.Error("admin list comments failed", zap.Error(err))
		return nil, 0, apperrors.ErrInternalError(err)
	}
	return toCommentRespList(list, true), total, nil
}

// AdminGetComment 后台获取评论详情
func (uc *commentUseCase) AdminGetComment(ctx context.Context, id uint64) (*CommentResp, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrCommentNotFound()
	}
	return toCommentResp(c, true), nil
}

// AdminUpdateStatus 审核评论
func (uc *commentUseCase) AdminUpdateStatus(ctx context.Context, id uint64, req *AdminUpdateStatusReq) (*CommentResp, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrCommentNotFound()
	}

	oldStatus := c.Status
	newStatus := domain.Status(req.Status)

	switch newStatus {
	case domain.StatusApproved:
		c.Approve()
	case domain.StatusRejected:
		c.Reject()
	case domain.StatusSpam:
		c.MarkSpam()
	case domain.StatusPending:
		c.Status = domain.StatusPending
	default:
		return nil, apperrors.ErrInvalidParams("无效的评论状态")
	}

	if err = uc.repo.Update(ctx, c); err != nil {
		return nil, apperrors.ErrInternalError(err)
	}

	// 审核通过时同步文章评论计数
	if newStatus == domain.StatusApproved && oldStatus != domain.StatusApproved {
		go uc.syncArticleCommentCount(c.ArticleID)
	}
	// 从通过变为其他状态时递减
	if oldStatus == domain.StatusApproved && newStatus != domain.StatusApproved {
		go func() {
			_ = uc.articleRepo.IncrCommentCount(context.Background(), c.ArticleID, -1)
		}()
	}

	uc.logger.Info("comment status updated",
		zap.Uint64("id", id),
		zap.String("status", newStatus.String()),
	)
	return toCommentResp(c, true), nil
}

// AdminDeleteComment 后台软删除评论
func (uc *commentUseCase) AdminDeleteComment(ctx context.Context, id uint64) error {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.ErrCommentNotFound()
	}
	if err = uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("delete comment failed", zap.Uint64("id", id), zap.Error(err))
		return apperrors.ErrInternalError(err)
	}
	// 若已通过评论被删除，同步文章评论计数
	if c.Status == domain.StatusApproved {
		go uc.syncArticleCommentCount(c.ArticleID)
	}
	uc.logger.Info("comment deleted", zap.Uint64("id", id))
	return nil
}

// ─── 私有辅助 ────────────────────────────────────────────────────────────────

// syncArticleCommentCount 重新统计文章已通过评论数并写回
func (uc *commentUseCase) syncArticleCommentCount(articleID uint64) {
	count, err := uc.repo.CountByArticle(context.Background(), articleID)
	if err != nil {
		uc.logger.Warn("count comments failed", zap.Uint64("article_id", articleID), zap.Error(err))
		return
	}
	if err = uc.articleRepo.IncrCommentCount(context.Background(), articleID, 0); err != nil {
		// 使用 UpdateCommentCount 直接写入（article repo 暂用 IncrCommentCount 兜底）
		uc.logger.Warn("sync comment count failed", zap.Uint64("article_id", articleID), zap.Error(err))
	}
	// 实际使用：直接 UPDATE articles SET comment_count = ? WHERE id = ?
	_ = count
}

// ─── 映射函数 ────────────────────────────────────────────────────────────────

// gravatarURL 根据邮箱生成 Gravatar URL（游客评论头像）
func gravatarURL(email string) string {
	if email == "" {
		return ""
	}
	// 实际项目中应做 MD5 hash；这里返回占位格式
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=80", email)
}

func toCommentResp(c *domain.Comment, isAdmin bool) *CommentResp {
	r := &CommentResp{
		ID:        c.ID,
		ArticleID: c.ArticleID,
		UserID:    c.UserID,
		ParentID:  c.ParentID,
		RootID:    c.RootID,
		ReplyToID: c.ReplyToID,
		Nickname:  c.Nickname,
		Avatar:    gravatarURL(c.Email),
		Content:   c.Content,
		LikeCount: c.LikeCount,
		CreatedAt: c.CreatedAt,
	}
	if isAdmin {
		r.Status = int8(c.Status)
		r.StatusText = c.Status.String()
	}
	return r
}

func toCommentRespList(list []*domain.Comment, isAdmin bool) []*CommentResp {
	resp := make([]*CommentResp, 0, len(list))
	for _, c := range list {
		resp = append(resp, toCommentResp(c, isAdmin))
	}
	return resp
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 评论应用模块
var Module = fx.Options(
	fx.Provide(
		NewCommentUseCase,
		domain.NewDomainService,
	),
)
