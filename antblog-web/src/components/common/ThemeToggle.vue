<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'
import { computed } from 'vue'
import { Moon, Sun } from 'lucide-vue-next'

const { theme, toggleDark } = useTheme()
const isDark = computed(() =>
  theme.value === 'dark' ||
  (theme.value === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches),
)
</script>

<template>
  <button
    type="button"
    class="inline-flex h-9 w-9 items-center justify-center rounded-md text-muted-foreground
           hover:bg-accent hover:text-foreground transition-colors"
    :aria-label="isDark ? '切换为亮色模式' : '切换为暗色模式'"
    @click="toggleDark"
  >
    <Sun v-if="isDark" class="h-4 w-4" />
    <Moon v-else class="h-4 w-4" />
  </button>
</template>
