package article

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	apptag "antblog/internal/application/tag"
	domain "antblog/internal/domain/article"
	domaincategory "antblog/internal/domain/category"
	domaintag "antblog/internal/domain/tag"
	apperrors "antblog/pkg/errors"
	"antblog/pkg/utils"
)

// ─── 依赖声明 ────────────────────────────────────────────────────────────────

// Deps fx 注入依赖
type Deps struct {
	fx.In
	Repo          domain.IArticleRepository
	CategoryRepo  domaincategory.ICategoryRepository
	TagRepo       domaintag.ITagRepository
	DomainService domain.IDomainService
	Logger        *zap.Logger
}

type articleUseCase struct {
	repo         domain.IArticleRepository
	categoryRepo domaincategory.ICategoryRepository
	tagRepo      domaintag.ITagRepository
	svc          domain.IDomainService
	logger       *zap.Logger
}

// NewArticleUseCase 创建文章用例（fx provider）
func NewArticleUseCase(deps Deps) IArticleUseCase {
	return &articleUseCase{
		repo:         deps.Repo,
		categoryRepo: deps.CategoryRepo,
		tagRepo:      deps.TagRepo,
		svc:          deps.DomainService,
		logger:       deps.Logger,
	}
}

// ─── 前台：列表 / 详情 ────────────────────────────────────────────────────────

func (uc *articleUseCase) ListArticles(ctx context.Context, req *ListArticleReq, userID *uint64) ([]*ArticleListItemResp, int64, error) {
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	status := domain.StatusPublished
	filter := &domain.ListFilter{
		Status:   &status,
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}
	if req.CategoryID > 0 {
		filter.CategoryID = &req.CategoryID
	}
	if req.TagID > 0 {
		filter.TagID = &req.TagID
	}

	list, total, err := uc.repo.FindList(ctx, filter)
	if err != nil {
		uc.logger.Error("list articles failed", zap.Error(err))
		return nil, 0, apperrors.ErrInternalError(err)
	}
	return uc.toListItemRespBatch(ctx, list, userID), total, nil
}

func (uc *articleUseCase) GetArticleBySlug(ctx context.Context, slug string, userID *uint64) (*ArticleResp, error) {
	a, err := uc.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrArticleNotFound()
	}
	if !a.IsVisibleToPublic() {
		return nil, apperrors.ErrArticleNotFound()
	}

	// 异步递增阅读数（不阻塞主流程）
	go func() {
		if e := uc.repo.IncrViewCount(context.Background(), a.ID); e != nil {
			uc.logger.Warn("incr view count failed", zap.Uint64("id", a.ID), zap.Error(e))
		}
	}()

	resp := uc.toArticleResp(ctx, a)

	// 注入互动状态
	if userID != nil {
		liked, _ := uc.repo.HasLiked(ctx, a.ID, *userID)
		bookmarked, _ := uc.repo.HasBookmarked(ctx, a.ID, *userID)
		resp.Liked = &liked
		resp.Bookmarked = &bookmarked
	}

	return resp, nil
}

func (uc *articleUseCase) GetFeaturedArticles(ctx context.Context, limit int, userID *uint64) ([]*ArticleListItemResp, error) {
	if limit <= 0 || limit > 20 {
		limit = 6
	}
	list, err := uc.repo.FindFeatured(ctx, limit)
	if err != nil {
		return nil, apperrors.ErrInternalError(err)
	}
	return uc.toListItemRespBatch(ctx, list, userID), nil
}

func (uc *articleUseCase) GetArchive(ctx context.Context) ([]*ArchiveItemResp, error) {
	items, err := uc.repo.GetArchive(ctx)
	if err != nil {
		return nil, apperrors.ErrInternalError(err)
	}
	resp := make([]*ArchiveItemResp, 0, len(items))
	for _, item := range items {
		resp = append(resp, &ArchiveItemResp{
			Year:         item.Year,
			Month:        item.Month,
			ArticleCount: item.ArticleCount,
		})
	}
	return resp, nil
}

