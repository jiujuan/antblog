<script setup lang="ts">
const year = new Date().getFullYear()
</script>

<template>
  <footer class="border-t border-border/60 mt-16">
    <div class="container py-8 flex flex-col md:flex-row items-center justify-between gap-4 text-sm text-muted-foreground">
      <p>
        <span class="font-serif text-foreground font-medium">✦ AntBlog</span>
        &nbsp;— 记录思考，分享技术
      </p>
      <p>© {{ year }} · Built with Vue 3 &amp; ❤</p>
    </div>
  </footer>
</template>
