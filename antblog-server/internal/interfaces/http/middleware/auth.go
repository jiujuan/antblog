// Package middleware HTTP 中间件集合。
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	infracache "antblog/internal/infrastructure/cache"
	"antblog/pkg/jwt"
	"antblog/pkg/response"
)

const (
	CtxKeyUserID   = "user_id"
	CtxKeyUsername = "username"
	CtxKeyRole     = "role"
	CtxKeyClaims   = "claims"

	bearerPrefix = "Bearer "
)

// JWTAuth JWT 鉴权中间件（必须登录）
func JWTAuth(tokenMgr jwt.ITokenManager, userCache infracache.IUserCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 检查 Token 是否在黑名单
		if userCache != nil {
			blocked, _ := userCache.IsTokenBlocked(c.Request.Context(), token)
			if blocked {
				response.Unauthorized(c)
				c.Abort()
				return
			}
		}

		claims, err := tokenMgr.ParseAccessToken(token)
		if err != nil {
			response.FailWithError(c, err)
			c.Abort()
			return
		}

		c.Set(CtxKeyUserID, claims.UserID)
		c.Set(CtxKeyUsername, claims.Username)
		c.Set(CtxKeyRole, claims.Role)
		c.Set(CtxKeyClaims, claims)
		c.Next()
	}
}

// AdminAuth 管理员权限中间件（必须是 admin 角色，需在 JWTAuth 之后使用）
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(CtxKeyRole)
		if !exists {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		if roleInt, ok := role.(int); !ok || roleInt != int(jwt.RoleAdmin) {
			response.Forbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// OptionalAuth 可选鉴权中间件（登录态注入，未登录也放行）
func OptionalAuth(tokenMgr jwt.ITokenManager, userCache infracache.IUserCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			c.Next()
			return
		}
		// Token 黑名单校验（失效的 token 不注入用户信息）
		if userCache != nil {
			if blocked, _ := userCache.IsTokenBlocked(c.Request.Context(), token); blocked {
				c.Next()
				return
			}
		}
		claims, err := tokenMgr.ParseAccessToken(token)
		if err == nil {
			c.Set(CtxKeyUserID, claims.UserID)
			c.Set(CtxKeyUsername, claims.Username)
			c.Set(CtxKeyRole, claims.Role)
			c.Set(CtxKeyClaims, claims)
		}
		c.Next()
	}
}

// ─── 辅助函数 ────────────────────────────────────────────────────────────────

// extractBearerToken 从 Authorization header 或 query 参数提取 Bearer Token
func extractBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		auth = c.Query("token")
		if auth == "" {
			return ""
		}
		return auth
	}
	if strings.HasPrefix(auth, bearerPrefix) {
		return strings.TrimPrefix(auth, bearerPrefix)
	}
	return ""
}

// GetCurrentUserID 从 gin.Context 获取当前用户 ID
func GetCurrentUserID(c *gin.Context) (uint64, bool) {
	v, exists := c.Get(CtxKeyUserID)
	if !exists {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}

// GetOptionalUserID 获取可选用户 ID（未登录时返回 nil）
func GetOptionalUserID(c *gin.Context) *uint64 {
	id, ok := GetCurrentUserID(c)
	if !ok {
		return nil
	}
	return &id
}

// GetCurrentRole 从 gin.Context 获取当前用户角色
func GetCurrentRole(c *gin.Context) (int, bool) {
	v, exists := c.Get(CtxKeyRole)
	if !exists {
		return 0, false
	}
	role, ok := v.(int)
	return role, ok
}

// MustGetCurrentUserID 获取当前用户 ID，不存在则 panic（内部使用，需确保已过 JWTAuth）
func MustGetCurrentUserID(c *gin.Context) uint64 {
	id, ok := GetCurrentUserID(c)
	if !ok {
		panic("user_id not found in context, JWTAuth middleware missing?")
	}
	return id
}
