<script setup lang="ts">
import { ref, onMounted } from 'vue'
import CommentItem from './CommentItem.vue'
import CommentForm from './CommentForm.vue'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import EmptyState from '@/components/common/EmptyState.vue'
import { commentApi } from '@/api/comment.api'
import type { Comment } from '@/types/comment.types'

const props = defineProps<{ articleId: number }>()
const comments = ref<Comment[]>([])
const total = ref(0)
const loading = ref(false)

async function fetchComments() {
  loading.value = true
  try {
    const res = await commentApi.list({ article_id: props.articleId, page: 1, page_size: 50 })
    comments.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(fetchComments)
</script>

<template>
  <section class="mt-12">
    <h2 class="font-serif text-xl font-semibold mb-6">
      评论
      <span v-if="total > 0" class="text-base font-normal text-muted-foreground ml-1">({{ total }})</span>
    </h2>

    <!-- Write comment -->
    <CommentForm :article-id="articleId" @submitted="fetchComments" />

    <Separator class="my-8" />

    <!-- Loading -->
    <div v-if="loading" class="space-y-6">
      <div v-for="i in 3" :key="i" class="flex gap-3">
        <Skeleton class="h-8 w-8 rounded-full shrink-0" />
        <div class="flex-1 space-y-2">
          <Skeleton class="h-4 w-24" />
          <Skeleton class="h-4 w-full" />
          <Skeleton class="h-4 w-3/4" />
        </div>
      </div>
    </div>

    <!-- Comments -->
    <div v-else-if="comments.length" class="space-y-6">
      <CommentItem
        v-for="comment in comments"
        :key="comment.id"
        :comment="comment"
        :article-id="articleId"
        @replied="fetchComments"
      />
    </div>

    <EmptyState v-else icon="💬" title="暂无评论" description="成为第一个留言的人吧！" />
  </section>
</template>
