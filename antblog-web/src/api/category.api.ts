import http, { adminHttp } from './http'
import type { Category, CreateCategoryReq, UpdateCategoryReq } from '@/types/category.types'

export const categoryApi = {
  /** 分类列表（前台公开） */
  list: () =>
    http.get<any, Category[]>('/categories'),

  /** 按 slug 获取分类 */
  getBySlug: (slug: string) =>
    http.get<any, Category>(`/categories/${slug}`),

  /** 创建分类（后台） */
  create: (data: CreateCategoryReq) =>
    adminHttp.post<any, Category>('/categories', data),

  /** 更新分类（后台） */
  update: (id: number, data: UpdateCategoryReq) =>
    adminHttp.put<any, Category>(`/categories/${id}`, data),

  /** 删除分类（后台） */
  delete: (id: number) =>
    adminHttp.delete<any, null>(`/categories/${id}`),
}
