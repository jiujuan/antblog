# AntBlog 后端项目结构

> 技术栈：Go + Gin + GORM + Zap + Viper + Validator + go-jwt + go-redis + uber/fx
> 架构：DDD + Clean Architecture + 依赖倒置

```
antblog-server/
├── cmd/
│   └── server/
│       └── main.go                     # 入口：fx.New() 注册所有模块
│
├── config/
│   ├── config.yaml                     # 主配置文件
│   ├── config.dev.yaml
│   └── config.prod.yaml
│
├── internal/                           # 业务核心（DDD 分层）
│   │
│   ├── domain/                         # 领域层（纯业务逻辑，零依赖）
│   │   ├── user/
│   │   │   ├── entity.go               # User 实体
│   │   │   ├── value_object.go         # Email, Password 值对象
│   │   │   ├── repository.go           # IUserRepository 接口
│   │   │   └── service.go              # 领域服务（纯业务规则）
│   │   ├── article/
│   │   │   ├── entity.go               # Article 实体
│   │   │   ├── value_object.go         # Status, Content 值对象
│   │   │   ├── repository.go           # IArticleRepository 接口
│   │   │   └── service.go
│   │   ├── category/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── tag/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── comment/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── media/
│   │       ├── entity.go               # 图片/媒体实体
│   │       ├── repository.go
│   │       └── service.go
│   │
│   ├── application/                    # 应用层（用例编排，依赖领域接口）
│   │   ├── user/
│   │   │   ├── dto.go                  # RegisterReq, LoginReq, UserResp
│   │   │   ├── usecase.go              # IUserUseCase 接口
│   │   │   └── usecase_impl.go         # 实现（注入 IUserRepository）
│   │   ├── article/
│   │   │   ├── dto.go                  # CreateArticleReq, ArticleListResp
│   │   │   ├── usecase.go
│   │   │   └── usecase_impl.go
│   │   ├── category/
│   │   │   ├── dto.go
│   │   │   ├── usecase.go
│   │   │   └── usecase_impl.go
│   │   ├── tag/
│   │   │   ├── dto.go
│   │   │   ├── usecase.go
│   │   │   └── usecase_impl.go
│   │   ├── comment/
│   │   │   ├── dto.go
│   │   │   ├── usecase.go
│   │   │   └── usecase_impl.go
│   │   └── media/
│   │       ├── dto.go
│   │       ├── usecase.go
│   │       └── usecase_impl.go
│   │
│   ├── infrastructure/                 # 基础设施层（实现领域接口）
│   │   ├── persistence/                # 数据库实现
│   │   │   ├── model/                  # GORM 模型（与 domain entity 解耦）
│   │   │   │   ├── user.go
│   │   │   │   ├── article.go
│   │   │   │   ├── category.go
│   │   │   │   ├── tag.go
│   │   │   │   ├── comment.go
│   │   │   │   └── media.go
│   │   │   ├── user_repo.go            # 实现 IUserRepository
│   │   │   ├── article_repo.go
│   │   │   ├── category_repo.go
│   │   │   ├── tag_repo.go
│   │   │   ├── comment_repo.go
│   │   │   └── media_repo.go
│   │   ├── cache/
│   │   │   ├── article_cache.go        # 文章缓存（Redis/Ristretto）
│   │   │   └── user_cache.go
│   │   └── storage/
│   │       └── local_storage.go        # 本地图片存储
│   │
│   └── interfaces/                     # 接口层（HTTP handler）
│       ├── http/
│       │   ├── middleware/
│       │   │   ├── auth.go             # JWT 鉴权中间件
│       │   │   ├── cors.go
│       │   │   ├── logger.go
│       │   │   ├── rate_limit.go
│       │   │   └── recovery.go
│       │   ├── handler/
│       │   │   ├── user_handler.go     # 注册/登录/登出
│       │   │   ├── article_handler.go  # 前台文章接口
│       │   │   ├── category_handler.go
│       │   │   ├── tag_handler.go
│       │   │   ├── comment_handler.go
│       │   │   └── media_handler.go
│       │   ├── admin/                  # 后台管理接口
│       │   │   ├── article_handler.go  # CRUD + Markdown
│       │   │   ├── category_handler.go
│       │   │   ├── tag_handler.go
│       │   │   ├── comment_handler.go
│       │   │   └── media_handler.go
│       │   └── router/
│       │       ├── router.go           # 路由注册入口
│       │       ├── api.go              # /api/v1 前台路由
│       │       └── admin.go            # /api/admin 后台路由
│       └── validator/
│           └── custom_validator.go     # 自定义校验规则
│
├── pkg/                                # 公共库（可复用，业务无关）
│   ├── config/
│   │   ├── config.go                   # Viper 封装，函数选项模式
│   │   └── options.go                  # WithConfigPath, WithEnvPrefix...
│   ├── db/
│   │   ├── db.go                       # GORM 封装（MySQL/PostgreSQL）
│   │   ├── options.go                  # WithDSN, WithMaxOpenConns...
│   │   └── migrate.go                  # 自动迁移
│   ├── logger/
│   │   ├── logger.go                   # Zap 封装
│   │   └── options.go                  # WithLevel, WithOutput...
│   ├── jwt/
│   │   ├── jwt.go                      # JWT 生成/解析接口
│   │   └── options.go                  # WithSecret, WithExpiry...
│   ├── crypto/
│   │   └── crypto.go                   # bcrypt 密码加密
│   ├── cache/
│   │   ├── cache.go                    # ICache 接口
│   │   ├── redis.go                    # Redis 实现
│   │   ├── ristretto.go                # 本地缓存实现
│   │   └── options.go
│   ├── response/
│   │   └── response.go                 # 统一响应结构 {code, msg, data}
│   ├── errors/
│   │   ├── errors.go                   # 业务错误码定义
│   │   └── codes.go                    # ErrUserNotFound, ErrUnauthorized...
│   ├── validator/
│   │   └── validator.go                # go-playground/validator 封装
│   └── utils/
│       ├── page.go                     # 分页工具
│       ├── slug.go                     # URL slug 生成
│       └── time.go                     # 时间工具
│
├── migrations/                         # SQL 迁移文件（可选）
│   ├── 001_create_users.sql
│   ├── 002_create_articles.sql
│   └── ...
│
├── scripts/
│   ├── build.sh
│   └── dev.sh
│
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## 关键文件说明

### `cmd/server/main.go` — fx 依赖注入入口
```go
func main() {
    fx.New(
        config.Module,
        logger.Module,
        db.Module,
        cache.Module,
        // Infrastructure
        persistence.Module,
        // Application
        userapp.Module,
        articleapp.Module,
        // Interface
        router.Module,
    ).Run()
}
```

### 依赖流向
```
interfaces → application → domain ← infrastructure
                ↑                        ↓
            (use case)          (实现 domain 接口)
```

### 路由设计
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout

GET    /api/v1/articles          # 列表（分页）
GET    /api/v1/articles/:id      # 详情
GET    /api/v1/articles/archive  # 时间线归档
GET    /api/v1/categories
GET    /api/v1/tags

# 后台（需 JWT Admin）
GET    /api/admin/articles
POST   /api/admin/articles
PUT    /api/admin/articles/:id
DELETE /api/admin/articles/:id
POST   /api/admin/media/upload
GET    /api/admin/comments
DELETE /api/admin/comments/:id
```
