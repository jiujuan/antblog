<script setup lang="ts">
import { ref } from 'vue'
import { Heart, Reply } from 'lucide-vue-next'
import { Avatar } from '@/components/ui/avatar'
import CommentForm from './CommentForm.vue'
import { commentApi } from '@/api/comment.api'
import { fromNow } from '@/utils/format'
import type { Comment } from '@/types/comment.types'

const props = defineProps<{ comment: Comment; articleId: number; depth?: number }>()
const emit = defineEmits<{ replied: [] }>()

const showReply = ref(false)
const liked = ref(false)
const likeCount = ref(props.comment.like_count)

async function toggleLike() {
  if (liked.value) return
  liked.value = true
  likeCount.value++
  try { await commentApi.like(props.comment.id) } catch { liked.value = false; likeCount.value-- }
}
</script>

<template>
  <div class="flex gap-3">
    <Avatar :fallback="comment.nickname || 'U'" size="sm" class="shrink-0 mt-0.5" />

    <div class="flex-1 min-w-0">
      <div class="flex items-baseline gap-2 mb-1">
        <span class="font-medium text-sm text-foreground">
          {{ comment.nickname || '匿名' }}
        </span>
        <span class="text-xs text-muted-foreground">{{ fromNow(comment.created_at) }}</span>
      </div>

      <p class="text-sm text-foreground leading-relaxed whitespace-pre-wrap">{{ comment.content }}</p>

      <div class="flex items-center gap-3 mt-2">
        <button
          type="button"
          class="flex items-center gap-1 text-xs text-muted-foreground hover:text-rose-500 transition-colors"
          @click="toggleLike"
        >
          <Heart class="h-3 w-3" :class="liked ? 'fill-rose-500 text-rose-500' : ''" />
          <span>{{ likeCount }}</span>
        </button>
        <button
          v-if="(depth ?? 0) < 1"
          type="button"
          class="flex items-center gap-1 text-xs text-muted-foreground hover:text-primary transition-colors"
          @click="showReply = !showReply"
        >
          <Reply class="h-3 w-3" />
          回复
        </button>
      </div>

      <!-- Reply form -->
      <div v-if="showReply" class="mt-3">
        <CommentForm
          :article-id="articleId"
          :parent-id="comment.id"
            :reply-to-id="comment.id"
            :placeholder="`回复 ${comment.nickname || '匿名'}…`"
          @submitted="showReply = false; emit('replied')"
        />
      </div>

      <!-- Nested children -->
      <div v-if="comment.children?.length" class="mt-4 pl-3 border-l border-border space-y-4">
        <CommentItem
          v-for="child in comment.children"
          :key="child.id"
          :comment="child"
          :article-id="articleId"
          :depth="(depth ?? 0) + 1"
          @replied="emit('replied')"
        />
      </div>
    </div>
  </div>
</template>
