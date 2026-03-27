// Package jwt 提供 JWT Access Token 和 Refresh Token 的生成与解析，
// 通过接口抽象保证可替换性，函数选项模式灵活配置。
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"

	"antblog/pkg/config"
	apperrors "antblog/pkg/errors"
)

// ─── 接口定义 ────────────────────────────────────────────────────────────────

// ITokenManager JWT Token 管理器接口
type ITokenManager interface {
	// GenerateAccessToken 生成 Access Token
	GenerateAccessToken(claims *UserClaims) (string, error)
	// GenerateRefreshToken 生成 Refresh Token
	GenerateRefreshToken(userID uint64, role int) (string, error)
	// ParseAccessToken 解析并验证 Access Token
	ParseAccessToken(tokenStr string) (*UserClaims, error)
	// ParseRefreshToken 解析并验证 Refresh Token
	ParseRefreshToken(tokenStr string) (*RefreshClaims, error)
}

// ─── Claims 结构 ─────────────────────────────────────────────────────────────

// UserClaims Access Token 载荷
type UserClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`     // 1=普通用户 2=管理员
	jwt.RegisteredClaims
}

// RefreshClaims Refresh Token 载荷（精简，减少体积）
type RefreshClaims struct {
	UserID uint64 `json:"user_id"`
	Role   int    `json:"role"`
	jwt.RegisteredClaims
}

// ─── 实现 ────────────────────────────────────────────────────────────────────

// TokenManager JWT Token 管理器实现
type TokenManager struct {
	opts *options
}

// New 创建 TokenManager
func New(opts ...Option) ITokenManager {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &TokenManager{opts: o}
}

// NewFromConfig 从 config.JWTConfig 创建
func NewFromConfig(cfg config.JWTConfig, opts ...Option) ITokenManager {
	baseOpts := []Option{
		WithSecret(cfg.Secret),
		WithAccessExpiry(time.Duration(cfg.AccessExpiry) * time.Minute),
		WithRefreshExpiry(time.Duration(cfg.RefreshExpiry) * 24 * time.Hour),
		WithIssuer(cfg.Issuer),
	}
	return New(append(baseOpts, opts...)...)
}

func (m *TokenManager) signingKey() []byte {
	return []byte(m.opts.secret)
}

func (m *TokenManager) signingMethod() jwt.SigningMethod {
	switch m.opts.signingMethod {
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}

// GenerateAccessToken 生成 Access Token
func (m *TokenManager) GenerateAccessToken(claims *UserClaims) (string, error) {
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    m.opts.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.opts.accessExpiry)),
	}

	token := jwt.NewWithClaims(m.signingMethod(), claims)
	signed, err := token.SignedString(m.signingKey())
	if err != nil {
		return "", fmt.Errorf("jwt: sign access token: %w", err)
	}
	return signed, nil
}

// GenerateRefreshToken 生成 Refresh Token
func (m *TokenManager) GenerateRefreshToken(userID uint64, role int) (string, error) {
	now := time.Now()
	claims := &RefreshClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.opts.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.opts.refreshExpiry)),
		},
	}

	token := jwt.NewWithClaims(m.signingMethod(), claims)
	signed, err := token.SignedString(m.signingKey())
	if err != nil {
		return "", fmt.Errorf("jwt: sign refresh token: %w", err)
	}
	return signed, nil
}

// ParseAccessToken 解析并验证 Access Token
func (m *TokenManager) ParseAccessToken(tokenStr string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.signingKey(), nil
	})

	if err != nil {
		return nil, mapJWTError(err)
	}
	if !token.Valid {
		return nil, apperrors.ErrTokenInvalid()
	}
	return claims, nil
}

// ParseRefreshToken 解析并验证 Refresh Token
func (m *TokenManager) ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.signingKey(), nil
	})

	if err != nil {
		return nil, mapJWTError(err)
	}
	if !token.Valid {
		return nil, apperrors.ErrTokenInvalid()
	}
	return claims, nil
}

// mapJWTError 将 jwt 库错误映射到业务错误
func mapJWTError(err error) error {
	if errors.Is(err, jwt.ErrTokenExpired) {
		return apperrors.ErrTokenExpired()
	}
	return apperrors.ErrTokenInvalid()
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 模块，依赖 *config.Config，提供 ITokenManager
var Module = fx.Options(
	fx.Provide(func(cfg *config.Config) ITokenManager {
		return NewFromConfig(cfg.JWT)
	}),
)

// ─── 角色常量（与 domain 层保持同步，避免循环依赖）────────────────────────────

const (
	RoleUser  = 1
	RoleAdmin = 2
)
