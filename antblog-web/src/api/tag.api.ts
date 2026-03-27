import http, { adminHttp } from './http'
import type { Tag, CreateTagReq, UpdateTagReq, BatchCreateTagReq } from '@/types/tag.types'

export const tagApi = {
  /** 标签列表（前台公开） */
  list: () =>
    http.get<any, Tag[]>('/tags'),

  /** 按 slug 获取标签 */
  getBySlug: (slug: string) =>
    http.get<any, Tag>(`/tags/${slug}`),

  /** 创建标签（后台） */
  create: (data: CreateTagReq) =>
    adminHttp.post<any, Tag>('/tags', data),

  /** 批量创建标签（后台） */
  batchCreate: (data: BatchCreateTagReq) =>
    adminHttp.post<any, Tag[]>('/tags/batch', data),

  /** 更新标签（后台） */
  update: (id: number, data: UpdateTagReq) =>
    adminHttp.put<any, Tag>(`/tags/${id}`, data),

  /** 删除标签（后台） */
  delete: (id: number) =>
    adminHttp.delete<any, null>(`/tags/${id}`),
}
