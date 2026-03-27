package router

import (
	"time"

	"github.com/gin-gonic/gin"

	"antblog/internal/interfaces/http/middleware"
)

// registerAPIV1 注册前台 /api/v1 路由
func registerAPIV1(rg *gin.RouterGroup, deps RouterDeps) {
	// ── Auth ──────────────────────────────────────────────────────────────
	auth := rg.Group("/auth")
	{
		loginLimit := middleware.RateLimit(deps.Cache, 10, time.Minute)
		auth.POST("/register", loginLimit, deps.UserHandler.Register)
		auth.POST("/login", loginLimit, deps.UserHandler.Login)
		auth.POST("/refresh", deps.UserHandler.RefreshToken)
		auth.POST("/logout",
			middleware.JWTAuth(deps.TokenMgr, deps.UserCache),
			deps.UserHandler.Logout,
		)
	}

	// ── 用户个人（需登录）────────────────────────────────────────────────
	jwtAuth := middleware.JWTAuth(deps.TokenMgr, deps.UserCache)
	user := rg.Group("/user", jwtAuth)
	{
		user.GET("/profile", deps.UserHandler.GetProfile)
		user.PUT("/profile", deps.UserHandler.UpdateProfile)
		user.PUT("/password", deps.UserHandler.ChangePassword)
		user.GET("/bookmarks", deps.ArticleHandler.GetUserBookmarks)
		user.GET("/media", deps.MediaHandler.ListMyMedia) // 当前用户的媒体库
	}

	// ── 分类（公开）──────────────────────────────────────────────────────
	categories := rg.Group("/categories")
	{
		categories.GET("", deps.CategoryHandler.ListCategories)
		categories.GET("/:slug", deps.CategoryHandler.GetCategoryBySlug)
	}

	// ── 标签（公开）──────────────────────────────────────────────────────
	tags := rg.Group("/tags")
	{
		tags.GET("", deps.TagHandler.ListTags)
		tags.GET("/:slug", deps.TagHandler.GetTagBySlug)
	}

	// ── 文章（公开 + 登录互动）──────────────────────────────────────────
	optAuth := middleware.OptionalAuth(deps.TokenMgr, deps.UserCache)
	articles := rg.Group("/articles")
	{
		// 列表/精选/归档 加 OptionalAuth：登录用户获得 liked/bookmarked 状态
		articles.GET("", optAuth, deps.ArticleHandler.ListArticles)
		articles.GET("/featured", optAuth, deps.ArticleHandler.GetFeaturedArticles)
		articles.GET("/archive", deps.ArticleHandler.GetArchive)
		articles.GET("/archive/detail", optAuth, deps.ArticleHandler.GetArchiveDetail)
		articles.GET("/:slug", optAuth, deps.ArticleHandler.GetArticleBySlug)
		articles.POST("/:id/like", jwtAuth, deps.ArticleHandler.LikeArticle)
		articles.DELETE("/:id/like", jwtAuth, deps.ArticleHandler.UnlikeArticle)
		articles.POST("/:id/bookmark", jwtAuth, deps.ArticleHandler.BookmarkArticle)
		articles.DELETE("/:id/bookmark", jwtAuth, deps.ArticleHandler.UnbookmarkArticle)
	}

	// ── 评论（公开读取，可选登录发表）──────────────────────────────────
	comments := rg.Group("/comments", optAuth)
	{
		comments.GET("", deps.CommentHandler.ListComments)
		comments.POST("", deps.CommentHandler.CreateComment)
		comments.POST("/:id/like", deps.CommentHandler.LikeComment)
	}
}
