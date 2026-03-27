package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"antblog/pkg/cache"
	apperrors "antblog/pkg/errors"
	"antblog/pkg/response"
)

// RateLimit 基于 Redis 滑动窗口的限流中间件
// limit: 窗口内最大请求数，window: 时间窗口大小
func RateLimit(c cache.ICache, limit int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := fmt.Sprintf(cache.KeyRateLimit, ctx.ClientIP(), ctx.FullPath())

		count, err := c.Incr(context.Background(), key)
		if err != nil {
			// 缓存不可用时放行（降级策略）
			ctx.Next()
			return
		}

		// 首次请求时设置过期时间
		if count == 1 {
			_ = c.Expire(context.Background(), key, window)
		}

		if int(count) > limit {
			response.Fail(ctx, apperrors.CodeTooManyRequests, apperrors.Message(apperrors.CodeTooManyRequests))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
