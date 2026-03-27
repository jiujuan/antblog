import http, { adminHttp } from './http'
import type { PageResult } from '@/types/api.types'
import type { Media, AdminListMediaQuery, BindArticleReq } from '@/types/media.types'

export const mediaApi = {
  /** 我的媒体库 */
  myMedia: (params?: { page?: number; page_size?: number }) =>
    http.get<any, PageResult<Media>>('/user/media', { params }),

  /** 后台：上传文件 */
  upload: (file: File, articleId?: number) => {
    const form = new FormData()
    form.append('file', file)
    if (articleId != null) form.append('article_id', String(articleId))
    return adminHttp.post<any, Media>('/media/upload', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  /** 后台：媒体列表 */
  adminList: (params?: AdminListMediaQuery) =>
    adminHttp.get<any, PageResult<Media>>('/media', { params }),

  /** 后台：媒体详情 */
  adminDetail: (id: number) =>
    adminHttp.get<any, Media>(`/media/${id}`),

  /** 后台：绑定/解绑文章 */
  bindArticle: (id: number, data: BindArticleReq) =>
    adminHttp.patch<any, Media>(`/media/${id}/bind`, data),

  /** 后台：删除媒体 */
  delete: (id: number) =>
    adminHttp.delete<any, null>(`/media/${id}`),
}
