import { defineStore } from 'pinia'
import { ref } from 'vue'
import { categoryApi } from '@/api/category.api'
import type { Category } from '@/types/category.types'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([])
  const loaded = ref(false)

  async function fetchAll() {
    if (loaded.value) return
    categories.value = await categoryApi.list()
    loaded.value = true
  }

  function invalidate() {
    loaded.value = false
  }

  return { categories, loaded, fetchAll, invalidate }
})
