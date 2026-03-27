import { defineStore } from 'pinia'
import { ref } from 'vue'
import { tagApi } from '@/api/tag.api'
import type { Tag } from '@/types/tag.types'

export const useTagStore = defineStore('tag', () => {
  const tags = ref<Tag[]>([])
  const loaded = ref(false)

  async function fetchAll() {
    if (loaded.value) return
    tags.value = await tagApi.list()
    loaded.value = true
  }

  function invalidate() {
    loaded.value = false
  }

  return { tags, loaded, fetchAll, invalidate }
})
