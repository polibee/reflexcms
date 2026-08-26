<script setup lang="ts">
definePageMeta({ layout: 'public' })

interface PublicArticle {
  id: number
  title: string
  slug: string
  summary: string
  cover: string
  published_at: string
  view_count: number
  comment_count: number
}

interface PublicTopic {
  id: number
  title: string
  board_name?: string
  author_name?: string
  last_replier_name?: string
  reply_count: number
  view_count: number
  created_at: string
  last_reply_at?: string
}

const { config } = useSiteConfig()
const { get } = usePublicApi()

const articles = ref<PublicArticle[]>([])
const latestTopics = ref<PublicTopic[]>([])   // sort=new — 刚发布的帖子
const activeTopics = ref<PublicTopic[]>([])   // sort=reply — 最近有新回复的帖子
const loading = ref(true)

onMounted(async () => {
  if (!config.value) {
    try {
      config.value = await get<SiteConfig>('site/config')
    } catch { /* silent */ }
  }
  const mode = config.value?.mode ?? 'hybrid'
  const home = config.value?.home ?? 'both'
  const showArticles = mode !== 'forum' && home !== 'forum'
  const showTopics = mode !== 'blog' && home !== 'cms'

  try {
    const jobs: Promise<void>[] = []
    if (showArticles) {
      jobs.push(get<{ items: PublicArticle[] }>('articles?perPage=6').then(res => {
        articles.value = res.items ?? []
      }))
    }
    if (showTopics) {
      jobs.push(get<{ items: PublicTopic[] }>('topics?sort=new&perPage=8').then(res => {
        latestTopics.value = res.items ?? []
      }))
      jobs.push(get<{ items: PublicTopic[] }>('topics?sort=reply&perPage=8').then(res => {
        activeTopics.value = res.items ?? []
      }))
    }
    await Promise.all(jobs)
  } catch { /* silent */ }
  loading.value = false
})

function fmtDate(s?: string) {
  return s?.slice(0, 10) ?? ''
}
</script>

<template>
  <div>
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <div
      v-else
      class="space-y-10"
    >
      <!-- ===== CMS 文章区（cms 模式 / hybrid 上半块） ===== -->
      <section
        v-if="articles.length"
        class="space-y-4"
      >
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold text-gray-900">
            最新文章
          </h2>
          <NuxtLink
            to="/articles"
            class="text-sm text-blue-600 hover:underline"
          >更多 →</NuxtLink>
        </div>
        <div class="grid gap-5 md:grid-cols-2 lg:grid-cols-3">
          <NuxtLink
            v-for="a in articles"
            :key="a.id"
            :to="`/articles/${a.slug}`"
            class="group overflow-hidden rounded-xl border bg-white transition-shadow hover:shadow-md"
          >
            <div class="h-36 w-full overflow-hidden">
              <img
                v-if="a.cover"
                :src="a.cover"
                :alt="a.title"
                class="h-full w-full object-cover transition-transform group-hover:scale-105"
              >
              <div
                v-else
                class="flex h-full w-full items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 text-2xl font-black text-indigo-300"
              >
                {{ a.title.slice(0, 1).toUpperCase() }}
              </div>
            </div>
            <div class="p-4">
              <h3 class="line-clamp-2 font-semibold leading-snug text-gray-900 group-hover:text-blue-600">
                {{ a.title }}
              </h3>
              <p class="mt-1.5 line-clamp-2 text-sm text-gray-500">
                {{ a.summary }}
              </p>
              <div class="mt-3 flex items-center gap-3 text-xs text-gray-400">
                <span>{{ fmtDate(a.published_at) }}</span>
                <span>👁 {{ a.view_count }}</span>
                <span>💬 {{ a.comment_count }}</span>
              </div>
            </div>
          </NuxtLink>
        </div>
      </section>

      <!-- ===== 论坛区（forum 模式 / hybrid 下半块） ===== -->
      <section
        v-if="latestTopics.length || activeTopics.length"
        class="space-y-4"
      >
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold text-gray-900">
            社区讨论
          </h2>
          <NuxtLink
            to="/forums"
            class="text-sm text-blue-600 hover:underline"
          >进入论坛 →</NuxtLink>
        </div>

        <div class="grid gap-5 lg:grid-cols-2">
          <!-- 最新帖子 -->
          <div class="rounded-xl border bg-white">
            <h3 class="border-b px-5 py-3 text-sm font-semibold text-gray-900">
              🆕 最新帖子
            </h3>
            <div class="divide-y divide-gray-50">
              <NuxtLink
                v-for="t in latestTopics"
                :key="t.id"
                :to="`/topics/${t.id}`"
                class="block px-5 py-3 hover:bg-gray-50"
              >
                <div class="flex items-start justify-between gap-3">
                  <span class="min-w-0 flex-1 truncate text-sm font-medium text-gray-900">{{ t.title }}</span>
                  <span class="shrink-0 text-xs text-gray-400">💬 {{ t.reply_count }} / 👁 {{ t.view_count }}</span>
                </div>
                <p class="mt-1 text-xs text-gray-400">
                  {{ t.author_name ? `${t.author_name} 发布` : '发布于' }} {{ fmtDate(t.created_at) }}
                </p>
              </NuxtLink>
            </div>
            <p
              v-if="!latestTopics.length"
              class="px-5 py-8 text-center text-sm text-gray-400"
            >
              暂无帖子
            </p>
          </div>

          <!-- 最新回复（只到帖子层面，点进去看内容） -->
          <div class="rounded-xl border bg-white">
            <h3 class="border-b px-5 py-3 text-sm font-semibold text-gray-900">
              🔥 最新回复
            </h3>
            <div class="divide-y divide-gray-50">
              <NuxtLink
                v-for="t in activeTopics"
                :key="t.id"
                :to="`/topics/${t.id}`"
                class="block px-5 py-3 hover:bg-gray-50"
              >
                <div class="flex items-start justify-between gap-3">
                  <span class="min-w-0 flex-1 truncate text-sm font-medium text-gray-900">{{ t.title }}</span>
                  <span class="shrink-0 text-xs text-gray-400">💬 {{ t.reply_count }} / 👁 {{ t.view_count }}</span>
                </div>
                <p
                  v-if="t.last_replier_name"
                  class="mt-1 truncate text-xs text-gray-400"
                >
                  {{ t.last_replier_name }} 回复于 {{ fmtDate(t.last_reply_at) }}
                </p>
                <p
                  v-else
                  class="mt-1 text-xs text-gray-400"
                >
                  {{ t.author_name ? `${t.author_name} 发布` : '发布于' }} {{ fmtDate(t.created_at) }}
                </p>
              </NuxtLink>
            </div>
            <p
              v-if="!activeTopics.length"
              class="px-5 py-8 text-center text-sm text-gray-400"
            >
              暂无回复
            </p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
