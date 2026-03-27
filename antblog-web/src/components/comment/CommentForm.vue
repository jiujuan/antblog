<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Input } from '@/components/ui/input'
import { useAuthStore } from '@/stores/auth.store'
import { commentApi } from '@/api/comment.api'
import { useToast } from '@/components/ui/toast'
import type { CreateCommentReq } from '@/types/comment.types'

const props = defineProps<{
  articleId: number
  parentId?: number
  replyToId?: number
  placeholder?: string
}>()

const emit = defineEmits<{ submitted: [] }>()

const authStore = useAuthStore()
const { toast } = useToast()
const content = ref('')
const guestName = ref('')
const guestEmail = ref('')
const loading = ref(false)

async function submit() {
  if (!content.value.trim()) return
  if (!authStore.isLoggedIn && !guestName.value.trim()) {
    toast({ title: '发表失败', description: '游客昵称不能为空', variant: 'destructive' })
    return
  }
  loading.value = true
  try {
    const req: CreateCommentReq = {
      article_id: props.articleId,
      parent_id: props.parentId,
      reply_to_id: props.replyToId,
      content: content.value,
    }
    if (!authStore.isLoggedIn) {
      req.nickname = guestName.value.trim()
      req.email = guestEmail.value.trim()
    }
    await commentApi.create(req)
    content.value = ''
    toast({ title: '评论已提交', description: '待审核后展示', variant: 'success' })
    emit('submitted')
  } catch (e: any) {
    toast({ title: '发表失败', description: e.message, variant: 'destructive' })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="space-y-3">
    <!-- Guest info (not logged in) -->
    <div v-if="!authStore.isLoggedIn" class="grid grid-cols-2 gap-2">
      <Input v-model="guestName" placeholder="昵称（必填）" />
      <Input v-model="guestEmail" type="email" placeholder="邮箱（可选）" />
    </div>

    <Textarea
      v-model="content"
      :placeholder="placeholder ?? '写下你的想法…'"
      :rows="3"
    />

    <div class="flex justify-end">
      <Button :disabled="loading || !content.trim()" @click="submit">
        {{ loading ? '提交中…' : '发表评论' }}
      </Button>
    </div>
  </div>
</template>
