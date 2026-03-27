import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth.api'
import type { UserInfo, LoginReq, RegisterReq } from '@/types/auth.types'
import { storage } from '@/utils/storage'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserInfo | null>(storage.get<UserInfo>('user_info'))
  const accessToken = ref<string>(storage.get<string>('access_token') ?? '')
  const refreshTokenVal = ref<string>(storage.get<string>('refresh_token') ?? '')

  const isLoggedIn = computed(() => !!accessToken.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 2)

  async function login(req: LoginReq) {
    const resp = await authApi.login(req)
    setAuth(resp.user, resp.access_token, resp.refresh_token)
    return resp
  }

  async function register(req: RegisterReq) {
    await authApi.register(req)
    return login({ email: req.email, password: req.password })
  }

  async function logout() {
    if (accessToken.value) {
      try { await authApi.logout() } catch {}
    }
    clearAuth()
  }

  async function refreshToken(): Promise<boolean> {
    if (!refreshTokenVal.value) return false
    try {
      const resp = await authApi.refresh(refreshTokenVal.value)
      accessToken.value = resp.access_token
      refreshTokenVal.value = resp.refresh_token
      storage.set('access_token', resp.access_token)
      storage.set('refresh_token', resp.refresh_token)
      return true
    } catch {
      return false
    }
  }

  async function fetchProfile() {
    const info = await authApi.getProfile()
    user.value = info
    storage.set('user_info', info)
  }

  function setAuth(userInfo: UserInfo, token: string, refresh: string) {
    user.value = userInfo
    accessToken.value = token
    refreshTokenVal.value = refresh
    storage.set('user_info', userInfo)
    storage.set('access_token', token)
    storage.set('refresh_token', refresh)
  }

  function clearAuth() {
    user.value = null
    accessToken.value = ''
    refreshTokenVal.value = ''
    storage.remove('user_info')
    storage.remove('access_token')
    storage.remove('refresh_token')
  }

  return {
    user,
    accessToken,
    isLoggedIn,
    isAdmin,
    login,
    register,
    logout,
    refreshToken,
    fetchProfile,
  }
})
