<script setup lang="ts">
defineProps<{ title?: string; description?: string; icon?: string }>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-20 text-center">
    <div class="mb-4 text-5xl opacity-30">{{ icon ?? '📭' }}</div>
    <h3 class="font-serif text-lg font-medium text-foreground mb-1">
      {{ title ?? '暂无内容' }}
    </h3>
    <p class="text-sm text-muted-foreground max-w-xs">
      {{ description ?? '这里还什么都没有，稍后再来看看吧。' }}
    </p>
    <slot />
  </div>
</template>
