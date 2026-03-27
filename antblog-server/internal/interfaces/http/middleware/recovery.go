package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	apperrors "antblog/pkg/errors"
	"antblog/pkg/response"
)

// Recovery panic 捕获中间件，防止 goroutine 崩溃，记录堆栈并返回 500
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.ByteString("stack", stack),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code: apperrors.CodeInternalError,
					Msg:  apperrors.Message(apperrors.CodeInternalError),
				})
			}
		}()
		c.Next()
	}
}