func (uc *articleUseCase) GetArchiveDetail(ctx context.Context, req *ArchiveDetailReq, userID *uint64) ([]*ArticleListItemResp, int64, error) {
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	list, total, err := uc.repo.FindByYearMonth(ctx, req.Year, req.Month, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, apperrors.ErrInternalError(err)
	}
	return uc.toListItemRespBatch(ctx, list, userID), total, nil
}

// ─── 前台：互动 ──────────────────────────────────────────────────────────────

func (uc *articleUseCase) LikeArticle(ctx context.Context, articleID, userID uint64) error {
	if _, err := uc.repo.FindByID(ctx, articleID); err != nil {
		return apperrors.ErrArticleNotFound()
	}
	if err := uc.repo.AddLike(ctx, articleID, userID); err != nil {
		return apperrors.ErrInternalError(err)
	}
	_ = uc.repo.IncrLikeCount(ctx, articleID, 1)
	return nil
}

func (uc *articleUseCase) UnlikeArticle(ctx context.Context, articleID, userID uint64) error {
	if err := uc.repo.RemoveLike(ctx, articleID, userID); err != nil {
		return apperrors.ErrInternalError(err)
	}
	_ = uc.repo.IncrLikeCount(ctx, articleID, -1)
	return nil
}

func (uc *articleUseCase) BookmarkArticle(ctx context.Context, articleID, userID uint64) error {
	if _, err := uc.repo.FindByID(ctx, articleID); err != nil {
		return apperrors.ErrArticleNotFound()
	}
	if err := uc.repo.AddBookmark(ctx, articleID, userID); err != nil {
		return apperrors.ErrInternalError(err)
	}
	_ = uc.repo.IncrBookmarkCount(ctx, articleID, 1)
	return nil
}

func (uc *articleUseCase) UnbookmarkArticle(ctx context.Context, articleID, userID uint64) error {
	if err := uc.repo.RemoveBookmark(ctx, articleID, userID); err != nil {
		return apperrors.ErrInternalError(err)
	}
	_ = uc.repo.IncrBookmarkCount(ctx, articleID, -1)
	return nil
}

func (uc *articleUseCase) GetUserBookmarks(ctx context.Context, userID uint64, page, pageSize int) ([]*ArticleListItemResp, int64, error) {
	page = utils.NormalizePage(page)
	pageSize = utils.NormalizePageSize(pageSize)
	list, total, err := uc.repo.GetUserBookmarks(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, apperrors.ErrInternalError(err)
	}
	// 收藏列表里所有文章对该用户一定是已收藏状态，直接注入
	uid := userID
	return uc.toListItemRespBatch(ctx, list, &uid), total, nil
}

// ─── 后台管理 ────────────────────────────────────────────────────────────────

func (uc *articleUseCase) AdminListArticles(ctx context.Context, req *AdminListArticleReq) ([]*ArticleListItemResp, int64, error) {
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	filter := &domain.ListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}
	if req.Status > 0 {
		s := domain.Status(req.Status)
		filter.Status = &s
	}
	if req.CategoryID > 0 {
		filter.CategoryID = &req.CategoryID
	}
	if req.TagID > 0 {
		filter.TagID = &req.TagID
	}

	list, total, err := uc.repo.FindList(ctx, filter)
	if err != nil {
		uc.logger.Error("admin list articles failed", zap.Error(err))
		return nil, 0, apperrors.ErrInternalError(err)
	}
	return uc.toListItemRespBatch(ctx, list, nil), total, nil
}

func (uc *articleUseCase) AdminGetArticle(ctx context.Context, id uint64) (*ArticleResp, error) {
	a, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrArticleNotFound()
	}
	return uc.toArticleResp(ctx, a), nil
}

