import http from './http'
import type { LoginReq, LoginResp, RegisterReq, UpdateProfileReq, ChangePasswordReq, UserInfo } from '@/types/auth.types'

export const authApi = {
  /** 注册 */
  register: (data: RegisterReq) =>
    http.post<any, UserInfo>('/auth/register', data),

  /** 登录 */
  login: (data: LoginReq) =>
    http.post<any, LoginResp>('/auth/login', data),

  /** 刷新 Token */
  refresh: (refreshToken: string) =>
    http.post<any, { access_token: string; refresh_token: string; expires_at: string; user: UserInfo }>('/auth/refresh', { refresh_token: refreshToken }),

  /** 登出 */
  logout: () =>
    http.post<any, null>('/auth/logout'),

  /** 获取当前用户资料 */
  getProfile: () =>
    http.get<any, UserInfo>('/user/profile'),

  /** 更新资料 */
  updateProfile: (data: UpdateProfileReq) =>
    http.put<any, UserInfo>('/user/profile', data),

  /** 修改密码 */
  changePassword: (data: ChangePasswordReq) =>
    http.put<any, null>('/user/password', data),
}
