import http, { adminHttp } from './http'
import type { PageResult } from '@/types/api.types'
import type { Comment, CreateCommentReq, ListCommentQuery, AdminListCommentQuery } from '@/types/comment.types'

export const commentApi = {
  /** 评论列表（前台，树形结构） */
  list: (params: ListCommentQuery) =>
    http.get<any, PageResult<Comment>>('/comments', { params }),

  /** 发表评论 */
  create: (data: CreateCommentReq) =>
    http.post<any, Comment>('/comments', data),

  /** 点赞评论 */
  like: (id: number) =>
    http.post<any, null>(`/comments/${id}/like`),

  /** 后台评论列表 */
  adminList: (params?: AdminListCommentQuery) =>
    adminHttp.get<any, PageResult<Comment>>('/comments', { params }),

  /** 后台审核评论 */
  adminReview: (id: number, status: number) =>
    adminHttp.patch<any, null>(`/comments/${id}/status`, { status }),

  /** 后台删除评论 */
  adminDelete: (id: number) =>
    adminHttp.delete<any, null>(`/comments/${id}`),
}
