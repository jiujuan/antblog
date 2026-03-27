package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	domain "antblog/internal/domain/tag"
	"antblog/internal/infrastructure/persistence/model"
	apperrors "antblog/pkg/errors"
)

// tagRepository ITagRepository 的 GORM 实现
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db *gorm.DB) domain.ITagRepository {
	return &tagRepository{db: db}
}

// ─── 单条操作 ────────────────────────────────────────────────────────────────

// Create 创建标签
func (r *tagRepository) Create(ctx context.Context, t *domain.Tag) (*domain.Tag, error) {
	m := tagDomainToModel(t)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return tagModelToDomain(m), nil
}

// FindByID 按 ID 查询标签
func (r *tagRepository) FindByID(ctx context.Context, id uint64) (*domain.Tag, error) {
	var m model.Tag
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrTagNotFound()
		}
		return nil, err
	}
	return tagModelToDomain(&m), nil
}

// FindBySlug 按 Slug 查询标签
func (r *tagRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	var m model.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrTagNotFound()
		}
		return nil, err
	}
	return tagModelToDomain(&m), nil
}

// FindByName 按名称精确查询
func (r *tagRepository) FindByName(ctx context.Context, name string) (*domain.Tag, error) {
	var m model.Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrTagNotFound()
		}
		return nil, err
	}
	return tagModelToDomain(&m), nil
}

// Update 更新标签
func (r *tagRepository) Update(ctx context.Context, t *domain.Tag) error {
	return r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Where("id = ?", t.ID).
		Updates(map[string]any{
			"name":  t.Name,
			"slug":  t.Slug,
			"color": t.Color,
		}).Error
}

// Delete 软删除标签
func (r *tagRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Tag{}, "id = ?", id).Error
}

// ─── 批量查询 ────────────────────────────────────────────────────────────────

// FindAll 查询所有标签，按文章数降序，article_count 相同时按 ID 升序
func (r *tagRepository) FindAll(ctx context.Context) ([]*domain.Tag, error) {
	var list []model.Tag
	if err := r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Select("tags.*, COUNT(articles.id) AS article_count").
		Joins("LEFT JOIN article_tags ON article_tags.tag_id = tags.id").
		Joins("LEFT JOIN articles ON articles.id = article_tags.article_id AND articles.status = ? AND articles.deleted_at IS NULL", 2).
		Group("tags.id").
		Order("article_count DESC, id ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.Tag, 0, len(list))
	for i := range list {
		result = append(result, tagModelToDomain(&list[i]))
	}
	return result, nil
}

// FindByIDs 按 ID 批量查询（保持传入顺序）
func (r *tagRepository) FindByIDs(ctx context.Context, ids []uint64) ([]*domain.Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []model.Tag
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&list).Error; err != nil {
		return nil, err
	}

	// 按传入 ids 顺序返回（构建 map 后按顺序取）
	tagMap := make(map[uint64]*domain.Tag, len(list))
	for i := range list {
		d := tagModelToDomain(&list[i])
		tagMap[d.ID] = d
	}
	result := make([]*domain.Tag, 0, len(ids))
	for _, id := range ids {
		if t, ok := tagMap[id]; ok {
			result = append(result, t)
		}
	}
	return result, nil
}

// FindByArticleID 查询某篇文章关联的所有标签（JOIN article_tags）
func (r *tagRepository) FindByArticleID(ctx context.Context, articleID uint64) ([]*domain.Tag, error) {
	var list []model.Tag
	err := r.db.WithContext(ctx).
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Order("tags.id ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Tag, 0, len(list))
	for i := range list {
		result = append(result, tagModelToDomain(&list[i]))
	}
	return result, nil
}

// ─── 唯一性校验 ──────────────────────────────────────────────────────────────

// ExistsBySlug 检查 Slug 是否已存在
func (r *tagRepository) ExistsBySlug(ctx context.Context, slug string, excludeID uint64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.Tag{}).Where("slug = ?", slug)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// ExistsByName 检查名称是否已存在
func (r *tagRepository) ExistsByName(ctx context.Context, name string, excludeID uint64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.Tag{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// ─── 计数管理 ────────────────────────────────────────────────────────────────

// IncrArticleCount 原子性增减文章计数（delta 为负数时防止变负）
func (r *tagRepository) IncrArticleCount(ctx context.Context, id uint64, delta int) error {
	if delta < 0 {
		return r.db.WithContext(ctx).
			Model(&model.Tag{}).
			Where("id = ? AND article_count >= ?", id, -delta).
			Update("article_count", gorm.Expr("article_count + ?", delta)).Error
	}
	return r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Where("id = ?", id).
		Update("article_count", gorm.Expr("article_count + ?", delta)).Error
}

// UpdateArticleCount 直接设置文章计数
func (r *tagRepository) UpdateArticleCount(ctx context.Context, id uint64, count int) error {
	return r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Where("id = ?", id).
		Update("article_count", count).Error
}

// ─── 模型映射 ────────────────────────────────────────────────────────────────

func tagDomainToModel(t *domain.Tag) *model.Tag {
	return &model.Tag{
		ID:           t.ID,
		Name:         t.Name,
		Slug:         t.Slug,
		Color:        t.Color,
		ArticleCount: t.ArticleCount,
	}
}

func tagModelToDomain(m *model.Tag) *domain.Tag {
	return &domain.Tag{
		ID:           m.ID,
		Name:         m.Name,
		Slug:         m.Slug,
		Color:        m.Color,
		ArticleCount: m.ArticleCount,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
