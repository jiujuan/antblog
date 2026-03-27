package tag

import (
	"context"
	"regexp"

	apperrors "antblog/pkg/errors"
	"antblog/pkg/utils"
)

// colorHexRegexp 合法的 Hex 颜色值校验（如 #00ADD8 或 #FFF）
var colorHexRegexp = regexp.MustCompile(`^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$`)

const defaultTagColor = "#6B7280"

// IDomainService 标签领域服务接口
type IDomainService interface {
	// ValidateCreate 校验创建参数（名称/Slug 唯一性 + 颜色格式）
	ValidateCreate(ctx context.Context, name, slug, color string) error

	// ValidateUpdate 校验更新参数（排除自身的唯一性）
	ValidateUpdate(ctx context.Context, id uint64, name, slug, color string) error

	// BuildTag 构建新标签实体
	BuildTag(name, slug, color string) (*Tag, error)

	// EnsureSlug 若 slug 为空则由 name 自动生成，并确保唯一
	EnsureSlug(ctx context.Context, name, slug string, excludeID uint64) (string, error)

	// NormalizeColor 规范化颜色值，若为空或格式无效则返回默认颜色
	NormalizeColor(color string) string
}

// DomainService 标签领域服务实现
type DomainService struct {
	repo ITagRepository
}

// NewDomainService 创建标签领域服务
func NewDomainService(repo ITagRepository) IDomainService {
	return &DomainService{repo: repo}
}

// ValidateCreate 校验创建参数
func (s *DomainService) ValidateCreate(ctx context.Context, name, slug, color string) error {
	if name == "" {
		return apperrors.ErrInvalidParams("标签名称不能为空")
	}

	// 名称唯一性
	exists, err := s.repo.ExistsByName(ctx, name, 0)
	if err != nil {
		return apperrors.ErrInternalError(err)
	}
	if exists {
		return apperrors.New(apperrors.CodeSlugAlreadyExists, "标签名称已存在")
	}

	// Slug 唯一性（有值时才校验）
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

	// 颜色格式校验
	if color != "" && !colorHexRegexp.MatchString(color) {
		return apperrors.ErrInvalidParams("颜色值格式无效，请使用 Hex 格式（如 #00ADD8）")
	}

	return nil
}

// ValidateUpdate 校验更新参数（排除自身）
func (s *DomainService) ValidateUpdate(ctx context.Context, id uint64, name, slug, color string) error {
	if name == "" {
		return apperrors.ErrInvalidParams("标签名称不能为空")
	}

	exists, err := s.repo.ExistsByName(ctx, name, id)
	if err != nil {
		return apperrors.ErrInternalError(err)
	}
	if exists {
		return apperrors.New(apperrors.CodeSlugAlreadyExists, "标签名称已存在")
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

	if color != "" && !colorHexRegexp.MatchString(color) {
		return apperrors.ErrInvalidParams("颜色值格式无效，请使用 Hex 格式（如 #00ADD8）")
	}

	return nil
}

// BuildTag 构建标签实体（不含 ID，由仓储层赋值）
func (s *DomainService) BuildTag(name, slug, color string) (*Tag, error) {
	if name == "" {
		return nil, apperrors.ErrInvalidParams("标签名称不能为空")
	}
	return &Tag{
		Name:  name,
		Slug:  slug,
		Color: s.NormalizeColor(color),
	}, nil
}

// EnsureSlug 若 slug 为空则由 name 自动生成唯一 slug
func (s *DomainService) EnsureSlug(ctx context.Context, name, slug string, excludeID uint64) (string, error) {
	if slug != "" {
		return slug, nil
	}

	base := utils.Slugify(name)
	if base == "" {
		base = "tag"
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

// NormalizeColor 规范化颜色值，空值或无效格式返回默认灰色
func (s *DomainService) NormalizeColor(color string) string {
	if color == "" || !colorHexRegexp.MatchString(color) {
		return defaultTagColor
	}
	return color
}
