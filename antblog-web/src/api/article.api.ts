import http, { adminHttp } from './http'
import type { PageResult } from '@/types/api.types'
import type {
  Article,
  ArticleListItem,
  ArchiveItem,
  ListArticleQuery,
  ArchiveDetailQuery,
  CreateArticleReq,
  UpdateArticleReq,
  UpdateStatusReq,
  AdminListArticleQuery,
} from '@/types/article.types'

export const articleApi = {
  // ── 前台 ──────────────────────────────────────────────────────────────────

  /** 文章列表 */
  list: (params?: ListArticleQuery) =>
    http.get<any, PageResult<ArticleListItem>>('/articles', { params }),

  /** 精选文章 */
  featured: (limit = 6) =>
    http.get<any, ArticleListItem[]>('/articles/featured', { params: { limit } }),

  /** 归档时间线 */
  archive: () =>
    http.get<any, ArchiveItem[]>('/articles/archive'),

  /** 归档详情 */
  archiveDetail: (params: ArchiveDetailQuery) =>
    http.get<any, PageResult<ArticleListItem>>('/articles/archive/detail', { params }),

  /** 文章详情（by slug） */
  detail: (slug: string) =>
    http.get<any, Article>(`/articles/${slug}`),

  /** 点赞 */
  like: (id: number) =>
    http.post<any, null>(`/articles/${id}/like`),

  /** 取消点赞 */
  unlike: (id: number) =>
    http.delete<any, null>(`/articles/${id}/like`),

  /** 收藏 */
  bookmark: (id: number) =>
    http.post<any, null>(`/articles/${id}/bookmark`),

  /** 取消收藏 */
  unbookmark: (id: number) =>
    http.delete<any, null>(`/articles/${id}/bookmark`),

  /** 我的收藏列表 */
  myBookmarks: (params?: { page?: number; page_size?: number }) =>
    http.get<any, PageResult<ArticleListItem>>('/user/bookmarks', { params }),

  // ── 后台 ──────────────────────────────────────────────────────────────────

  /** 后台文章列表 */
  adminList: (params?: AdminListArticleQuery) =>
    adminHttp.get<any, PageResult<ArticleListItem>>('/articles', { params }),

  /** 后台文章详情 */
  adminDetail: (id: number) =>
    adminHttp.get<any, Article>(`/articles/${id}`),

  /** 创建文章 */
  create: (data: CreateArticleReq) =>
    adminHttp.post<any, Article>('/articles', data),

  /** 更新文章 */
  update: (id: number, data: UpdateArticleReq) =>
    adminHttp.put<any, Article>(`/articles/${id}`, data),

  /** 更新状态 */
  updateStatus: (id: number, data: UpdateStatusReq) =>
    adminHttp.patch<any, Article>(`/articles/${id}/status`, data),

  /** 删除文章 */
  delete: (id: number) =>
    adminHttp.delete<any, null>(`/articles/${id}`),
}