func (uc *articleUseCase) CreateArticle(ctx context.Context, authorID uint64, req *CreateArticleReq) (*ArticleResp, error) {
	// 1. 领域校验
	if err := uc.svc.ValidateCreate(ctx, req.Title, req.Slug); err != nil {
		return nil, err
	}

	// 2. 确保 Slug 唯一
	slug, err := uc.svc.EnsureSlug(ctx, req.Title, req.Slug, 0)
	if err != nil {
		return nil, err
	}

	// 3. 构建实体
	a, err := uc.svc.BuildArticle(authorID, domain.BuildArticleReq{
		Title:        req.Title,
		Slug:         slug,
		Summary:      req.Summary,
		Content:      req.Content,
		ContentHTML:  req.Content, // 前端传入 HTML，或服务端后续渲染
		Cover:        req.Cover,
		CategoryID:   req.CategoryID,
		TagIDs:       req.TagIDs,
		Status:       domain.Status(req.Status),
		IsTop:        req.IsTop,
		IsFeatured:   req.IsFeatured,
		AllowComment: req.AllowComment,
	})
	if err != nil {
		return nil, err
	}

	// 4. 持久化（含标签关联）
	created, err := uc.repo.Create(ctx, a)
	if err != nil {
		uc.logger.Error("create article failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	// 5. 更新标签文章计数
	go uc.incrTagCounts(created.TagIDs, 1)
	uc.syncCategoryCountOnCreate(created)

	uc.logger.Info("article created",
		zap.Uint64("id", created.ID),
		zap.String("title", created.Title),
		zap.String("status", created.Status.String()),
	)
	return uc.toArticleResp(ctx, created), nil
}

func (uc *articleUseCase) UpdateArticle(ctx context.Context, id uint64, req *UpdateArticleReq) (*ArticleResp, error) {
	// 1. 查询存在
	a, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrArticleNotFound()
	}
	oldTagIDs := a.TagIDs
	oldCategoryID := a.CategoryID

	// 2. 领域校验
	if err = uc.svc.ValidateUpdate(ctx, id, req.Title, req.Slug); err != nil {
		return nil, err
	}

	// 3. 确保 Slug 唯一（排除自身）
	slug, err := uc.svc.EnsureSlug(ctx, req.Title, req.Slug, id)
	if err != nil {
		return nil, err
	}

	// 4. 调用领域方法更新
	wordCount := domain.CountWords(req.Content)
	summary := domain.AutoSummary(req.Content, req.Summary)
	a.UpdateContent(
		req.Title, slug, summary, req.Content, req.Content,
		req.Cover, req.CategoryID, req.TagIDs,
		req.AllowComment, req.IsTop, req.IsFeatured,
	)
	a.WordCount = wordCount

	// 5. 持久化
	if err = uc.repo.Update(ctx, a); err != nil {
		uc.logger.Error("update article failed", zap.Uint64("id", id), zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	// 6. 同步标签计数变化（异步）
	go uc.syncTagCounts(oldTagIDs, req.TagIDs)
	uc.syncCategoryCountOnCategoryChange(a.Status, oldCategoryID, req.CategoryID)

	uc.logger.Info("article updated", zap.Uint64("id", id), zap.String("title", a.Title))
	return uc.toArticleResp(ctx, a), nil
}

func (uc *articleUseCase) UpdateArticleStatus(ctx context.Context, id uint64, req *UpdateStatusReq) (*ArticleResp, error) {
	a, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrArticleNotFound()
	}

	newStatus := domain.Status(req.Status)
	oldStatus := a.Status
	switch newStatus {
	case domain.StatusPublished:
		a.Publish()
	case domain.StatusArchived:
		a.Archive()
	case domain.StatusDraft:
		a.Unpublish()
	default:
		return nil, apperrors.ErrInvalidParams("无效的文章状态")
	}

	if err = uc.repo.UpdateStatus(ctx, id, a.Status, a.PublishedAt); err != nil {
		return nil, apperrors.ErrInternalError(err)
	}
	uc.syncCategoryCountOnStatusChange(oldStatus, a.Status, a.CategoryID)

	uc.logger.Info("article status updated",
		zap.Uint64("id", id),
		zap.String("status", newStatus.String()),
	)
	return uc.toArticleResp(ctx, a), nil
}

func (uc *articleUseCase) DeleteArticle(ctx context.Context, id uint64) error {
	a, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.ErrArticleNotFound()
	}
	if err = uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("delete article failed", zap.Uint64("id", id), zap.Error(err))
		return apperrors.ErrInternalError(err)
	}
	// 删除后更新标签计数
	go uc.incrTagCounts(a.TagIDs, -1)
	uc.syncCategoryCountOnDelete(a)
	uc.logger.Info("article deleted", zap.Uint64("id", id))
	return nil
}

