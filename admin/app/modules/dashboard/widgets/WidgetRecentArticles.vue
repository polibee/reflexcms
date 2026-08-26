<script setup lang="ts">
interface RecentArticle {
  id: number
  title: string
  slug: string
  status: string
  created_at: string
}

const articles = ref<RecentArticle[]>([])

onMounted(async () => {
  try {
    const stats = await $fetch<{ recentArticles: RecentArticle[] }>('/api/admin/stats')
    articles.value = stats.recentArticles ?? []
  } catch {
    articles.value = []
  }
})
</script>

<template>
  <div class="space-y-3">
    <p
      v-if="articles.length === 0"
      class="text-sm text-muted-foreground"
    >
      No articles yet.
    </p>
    <div
      v-for="article in articles"
      :key="article.id"
      class="flex items-center justify-between gap-4"
    >
      <span class="truncate text-sm">{{ article.title }}</span>
      <span class="shrink-0 text-xs text-muted-foreground">{{ article.status }}</span>
    </div>
  </div>
</template>
