import axios, { type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api.types'
import { useAuthStore } from '@/stores/auth.store'
import router from '@/router'
import { storage } from '@/utils/storage'

function createClient(baseURL: string): AxiosInstance {
  const client: AxiosInstance = axios.create({
    baseURL,
    timeout: 15000,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  client.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const token = storage.get<string>('access_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    },
    (error) => Promise.reject(error),
  )

  client.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { code, msg, data } = response.data
    if (code === 0 || code === 200) {
      return data as any
    }
    const err = new Error(msg || '请求失败') as any
    err.code = code
    return Promise.reject(err)
  },
  async (error) => {
    const status = error.response?.status
    const url = String(error.config?.url || '')
    const isAuthEndpoint =
      url.includes('/auth/login') ||
      url.includes('/auth/register') ||
      url.includes('/auth/refresh') ||
      url.includes('/auth/logout')

    if (status === 401 && !isAuthEndpoint && !error.config?._retry) {
      error.config._retry = true
      const authStore = useAuthStore()
      const refreshed = await authStore.refreshToken()
      if (!refreshed) {
        await authStore.logout()
        router.push({ name: 'login', query: { redirect: router.currentRoute.value.fullPath } })
      } else {
        const config = error.config
        const token = storage.get<string>('access_token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return client.request(config)
      }
    }
    const msg = error.response?.data?.msg || error.message || '网络错误'
    const err = new Error(msg) as Error & { code?: number }
    err.code = error.response?.data?.code
    return Promise.reject(err)
  },
  )

  return client
}

const http = createClient('/api/v1')
export const adminHttp = createClient('/api/admin')

export default http