// ─── 映射函数 ────────────────────────────────────────────────────────────────

func (uc *articleUseCase) toArticleResp(ctx context.Context, a *domain.Article) *ArticleResp {
	tags := uc.fetchTagSimpleResps(ctx, a.TagIDs)
	return &ArticleResp{
		ID:            a.ID,
		UUID:          a.UUID,
		AuthorID:      a.AuthorID,
		CategoryID:    a.CategoryID,
		Title:         a.Title,
		Slug:          a.Slug,
		Summary:       a.Summary,
		Content:       a.Content,
		ContentHTML:   a.ContentHTML,
		Cover:         a.Cover,
		Status:        int8(a.Status),
		StatusText:    a.Status.String(),
		IsTop:         a.IsTop,
		IsFeatured:    a.IsFeatured,
		AllowComment:  a.AllowComment,
		ViewCount:     a.ViewCount,
		LikeCount:     a.LikeCount,
		CommentCount:  a.CommentCount,
		BookmarkCount: a.BookmarkCount,
		WordCount:     a.WordCount,
		PublishedAt:   a.PublishedAt,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
		Tags:          tags,
	}
}

func (uc *articleUseCase) toListItemResp(ctx context.Context, a *domain.Article) *ArticleListItemResp {
	tags := uc.fetchTagSimpleResps(ctx, a.TagIDs)
	return &ArticleListItemResp{
		ID:            a.ID,
		UUID:          a.UUID,
		AuthorID:      a.AuthorID,
		CategoryID:    a.CategoryID,
		Title:         a.Title,
		Slug:          a.Slug,
		Summary:       a.Summary,
		Cover:         a.Cover,
		Status:        int8(a.Status),
		IsTop:         a.IsTop,
		IsFeatured:    a.IsFeatured,
		ViewCount:     a.ViewCount,
		LikeCount:     a.LikeCount,
		CommentCount:  a.CommentCount,
		BookmarkCount: a.BookmarkCount,
		WordCount:     a.WordCount,
		PublishedAt:   a.PublishedAt,
		CreatedAt:     a.CreatedAt,
		Tags:          tags,
	}
}

// toListItemRespBatch 批量转换列表，若 userID 非 nil 则一次查询所有互动状态（避免 N+1）
func (uc *articleUseCase) toListItemRespBatch(ctx context.Context, list []*domain.Article, userID *uint64) []*ArticleListItemResp {
	if len(list) == 0 {
		return []*ArticleListItemResp{}
	}

	resp := make([]*ArticleListItemResp, 0, len(list))
	for _, a := range list {
		resp = append(resp, uc.toListItemResp(ctx, a))
	}

	// 批量注入互动状态（登录用户才查）
	if userID != nil {
		articleIDs := make([]uint64, 0, len(list))
		for _, a := range list {
			articleIDs = append(articleIDs, a.ID)
		}

		likedMap, err := uc.repo.BatchHasLiked(ctx, articleIDs, *userID)
		if err != nil {
			uc.logger.Warn("batch has liked failed", zap.Error(err))
			likedMap = map[uint64]bool{}
		}

		bookmarkedMap, err := uc.repo.BatchHasBookmarked(ctx, articleIDs, *userID)
		if err != nil {
			uc.logger.Warn("batch has bookmarked failed", zap.Error(err))
			bookmarkedMap = map[uint64]bool{}
		}

		for i, item := range resp {
			liked := likedMap[list[i].ID]
			bookmarked := bookmarkedMap[list[i].ID]
			item.Liked = &liked
			item.Bookmarked = &bookmarked
		}
	}

	return resp
}

