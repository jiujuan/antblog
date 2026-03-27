// Package router 路由注册入口，使用 uber/fx 注入所有依赖。
package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	infracache "antblog/internal/infrastructure/cache"
	adminhandler "antblog/internal/interfaces/http/admin"
	"antblog/internal/interfaces/http/handler"
	"antblog/internal/interfaces/http/middleware"
	"antblog/pkg/cache"
	"antblog/pkg/config"
	"antblog/pkg/jwt"
)

// RouterDeps fx 路由依赖
type RouterDeps struct {
	fx.In
	Config    *config.Config
	Logger    *zap.Logger
	TokenMgr  jwt.ITokenManager
	Cache     cache.ICache
	UserCache infracache.IUserCache

	// 前台 Handlers
	UserHandler     *handler.UserHandler
	CategoryHandler *handler.CategoryHandler
	TagHandler      *handler.TagHandler
	ArticleHandler  *handler.ArticleHandler
	CommentHandler  *handler.CommentHandler
	MediaHandler    *handler.MediaHandler

	// 后台 Handlers
	AdminCategoryHandler *adminhandler.CategoryHandler
	AdminTagHandler      *adminhandler.TagHandler
	AdminArticleHandler  *adminhandler.ArticleHandler
	AdminCommentHandler  *adminhandler.CommentHandler
	AdminMediaHandler    *adminhandler.MediaHandler
}

// NewRouter 创建并配置 Gin 引擎
func NewRouter(deps RouterDeps) *gin.Engine {
	if !deps.Config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(
		middleware.Recovery(deps.Logger),
		middleware.Logger(deps.Logger),
		middleware.CORS(),
	)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": deps.Config.App.Version})
	})

	// 静态文件服务：对外暴露上传目录
	r.Static("/uploads", deps.Config.Upload.LocalPath)

	v1 := r.Group("/api/v1")
	registerAPIV1(v1, deps)

	admin := r.Group("/api/admin")
	registerAdmin(admin, deps)

	deps.Logger.Info("router initialized",
		zap.String("addr", fmt.Sprintf("%s:%d", deps.Config.Server.Host, deps.Config.Server.Port)),
	)
	return r
}

// Module fx 路由模块
var Module = fx.Options(
	fx.Provide(NewRouter),
	// 前台 handlers
	fx.Provide(handler.NewUserHandler),
	fx.Provide(handler.NewCategoryHandler),
	fx.Provide(handler.NewTagHandler),
	fx.Provide(handler.NewArticleHandler),
	fx.Provide(handler.NewCommentHandler),
	fx.Provide(handler.NewMediaHandler),
	// 后台 handlers
	fx.Provide(adminhandler.NewAdminCategoryHandler),
	fx.Provide(adminhandler.NewAdminTagHandler),
	fx.Provide(adminhandler.NewAdminArticleHandler),
	fx.Provide(adminhandler.NewAdminCommentHandler),
	fx.Provide(adminhandler.NewAdminMediaHandler),
	// 缓存
	fx.Provide(infracache.NewUserCache),
	fx.Provide(infracache.NewArticleCache),
)
