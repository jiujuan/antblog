<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Search, Menu, X, PenSquare } from 'lucide-vue-next'
import ThemeToggle from './ThemeToggle.vue'
import { Avatar } from '@/components/ui/avatar'
import { useAuthStore } from '@/stores/auth.store'

const authStore = useAuthStore()
const router = useRouter()
const menuOpen = ref(false)
const searchQuery = ref('')
const isLoggedIn = computed(() => authStore.isLoggedIn)
const isAdmin    = computed(() => authStore.isAdmin)

function submitSearch() {
  if (searchQuery.value.trim()) {
    router.push({ name: 'search', query: { q: searchQuery.value } })
    searchQuery.value = ''
  }
}

async function handleLogout() {
  await authStore.logout()
  router.push('/')
}

const navLinks = [
  { name: '首页',   to: '/' },
  { name: '归档',   to: '/archive' },
]
</script>

<template>
  <header class="sticky top-0 z-40 w-full border-b border-border/60 bg-background/90 backdrop-blur-sm">
    <div class="container flex h-14 items-center justify-between gap-4">
      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-2 font-serif text-xl font-semibold tracking-tight shrink-0">
        <span class="text-primary">✦</span>
        <span>AntBlog</span>
      </RouterLink>

      <!-- Desktop nav -->
      <nav class="hidden md:flex items-center gap-6 text-sm font-medium">
        <RouterLink
          v-for="link in navLinks"
          :key="link.to"
          :to="link.to"
          class="text-muted-foreground hover:text-foreground transition-colors"
          active-class="!text-foreground"
        >
          {{ link.name }}
        </RouterLink>
      </nav>

      <!-- Right actions -->
      <div class="flex items-center gap-1">
        <!-- Search -->
        <form class="hidden md:flex items-center" @submit.prevent="submitSearch">
          <div class="relative">
            <Search class="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
            <input
              v-model="searchQuery"
              type="search"
              placeholder="搜索..."
              class="h-8 w-44 rounded-md border border-input bg-background pl-8 pr-3 text-sm
                     focus:outline-none focus:ring-2 focus:ring-ring focus:w-56 transition-all"
            />
          </div>
        </form>

        <ThemeToggle />

        <!-- Admin link -->
        <RouterLink
          v-if="isAdmin"
          to="/admin"
          class="hidden md:inline-flex items-center gap-1.5 rounded-md px-3 h-9 text-sm
                 bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          <PenSquare class="h-3.5 w-3.5" />
          后台
        </RouterLink>

        <!-- User -->
        <template v-if="isLoggedIn">
          <RouterLink to="/admin/dashboard" v-if="isAdmin">
            <Avatar :fallback="authStore.user?.nickname?.charAt(0)" size="sm" class="cursor-pointer" />
          </RouterLink>
          <button v-else type="button" class="text-sm text-muted-foreground hover:text-foreground" @click="handleLogout">
            登出
          </button>
        </template>
        <template v-else>
          <RouterLink
            to="/login"
            class="hidden md:inline-flex text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            登录
          </RouterLink>
        </template>

        <!-- Mobile menu toggle -->
        <button
          type="button"
          class="md:hidden inline-flex h-9 w-9 items-center justify-center rounded-md hover:bg-accent transition-colors"
          @click="menuOpen = !menuOpen"
        >
          <Menu v-if="!menuOpen" class="h-4 w-4" />
          <X v-else class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- Mobile nav -->
    <Transition name="slide-down">
      <div v-if="menuOpen" class="md:hidden border-t border-border bg-background px-4 pb-4 pt-2">
        <nav class="flex flex-col gap-2 text-sm">
          <RouterLink
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="py-2 text-muted-foreground hover:text-foreground transition-colors"
            @click="menuOpen = false"
          >
            {{ link.name }}
          </RouterLink>
          <RouterLink v-if="!isLoggedIn" to="/login" class="py-2 text-muted-foreground" @click="menuOpen = false">
            登录
          </RouterLink>
          <button v-else type="button" class="py-2 text-left text-muted-foreground" @click="handleLogout">
            登出
          </button>
        </nav>
        <!-- Mobile search -->
        <form class="mt-3" @submit.prevent="submitSearch">
          <input
            v-model="searchQuery"
            type="search"
            placeholder="搜索文章..."
            class="w-full h-9 rounded-md border border-input bg-muted px-3 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
          />
        </form>
      </div>
    </Transition>
  </header>
</template>

<style scoped>
.slide-down-enter-active, .slide-down-leave-active { transition: all 0.2s ease; }
.slide-down-enter-from, .slide-down-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