// fetchTagSimpleResps 根据 TagIDs 查询标签精简信息
func (uc *articleUseCase) fetchTagSimpleResps(ctx context.Context, tagIDs []uint64) []*apptag.TagSimpleResp {
	if len(tagIDs) == 0 {
		return []*apptag.TagSimpleResp{}
	}
	tags, err := uc.tagRepo.FindByIDs(ctx, tagIDs)
	if err != nil {
		uc.logger.Warn("fetch tags failed", zap.Error(err))
		return []*apptag.TagSimpleResp{}
	}
	resp := make([]*apptag.TagSimpleResp, 0, len(tags))
	for _, t := range tags {
		resp = append(resp, &apptag.TagSimpleResp{
			ID:    t.ID,
			Name:  t.Name,
			Slug:  t.Slug,
			Color: t.Color,
		})
	}
	return resp
}

// incrTagCounts 批量增减标签文章计数
func (uc *articleUseCase) incrTagCounts(tagIDs []uint64, delta int) {
	for _, tid := range tagIDs {
		if err := uc.tagRepo.IncrArticleCount(context.Background(), tid, delta); err != nil {
			uc.logger.Warn("incr tag article_count failed",
				zap.Uint64("tag_id", tid), zap.Int("delta", delta), zap.Error(err))
		}
	}
}

// syncTagCounts 比对新旧标签列表，差量更新计数
func (uc *articleUseCase) syncTagCounts(oldIDs, newIDs []uint64) {
	oldSet := make(map[uint64]struct{}, len(oldIDs))
	for _, id := range oldIDs {
		oldSet[id] = struct{}{}
	}
	newSet := make(map[uint64]struct{}, len(newIDs))
	for _, id := range newIDs {
		newSet[id] = struct{}{}
	}
	// 新增的标签 +1
	for id := range newSet {
		if _, ok := oldSet[id]; !ok {
			_ = uc.tagRepo.IncrArticleCount(context.Background(), id, 1)
		}
	}
	// 移除的标签 -1
	for id := range oldSet {
		if _, ok := newSet[id]; !ok {
			_ = uc.tagRepo.IncrArticleCount(context.Background(), id, -1)
		}
	}
}

func (uc *articleUseCase) syncCategoryCountOnCreate(a *domain.Article) {
	if a == nil || a.CategoryID == nil || a.Status != domain.StatusPublished {
		return
	}
	uc.incrCategoryCount(*a.CategoryID, 1)
}

func (uc *articleUseCase) syncCategoryCountOnDelete(a *domain.Article) {
	if a == nil || a.CategoryID == nil || a.Status != domain.StatusPublished {
		return
	}
	uc.incrCategoryCount(*a.CategoryID, -1)
}

func (uc *articleUseCase) syncCategoryCountOnCategoryChange(status domain.Status, oldCategoryID, newCategoryID *uint64) {
	if status != domain.StatusPublished {
		return
	}
	if categoryIDEq(oldCategoryID, newCategoryID) {
		return
	}
	if oldCategoryID != nil {
		uc.incrCategoryCount(*oldCategoryID, -1)
	}
	if newCategoryID != nil {
		uc.incrCategoryCount(*newCategoryID, 1)
	}
}

func (uc *articleUseCase) syncCategoryCountOnStatusChange(oldStatus, newStatus domain.Status, categoryID *uint64) {
	if categoryID == nil || oldStatus == newStatus {
		return
	}
	if oldStatus != domain.StatusPublished && newStatus == domain.StatusPublished {
		uc.incrCategoryCount(*categoryID, 1)
		return
	}
	if oldStatus == domain.StatusPublished && newStatus != domain.StatusPublished {
		uc.incrCategoryCount(*categoryID, -1)
	}
}

func (uc *articleUseCase) incrCategoryCount(categoryID uint64, delta int) {
	if err := uc.categoryRepo.IncrArticleCount(context.Background(), categoryID, delta); err != nil {
		uc.logger.Warn("incr category article_count failed",
			zap.Uint64("category_id", categoryID), zap.Int("delta", delta), zap.Error(err))
	}
}

func categoryIDEq(a, b *uint64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 文章应用模块
var Module = fx.Options(
	fx.Provide(
		NewArticleUseCase,
		domain.NewDomainService,
	),
)
