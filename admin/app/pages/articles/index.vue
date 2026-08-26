<script setup lang="ts">
definePageMeta({ layout: 'public' })

interface Article {
  id: number
  title: string
  slug: string
  summary: string
  cover: string
  published_at: string
  view_count: number
}

const route = useRoute()
const { get } = usePublicApi()

const articles = ref<Article[]>([])
const total = ref(0)
const page = ref(Number(route.query.page) || 1)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / 20)))
const loading = ref(true)

async function fetchArticles() {
  loading.value = true
  try {
    const res = await get<{ items: Article[], total: number }>(`articles?page=${page.value}&perPage=20`)
    articles.value = res.items ?? []
    total.value = res.total ?? 0
  } catch { /* silent */ }
  loading.value = false
}

watch(page, fetchArticles, { immediate: true })
</script>

<template>
  <div>
    <h1 class="mb-6 text-2xl font-bold">
      文章
    </h1>
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>
    <div
      v-else
      class="grid gap-6 md:grid-cols-2"
    >
      <NuxtLink
        v-for="a in articles"
        :key="a.id"
        :to="`/articles/${a.slug}`"
        class="group overflow-hidden rounded-lg border bg-white hover:shadow-md transition-shadow"
      >
        <img
          v-if="a.cover"
          :src="a.cover"
          :alt="a.title"
          class="h-44 w-full object-cover"
        >
        <div
          v-else
          class="flex h-44 w-full items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100"
        >
          <span class="text-3xl font-bold text-blue-200">{{ a.title.charAt(0) }}</span>
        </div>
        <div class="p-5">
          <h2 class="font-semibold text-gray-900 group-hover:text-blue-600">{{ a.title }}</h2>
          <p class="mt-1 line-clamp-2 text-sm text-gray-500">{{ a.summary }}</p>
          <div class="mt-3 flex items-center gap-4 text-xs text-gray-400">
            <time>{{ a.published_at?.slice(0, 10) }}</time>
            <span>{{ a.view_count }} 阅读</span>
          </div>
        </div>
      </NuxtLink>
    </div>
    <div
      v-if="totalPages > 1"
      class="mt-8 flex justify-center gap-2"
    >
      <button
        :disabled="page <= 1"
        class="rounded border px-3 py-1 text-sm disabled:opacity-40"
        @click="page--"
      >
        上一页
      </button>
      <span class="px-3 py-1 text-sm">{{ page }} / {{ totalPages }}</span>
      <button
        :disabled="page >= totalPages"
        class="rounded border px-3 py-1 text-sm disabled:opacity-40"
        @click="page++"
      >
        下一页
      </button>
    </div>
  </div>
</template>
