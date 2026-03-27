package category

import (
	"context"

	apperrors "antblog/pkg/errors"
	"antblog/pkg/utils"
)

// IDomainService 分类领域服务接口
type IDomainService interface {
	// ValidateCreate 校验创建参数（唯一性等业务规则）
	ValidateCreate(ctx context.Context, name, slug string) error

	// ValidateUpdate 校验更新参数（排除自身的唯一性）
	ValidateUpdate(ctx context.Context, id uint64, name, slug string) error

	// BuildCategory 构建新分类实体（含 Slug 生成）
	BuildCategory(name, slug, description, cover string, sortOrder int) (*Category, error)

	// EnsureSlug 若 slug 为空则由 name 自动生成，并确保唯一
	EnsureSlug(ctx context.Context, name, slug string, excludeID uint64) (string, error)
}

// DomainService 分类领域服务实现
type DomainService struct {
	repo ICategoryRepository
}

// NewDomainService 创建分类领域服务
func NewDomainService(repo ICategoryRepository) IDomainService {
	return &DomainService{repo: repo}
}

// ValidateCreate 校验创建时的名称/slug 唯一性
func (s *DomainService) ValidateCreate(ctx context.Context, name, slug string) error {
	if name == "" {
		return apperrors.ErrInvalidParams("分类名称不能为空")
	}

	exists, err := s.repo.ExistsByName(ctx, name, 0)
	if err != nil {
		return apperrors.ErrInternalError(err)
	}
	if exists {
		return apperrors.New(apperrors.CodeSlugAlreadyExists, "分类名称已存在")
	}

	if slug != "" {
		if !utils.IsValidSlug(slug) {
			return apperrors.ErrInvalidParams("Slug 格式无效，只能包含小写字母、数字和连字符")
		}
		exists, err = s.repo.ExistsBySlug(ctx, slug, 0)
		if err != nil {
			return apperrors.ErrInternalError(err)
		}
		if exists {
			return apperrors.FromCode(apperrors.CodeSlugAlreadyExists)
		}
	}
	return nil
}

// ValidateUpdate 校验更新时排除自身的唯一性
func (s *DomainService) ValidateUpdate(ctx context.Context, id uint64, name, slug string) error {
	if name == "" {
		return apperrors.ErrInvalidParams("分类名称不能为空")
	}

	exists, err := s.repo.ExistsByName(ctx, name, id)
	if err != nil {
		return apperrors.ErrInternalError(err)
	}
	if exists {
		return apperrors.New(apperrors.CodeSlugAlreadyExists, "分类名称已存在")
	}

	if slug != "" {
		if !utils.IsValidSlug(slug) {
			return apperrors.ErrInvalidParams("Slug 格式无效，只能包含小写字母、数字和连字符")
		}
		exists, err = s.repo.ExistsBySlug(ctx, slug, id)
		if err != nil {
			return apperrors.ErrInternalError(err)
		}
		if exists {
			return apperrors.FromCode(apperrors.CodeSlugAlreadyExists)
		}
	}
	return nil
}

// BuildCategory 构建分类实体
func (s *DomainService) BuildCategory(name, slug, description, cover string, sortOrder int) (*Category, error) {
	if name == "" {
		return nil, apperrors.ErrInvalidParams("分类名称不能为空")
	}
	return &Category{
		Name:        name,
		Slug:        slug,
		Description: description,
		Cover:       cover,
		SortOrder:   sortOrder,
	}, nil
}

// EnsureSlug 若 slug 为空则由 name 自动生成唯一 slug
func (s *DomainService) EnsureSlug(ctx context.Context, name, slug string, excludeID uint64) (string, error) {
	if slug != "" {
		return slug, nil
	}

	base := utils.Slugify(name)
	if base == "" {
		base = "category"
	}

	candidate := base
	for i := 1; i <= 20; i++ {
		exists, err := s.repo.ExistsBySlug(ctx, candidate, excludeID)
		if err != nil {
			return "", apperrors.ErrInternalError(err)
		}
		if !exists {
			return candidate, nil
		}
		candidate = utils.SlugifyWithSuffix(name, i)
	}
	return "", apperrors.ErrInvalidParams("无法生成唯一的 Slug，请手动指定")
}
