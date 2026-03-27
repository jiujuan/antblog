import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Article, ArticleListItem } from '@/types/article.types'

export const useArticleStore = defineStore('article', () => {
  /** 当前查看的文章详情缓存（key=slug） */
  const detailCache = ref<Map<string, Article>>(new Map())

  /** 最近一次列表查询结果 */
  const list = ref<ArticleListItem[]>([])
  const total = ref(0)

  function cacheDetail(slug: string, article: Article) {
    detailCache.value.set(slug, article)
  }

  function getCachedDetail(slug: string): Article | undefined {
    return detailCache.value.get(slug)
  }

  function setList(items: ArticleListItem[], t: number) {
    list.value = items
    total.value = t
  }

  /** 更新列表中某文章的点赞/收藏状态（乐观更新） */
  function updateInteraction(id: number, field: 'liked' | 'bookmarked', value: boolean) {
    const item = list.value.find((a) => a.id === id)
    if (item) {
      item[field] = value
      if (field === 'liked') item.like_count += value ? 1 : -1
      if (field === 'bookmarked') item.bookmark_count += value ? 1 : -1
    }
    // Also update detail cache
    detailCache.value.forEach((article) => {
      if (article.id === id) {
        article[field] = value
        if (field === 'liked') article.like_count += value ? 1 : -1
        if (field === 'bookmarked') article.bookmark_count += value ? 1 : -1
      }
    })
  }

  return { detailCache, list, total, cacheDetail, getCachedDetail, setList, updateInteraction }
})
