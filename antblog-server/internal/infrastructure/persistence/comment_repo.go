package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	domain "antblog/internal/domain/comment"
	"antblog/internal/infrastructure/persistence/model"
	apperrors "antblog/pkg/errors"
)

// commentRepository ICommentRepository 的 GORM 实现
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db *gorm.DB) domain.ICommentRepository {
	return &commentRepository{db: db}
}

// ─── 单条操作 ────────────────────────────────────────────────────────────────

func (r *commentRepository) Create(ctx context.Context, c *domain.Comment) (*domain.Comment, error) {
	m := commentDomainToModel(c)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return commentModelToDomain(m), nil
}

func (r *commentRepository) FindByID(ctx context.Context, id uint64) (*domain.Comment, error) {
	var m model.Comment
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCommentNotFound()
		}
		return nil, err
	}
	return commentModelToDomain(&m), nil
}

func (r *commentRepository) Update(ctx context.Context, c *domain.Comment) error {
	return r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", c.ID).
		Updates(map[string]any{
			"status": int8(c.Status),
		}).Error
}

func (r *commentRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, "id = ?", id).Error
}

// ─── 前台列表 ────────────────────────────────────────────────────────────────

// FindTopLevelByArticle 查询文章的顶级已通过评论（分页，按时间升序）
func (r *commentRepository) FindTopLevelByArticle(ctx context.Context, articleID uint64, page, pageSize int) ([]*domain.Comment, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("article_id = ? AND parent_id IS NULL AND status = ?",
			articleID, int8(domain.StatusApproved))

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	var list []model.Comment
	offset := (page - 1) * pageSize
	if err := q.Order("created_at ASC").
		Offset(offset).Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return commentModelsToDomainsSlice(list), total, nil
}

// FindChildrenByRoot 查询某根评论下的所有子评论（已通过，按时间升序）
func (r *commentRepository) FindChildrenByRoot(ctx context.Context, rootID uint64) ([]*domain.Comment, error) {
	var list []model.Comment
	err := r.db.WithContext(ctx).
		Where("root_id = ? AND status = ?", rootID, int8(domain.StatusApproved)).
		Order("created_at ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return commentModelsToDomainsSlice(list), nil
}

// ─── 后台管理 ────────────────────────────────────────────────────────────────

// AdminFind 后台评论列表（多条件过滤，分页，按创建时间倒序）
func (r *commentRepository) AdminFind(ctx context.Context, filter *domain.AdminFilter) ([]*domain.Comment, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Comment{})

	if filter.ArticleID != nil {
		q = q.Where("article_id = ?", *filter.ArticleID)
	}
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", int8(*filter.Status))
	}
	if filter.Keyword != "" {
		q = q.Where("content LIKE ?", "%"+filter.Keyword+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	var list []model.Comment
	offset := (filter.Page - 1) * filter.PageSize
	if err := q.Order("created_at DESC").
		Offset(offset).Limit(filter.PageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return commentModelsToDomainsSlice(list), total, nil
}

// ─── 计数 ────────────────────────────────────────────────────────────────────

func (r *commentRepository) CountByArticle(ctx context.Context, articleID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("article_id = ? AND status = ?", articleID, int8(domain.StatusApproved)).
		Count(&count).Error
	return count, err
}

func (r *commentRepository) IncrLikeCount(ctx context.Context, id uint64, delta int) error {
	if delta < 0 {
		return r.db.WithContext(ctx).Model(&model.Comment{}).
			Where("id = ? AND like_count >= ?", id, -delta).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
	}
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

// ─── 模型映射 ────────────────────────────────────────────────────────────────

func commentDomainToModel(c *domain.Comment) *model.Comment {
	return &model.Comment{
		ID:        c.ID,
		ArticleID: c.ArticleID,
		UserID:    c.UserID,
		ParentID:  c.ParentID,
		RootID:    c.RootID,
		ReplyToID: c.ReplyToID,
		Nickname:  c.Nickname,
		Email:     c.Email,
		Content:   c.Content,
		IP:        c.IP,
		UserAgent: c.UserAgent,
		Status:    int8(c.Status),
		LikeCount: c.LikeCount,
	}
}

func commentModelToDomain(m *model.Comment) *domain.Comment {
	return &domain.Comment{
		ID:        m.ID,
		ArticleID: m.ArticleID,
		UserID:    m.UserID,
		ParentID:  m.ParentID,
		RootID:    m.RootID,
		ReplyToID: m.ReplyToID,
		Nickname:  m.Nickname,
		Email:     m.Email,
		Content:   m.Content,
		IP:        m.IP,
		UserAgent: m.UserAgent,
		Status:    domain.Status(m.Status),
		LikeCount: m.LikeCount,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func commentModelsToDomainsSlice(list []model.Comment) []*domain.Comment {
	result := make([]*domain.Comment, 0, len(list))
	for i := range list {
		result = append(result, commentModelToDomain(&list[i]))
	}
	return result
}
