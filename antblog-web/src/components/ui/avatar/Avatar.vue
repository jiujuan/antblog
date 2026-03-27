<script setup lang="ts">
import { cn } from '@/utils/cn'
interface Props { src?: string; alt?: string; fallback?: string; class?: string; size?: 'sm' | 'md' | 'lg' }
const props = withDefaults(defineProps<Props>(), { size: 'md', alt: '' })
const sizes = { sm: 'h-7 w-7 text-xs', md: 'h-9 w-9 text-sm', lg: 'h-12 w-12 text-base' }
</script>
<template>
  <div :class="cn('relative flex shrink-0 overflow-hidden rounded-full bg-muted', sizes[size ?? 'md'], props.class)">
    <img v-if="src" :src="src" :alt="alt" class="aspect-square h-full w-full object-cover" />
    <span v-else class="flex h-full w-full items-center justify-center font-medium uppercase text-muted-foreground">
      {{ fallback ?? alt?.charAt(0) ?? '?' }}
    </span>
  </div>
</template>
