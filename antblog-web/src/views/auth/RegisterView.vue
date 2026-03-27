<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/composables/useAuth'
import { useToast } from '@/components/ui/toast'

const { register } = useAuth()
const { toast } = useToast()
const username = ref('')
const email = ref('')
const password = ref('')
const nickname = ref('')
const loading = ref(false)

async function handleSubmit() {
  loading.value = true
  try {
    await register({ username: username.value, email: email.value, password: password.value, nickname: nickname.value })
  } catch (e: any) {
    toast({ title: '注册失败', description: e.message, variant: 'destructive' })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-3.5rem)] flex items-center justify-center px-4">
    <div class="w-full max-w-sm space-y-6 animate-fade-up">
      <div class="text-center">
        <p class="font-serif text-3xl font-semibold">✦ AntBlog</p>
        <p class="text-muted-foreground mt-1 text-sm">创建你的账号</p>
      </div>

      <div class="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
        <div class="space-y-2">
          <label class="text-sm font-medium">用户名</label>
          <Input v-model="username" placeholder="username" />
        </div>
        <div class="space-y-2">
          <label class="text-sm font-medium">昵称（可选）</label>
          <Input v-model="nickname" placeholder="显示名称" />
        </div>
        <div class="space-y-2">
          <label class="text-sm font-medium">邮箱</label>
          <Input v-model="email" type="email" placeholder="you@example.com" />
        </div>
        <div class="space-y-2">
          <label class="text-sm font-medium">密码</label>
          <Input v-model="password" type="password" placeholder="至少 8 位" />
        </div>
        <Button class="w-full" :disabled="loading" @click="handleSubmit">
          {{ loading ? '注册中…' : '注册' }}
        </Button>
      </div>

      <p class="text-center text-sm text-muted-foreground">
        已有账号？
        <RouterLink to="/login" class="text-primary hover:underline font-medium">立即登录</RouterLink>
      </p>
    </div>
  </div>
</template>
