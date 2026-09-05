<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Forum index: topic-centric feed (boards live in the sidebar widgets).
 * Two tabs — newest topics and topics with the newest replies — both
 * paginated, with tab/page mirrored into the URL query for shareable links.
 * Reply activity lists show WHICH topic was replied to (title + last
 * replier); readers click through to see the actual content. */

interface TopicRow {
  id: number
  title: string
  board_name?: string
  author_name?: string
  last_replier_name?: string
  is_pinned: boolean
  is_featured: boolean
  reply_count: number
  view_count: number
  created_at: string
  last_reply_at?: string
}

const route = useRoute()
const router = useRouter()
const { t } = useLocale()

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

type Tab = 'topics' | 'replies'
const tab = ref<Tab>(route.query.tab === 'replies' ? 'replies' : 'topics')
const page = ref(1)
const perPage = 20

const topics = ref<TopicRow[]>([])
const total = ref(0)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / perPage)))
const loading = ref(true)

async function fetchList() {
  loading.value = true
  try {
    const sort = tab.value === 'replies' ? 'reply' : 'new'
    const res = await rawFetch(`/api/v1/topics?sort=${sort}&page=${page.value}&perPage=${perPage}`) as { items: TopicRow[], total: number }
    topics.value = res.items ?? []
    total.value = res.total ?? 0
  } catch { /* silent */ }
  loading.value = false
}

function switchTab(t: Tab) {
  if (tab.value === t) return
  tab.value = t
  page.value = 1
  router.replace({ query: t === 'replies' ? { tab: 'replies' } : {} })
  fetchList()
}

function goToPage(p: number) {
  page.value = p
  router.replace({ query: { ...(tab.value === 'replies' ? { tab: 'replies' } : {}), page: String(p) } })
  fetchList()
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

function fmtDate(s?: string) {
  return s?.slice(0, 10) ?? ''
}

onMounted(fetchList)

useHead({ title: '论坛' })
</script>

<template>
  <div>
    <div class="mb-5 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">
        {{ t('forum.title') }}
      </h1>
      <div class="flex rounded-lg border border-gray-200 bg-white p-1">
        <button
          type="button"
          class="rounded-md px-4 py-1.5 text-sm font-medium"
          :class="tab === 'topics'
            ? 'bg-blue-600 text-white'
            : 'text-gray-600 hover:bg-gray-50'"
          @click="switchTab('topics')"
        >
          {{ t('forum.latest') }}
        </button>
        <button
          type="button"
          class="rounded-md px-4 py-1.5 text-sm font-medium"
          :class="tab === 'replies'
            ? 'bg-blue-600 text-white'
            : 'text-gray-600 hover:bg-gray-50'"
          @click="switchTab('replies')"
        >
          {{ t('forum.latest_replies') }}
        </button>
      </div>
    </div>

    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <template v-else>
      <div class="divide-y divide-gray-100 rounded-xl border bg-white">
        <NuxtLink
          v-for="t in topics"
          :key="t.id"
          :to="`/topics/${t.id}`"
          class="flex items-center gap-3 px-5 py-3.5 hover:bg-gray-50"
        >
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <span
                v-if="t.is_pinned"
                class="shrink-0 rounded bg-red-100 px-1.5 py-0.5 text-xs text-red-600"
              >置顶</span>
              <span
                v-if="t.is_featured"
                class="shrink-0 rounded bg-blue-100 px-1.5 py-0.5 text-xs text-blue-600"
              >精华</span>
              <span
                v-if="t.board_name"
                class="shrink-0 rounded bg-gray-100 px-1.5 py-0.5 text-xs text-gray-500"
              >{{ t.board_name }}</span>
              <h2 class="truncate font-medium text-gray-900">{{ t.title }}</h2>
            </div>
            <!-- 最新回复 tab 显示最后回复人；最新帖子 tab 显示发帖人 -->
            <p
              v-if="tab === 'replies' && t.last_replier_name"
              class="mt-1 truncate text-xs text-gray-400"
            >
              {{ t.last_replier_name }} 回复于 {{ fmtDate(t.last_reply_at) }}
            </p>
            <p
              v-else
              class="mt-1 truncate text-xs text-gray-400"
            >
              {{ t.author_name ? `${t.author_name} 发布` : '发布于' }} {{ fmtDate(t.created_at) }}
            </p>
          </div>
          <div class="flex shrink-0 items-center gap-4 text-xs text-gray-400">
            <span>💬 {{ t.reply_count }}</span>
            <span>👁 {{ t.view_count }}</span>
          </div>
        </NuxtLink>

        <p
          v-if="!topics.length"
          class="py-16 text-center text-gray-400"
        >
          {{ tab === 'topics' ? t('forum.empty') : t('forum.empty_replies') }}
        </p>
      </div>

      <PublicPagination
        :page="page"
        :total-pages="totalPages"
        @change="goToPage"
      />
    </template>
  </div>
</template>
