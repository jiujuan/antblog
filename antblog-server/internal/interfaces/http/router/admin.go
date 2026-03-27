package router

import (
	"github.com/gin-gonic/gin"

	"antblog/internal/interfaces/http/middleware"
)

// registerAdmin 注册后台 /api/admin 路由（全部需要 Admin 权限）
func registerAdmin(rg *gin.RouterGroup, deps RouterDeps) {
	rg.Use(
		middleware.JWTAuth(deps.TokenMgr, deps.UserCache),
		middleware.AdminAuth(),
	)

	// ── 分类管理 ─────────────────────────────────────────────────────────
	categories := rg.Group("/categories")
	{
		categories.GET("", deps.AdminCategoryHandler.ListCategories)
		categories.GET("/:id", deps.AdminCategoryHandler.GetCategory)
		categories.POST("", deps.AdminCategoryHandler.CreateCategory)
		categories.PUT("/:id", deps.AdminCategoryHandler.UpdateCategory)
		categories.DELETE("/:id", deps.AdminCategoryHandler.DeleteCategory)
	}

	// ── 标签管理 ─────────────────────────────────────────────────────────
	tags := rg.Group("/tags")
	{
		tags.GET("", deps.AdminTagHandler.ListTags)
		tags.GET("/:id", deps.AdminTagHandler.GetTag)
		tags.POST("", deps.AdminTagHandler.CreateTag)
		tags.POST("/batch", deps.AdminTagHandler.BatchCreateTags)
		tags.PUT("/:id", deps.AdminTagHandler.UpdateTag)
		tags.DELETE("/:id", deps.AdminTagHandler.DeleteTag)
	}

	// ── 文章管理 ─────────────────────────────────────────────────────────
	articles := rg.Group("/articles")
	{
		articles.GET("", deps.AdminArticleHandler.ListArticles)
		articles.GET("/:id", deps.AdminArticleHandler.GetArticle)
		articles.POST("", deps.AdminArticleHandler.CreateArticle)
		articles.PUT("/:id", deps.AdminArticleHandler.UpdateArticle)
		articles.PATCH("/:id/status", deps.AdminArticleHandler.UpdateArticleStatus)
		articles.DELETE("/:id", deps.AdminArticleHandler.DeleteArticle)
	}

	// ── 评论管理 ─────────────────────────────────────────────────────────
	comments := rg.Group("/comments")
	{
		comments.GET("", deps.AdminCommentHandler.ListComments)
		comments.GET("/:id", deps.AdminCommentHandler.GetComment)
		comments.PATCH("/:id/status", deps.AdminCommentHandler.UpdateStatus)
		comments.DELETE("/:id", deps.AdminCommentHandler.DeleteComment)
	}

	// ── 媒体管理 ─────────────────────────────────────────────────────────
	// POST /api/admin/media/upload  上传文件
	// GET  /api/admin/media         列表
	// GET  /api/admin/media/:id     详情
	// PATCH /api/admin/media/:id/bind  绑定/解绑文章
	// DELETE /api/admin/media/:id   删除
	media := rg.Group("/media")
	{
		media.POST("/upload", deps.AdminMediaHandler.Upload)
		media.GET("", deps.AdminMediaHandler.ListMedia)
		media.GET("/:id", deps.AdminMediaHandler.GetMedia)
		media.PATCH("/:id/bind", deps.AdminMediaHandler.BindArticle)
		media.DELETE("/:id", deps.AdminMediaHandler.DeleteMedia)
	}
}
