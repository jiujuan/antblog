<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Folder, Tag } from 'lucide-vue-next'
import { Skeleton } from '@/components/ui/skeleton'
import { useCategoryStore } from '@/stores/category.store'
import { useTagStore } from '@/stores/tag.store'

const catStore = useCategoryStore()
const tagStore = useTagStore()

onMounted(async () => {
  await Promise.all([catStore.fetchAll(), tagStore.fetchAll()])
})
</script>

<template>
  <aside class="space-y-8">
    <!-- Categories -->
    <section>
      <h3 class="flex items-center gap-2 font-serif text-sm font-semibold uppercase tracking-widest text-muted-foreground mb-3">
        <Folder class="h-3.5 w-3.5" /> 分类
      </h3>
      <ul v-if="catStore.loaded" class="space-y-1.5">
        <li v-for="cat in catStore.categories" :key="cat.id">
          <RouterLink
            :to="`/categories/${cat.slug}`"
            class="flex items-center justify-between text-sm text-muted-foreground hover:text-foreground transition-colors group"
          >
            <span class="group-hover:translate-x-0.5 transition-transform">{{ cat.name }}</span>
            <span class="text-xs bg-muted rounded-full px-2 py-0.5">{{ cat.article_count }}</span>
          </RouterLink>
        </li>
      </ul>
      <div v-else class="space-y-2">
        <Skeleton v-for="i in 4" :key="i" class="h-5 w-full" />
      </div>
    </section>

    <!-- Tags -->
    <section>
      <h3 class="flex items-center gap-2 font-serif text-sm font-semibold uppercase tracking-widest text-muted-foreground mb-3">
        <Tag class="h-3.5 w-3.5" /> 标签
      </h3>
      <div v-if="tagStore.loaded" class="flex flex-wrap gap-2">
        <RouterLink
          v-for="tag in tagStore.tags"
          :key="tag.id"
          :to="`/tags/${tag.slug}`"
          class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium border transition-colors hover:border-primary/50 hover:text-primary"
          :style="{ borderColor: tag.color + '66', color: tag.color }"
        >
          {{ tag.name }}
          <span class="text-[10px] opacity-80">({{ tag.article_count ?? 0 }})</span>
        </RouterLink>
      </div>
      <div v-else class="flex flex-wrap gap-2">
        <Skeleton v-for="i in 8" :key="i" class="h-6 w-16 rounded-full" />
      </div>
    </section>
  </aside>
</template>
