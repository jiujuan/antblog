<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/composables/useAuth'
import { useToast } from '@/components/ui/toast'

const { login } = useAuth()
const { toast } = useToast()
const email = ref('')
const password = ref('')
const loading = ref(false)

function resolveLoginErrorMessage(err: any) {
  const code = err?.code
  if (code === 20001) return '用户不存在，请检查用户名/邮箱是否正确'
  if (code === 20003) return '密码错误，请重新输入'
  if (code === 10001) return '参数格式错误，请检查输入后重试'
  return err?.message || '登录失败，请稍后重试'
}

async function handleSubmit() {
  loading.value = true
  try {
    await login({ email: email.value, password: password.value })
  } catch (e: any) {
    toast({
      title: '登录失败',
      description: resolveLoginErrorMessage(e),
      variant: 'destructive',
    })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-3.5rem)] flex items-center justify-center px-4">
    <div class="w-full max-w-sm space-y-6 animate-fade-up">
      <!-- Header -->
      <div class="text-center">
        <p class="font-serif text-3xl font-semibold">✦ AntBlog</p>
        <p class="text-muted-foreground mt-1 text-sm">欢迎回来，请登录你的账号</p>
      </div>

      <!-- Form -->
      <div class="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
        <div class="space-y-2">
          <label class="text-sm font-medium">用户名或邮箱</label>
          <Input v-model="email" type="text" placeholder="username / you@example.com" />
        </div>
        <div class="space-y-2">
          <label class="text-sm font-medium">密码</label>
          <Input v-model="password" type="password" placeholder="••••••••" @keyup.enter="handleSubmit" />
        </div>
        <Button class="w-full" :disabled="loading" @click="handleSubmit">
          {{ loading ? '登录中…' : '登录' }}
        </Button>
      </div>

      <p class="text-center text-sm text-muted-foreground">
        还没有账号？
        <RouterLink to="/register" class="text-primary hover:underline font-medium">立即注册</RouterLink>
      </p>
    </div>
  </div>
</template>
