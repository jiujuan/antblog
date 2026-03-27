# AntBlog 前端项目结构

> 技术栈：Vue3 + TypeScript + Tailwind CSS + shadcn/ui + Vite

```
antblog-web/
├── public/
│   └── favicon.ico
│
├── src/
│   ├── main.ts                         # 应用入口
│   ├── App.vue                         # 根组件
│   │
│   ├── assets/
│   │   ├── styles/
│   │   │   ├── main.css                # Tailwind 指令 + 全局样式
│   │   │   └── markdown.css            # Markdown 内容样式
│   │   └── images/
│   │
│   ├── router/
│   │   ├── index.ts                    # 路由实例 + 全局守卫
│   │   ├── routes/
│   │   │   ├── blog.ts                 # 前台路由
│   │   │   └── admin.ts                # 后台路由
│   │   └── guards/
│   │       └── auth.guard.ts           # 登录态路由守卫
│   │
│   ├── stores/                         # Pinia 状态管理
│   │   ├── auth.store.ts               # 用户认证状态（token, user info）
│   │   ├── article.store.ts            # 文章列表/详情缓存
│   │   ├── category.store.ts
│   │   └── tag.store.ts
│   │
│   ├── api/                            # API 层（与后端对应）
│   │   ├── http.ts                     # axios 实例（拦截器、token注入）
│   │   ├── auth.api.ts                 # register / login / logout
│   │   ├── article.api.ts              # 文章 CRUD
│   │   ├── category.api.ts
│   │   ├── tag.api.ts
│   │   ├── comment.api.ts
│   │   └── media.api.ts
│   │
│   ├── types/                          # TypeScript 类型定义
│   │   ├── api.types.ts                # 通用响应 ApiResponse<T>
│   │   ├── auth.types.ts               # LoginReq, UserInfo
│   │   ├── article.types.ts            # Article, ArticleList, CreateArticleDto
│   │   ├── category.types.ts
│   │   ├── tag.types.ts
│   │   └── comment.types.ts
│   │
│   ├── composables/                    # Vue3 组合式函数
│   │   ├── useAuth.ts                  # 登录/登出逻辑
│   │   ├── useArticleList.ts           # 分页文章列表
│   │   ├── useArticleDetail.ts         # 文章详情 + 点赞收藏
│   │   ├── usePagination.ts            # 通用分页
│   │   ├── useTheme.ts                 # 主题切换（暗/亮）
│   │   └── useUpload.ts                # 图片上传
│   │
│   ├── components/                     # 通用组件
│   │   ├── ui/                         # shadcn/ui 组件（自动生成）
│   │   │   ├── button/
│   │   │   ├── card/
│   │   │   ├── dialog/
│   │   │   ├── input/
│   │   │   ├── badge/
│   │   │   ├── pagination/
│   │   │   ├── toast/
│   │   │   └── ...
│   │   ├── common/
│   │   │   ├── AppHeader.vue           # 顶部导航栏
│   │   │   ├── AppFooter.vue
│   │   │   ├── AppSidebar.vue          # 分类/标签侧边栏
│   │   │   ├── ThemeToggle.vue         # 暗/亮主题切换
│   │   │   ├── LoadingSpinner.vue
│   │   │   └── EmptyState.vue
│   │   ├── article/
│   │   │   ├── ArticleCard.vue         # 文章列表卡片
│   │   │   ├── ArticleList.vue         # 文章列表容器（含分页）
│   │   │   ├── ArticleMeta.vue         # 时间/分类/标签/阅读数
│   │   │   ├── ArticleContent.vue      # Markdown 渲染
│   │   │   ├── ArticleTags.vue         # 标签组
│   │   │   ├── ArticleLike.vue         # 点赞按钮
│   │   │   └── ArticleBookmark.vue     # 收藏按钮
│   │   └── comment/
│   │       ├── CommentList.vue
│   │       ├── CommentItem.vue
│   │       └── CommentForm.vue
│   │
│   ├── views/                          # 页面视图
│   │   ├── blog/                       # 前台页面
│   │   │   ├── HomeView.vue            # 首页（文章列表）
│   │   │   ├── ArticleDetailView.vue   # 文章详情
│   │   │   ├── CategoryView.vue        # 分类页
│   │   │   ├── TagView.vue             # 标签页
│   │   │   ├── ArchiveView.vue         # 时间线归档
│   │   │   └── SearchView.vue          # 搜索结果
│   │   ├── auth/
│   │   │   ├── LoginView.vue
│   │   │   └── RegisterView.vue
│   │   └── admin/                      # 后台管理页面
│   │       ├── AdminLayout.vue         # 后台布局（侧边栏导航）
│   │       ├── DashboardView.vue       # 仪表盘
│   │       ├── article/
│   │       │   ├── ArticleListView.vue # 文章管理列表
│   │       │   └── ArticleEditView.vue # 文章编辑（Markdown 编辑器）
│   │       ├── CategoryManageView.vue
│   │       ├── TagManageView.vue
│   │       ├── CommentManageView.vue
│   │       └── MediaManageView.vue     # 图片管理
│   │
│   └── utils/
│       ├── format.ts                   # 日期/数字格式化
│       ├── markdown.ts                 # markdown-it 配置实例
│       └── storage.ts                  # localStorage 封装
│
├── index.html
├── vite.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── tsconfig.app.json
├── components.json                     # shadcn/ui 配置
├── package.json
└── README.md
```

---

## 关键文件说明

### `src/api/http.ts` — axios 实例
```typescript
// 请求拦截：自动注入 Authorization: Bearer <token>
// 响应拦截：统一处理 401 跳登录、业务错误 toast 提示
```

### `src/types/api.types.ts` — 通用响应类型
```typescript
interface ApiResponse<T> {
  code: number
  msg: string
  data: T
}
interface PageResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}
```

### 路由结构
```
/                    → HomeView（文章列表）
/articles/:id        → ArticleDetailView
/categories/:slug    → CategoryView
/tags/:slug          → TagView
/archive             → ArchiveView
/login               → LoginView
/register            → RegisterView

/admin               → AdminLayout（需登录）
/admin/dashboard     → DashboardView
/admin/articles      → ArticleListView
/admin/articles/new  → ArticleEditView
/admin/articles/:id  → ArticleEditView
/admin/categories    → CategoryManageView
/admin/tags          → TagManageView
/admin/comments      → CommentManageView
/admin/media         → MediaManageView
```

### 主要 npm 依赖
```json
{
  "dependencies": {
    "vue": "^3.x",
    "vue-router": "^4.x",
    "pinia": "^2.x",
    "axios": "^1.x",
    "markdown-it": "^14.x",
    "@vueuse/core": "^10.x",
    "lucide-vue-next": "latest"
  },
  "devDependencies": {
    "typescript": "^5.x",
    "vite": "^5.x",
    "@vitejs/plugin-vue": "^5.x",
    "tailwindcss": "^3.x",
    "autoprefixer": "^10.x",
    "@types/node": "latest"
  }
}
```
