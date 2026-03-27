antblog 蚂蚁博客

## antblog蚂蚁博客功能

博客antblog功能：

普通用户模块：

普通用户只能查看文章，发表评论、点缀、收藏文章。

- 用户注册、登录、登出

文章模块：

- 文章列表展示（分页）
- 文章详情页
- 文章tag
- 文章分类

admin管理员后台管理：

管理员后台发表文章，管理文章，管理评论等功能。

后台管理员默认账号：

用户名：admin
密码：Admin@2026

- 文章分类、归档（按时间线）

- 文章的增、删、改、查（markdown编辑器）
- 文章tag标签
- 文章点赞、收藏管理
- 文章的图片管理
- 文章评论管理

## 技术栈

### 后端技术

用 Go 语言，结合 Gin + GORM + Zap + Viper +validator+go-jwt+go-redis 技术栈。

- MySQL
- Redis

### 前端技术

Vue3 + Tailwind CSS + shadcn/ui +Vite+TS

## 架构设计文档

[后端架构架构设计文档](./antblog-server/docs/antblog-backend-structure.md)


[前端架构架构设计文档](./antblog-server/docs/antblog-frontend-structure.md)

## 数据库

[数据库建表语句](./antblog-server/docs/sqls/sqls.sql)

## 博客页面预览

