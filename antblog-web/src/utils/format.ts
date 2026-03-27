/**
 * 日期/数字格式化工具
 */

/** 格式化为 YYYY-MM-DD */
export function formatDate(dateStr: string | null | undefined): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

/** 格式化为相对时间（如：3 天前） */
export function fromNow(dateStr: string | null | undefined): string {
  if (!dateStr) return ''
  const now = Date.now()
  const diff = now - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins} 分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} 小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days} 天前`
  const months = Math.floor(days / 30)
  if (months < 12) return `${months} 个月前`
  return `${Math.floor(months / 12)} 年前`
}

/** 数字缩略（如：1200 → 1.2k） */
export function formatCount(n: number): string {
  if (n >= 10000) return (n / 10000).toFixed(1).replace(/\.0$/, '') + 'w'
  if (n >= 1000) return (n / 1000).toFixed(1).replace(/\.0$/, '') + 'k'
  return String(n)
}

/** 字数转阅读时间（分钟） */
export function readingTime(wordCount: number): number {
  return Math.max(1, Math.round(wordCount / 300))
}

/** 月份名 */
export const MONTHS = ['一月','二月','三月','四月','五月','六月',
                       '七月','八月','九月','十月','十一月','十二月']
