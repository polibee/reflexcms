<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* My notification feed: mentions, moderation results, reply notices.
 * Rows link to their source; 全部已读 clears unread badges. */

interface NotificationRow {
  id: number
  type: string
  data: string
  read_at: string | null
  created_at: string
}

const { user: me } = useAuthUser()
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

const items = ref<NotificationRow[]>([])
const page = ref(1)
const totalPages = ref(1)
const loading = ref(true)

const TYPE_LABELS: Record<string, string> = {
  mention: '提及了我',
  'forum.reply': '帖子有了新回复',
  'comment.approved': '评论已通过',
  'comment.rejected': '评论未通过',
}

function parseData(row: NotificationRow): Record<string, any> {
  try {
    return JSON.parse(row.data || '{}')
  } catch {
    return {}
  }
}

function linkFor(row: NotificationRow): string {
  const d = parseData(row)
  if (row.type === 'mention' || row.type === 'forum.reply') {
    return d.topic_id ? `/topics/${d.topic_id}` : d.link_id ? `/topics/${d.link_id}` : '/forums'
  }
  return d.comment_id ? `/articles` : '#'
}

function titleFor(row: NotificationRow): string {
  const d = parseData(row)
  if (row.type === 'mention') return `${d.from_name ?? '有人'} 在回复中提到了你`
  return TYPE_LABELS[row.type] ?? row.type
}

async function fetchList() {
  loading.value = true
  try {
    const res = await rawFetch(`/api/v1/me/notifications?page=${page.value}&perPage=15`) as { items: NotificationRow[], totalPages: number }
    items.value = res.items ?? []
    totalPages.value = res.totalPages ?? 1
  } catch { /* silent */ }
  loading.value = false
}

async function readAll() {
  try {
    await rawFetch('/api/v1/me/notifications/read-all', { method: 'POST' })
    await fetchList()
  } catch { /* silent */ }
}

function goToPage(p: number) {
  page.value = p
  void fetchList()
}

onMounted(async () => {
  // unconditional: anonymous callers just get a 401 that we swallow
  await fetchList()
})

useHead({ title: '通知中心' })
</script>

<template>
  <div class="mx-auto max-w-3xl">
    <div class="mb-5 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">
        通知中心
      </h1>
      <button
        v-if="me"
        type="button"
        class="rounded-md border border-gray-300 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50"
        @click="readAll"
      >
        全部已读
      </button>
    </div>

    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <div
      v-else-if="me"
      class="divide-y divide-gray-100 rounded-xl border bg-white"
    >
      <NuxtLink
        v-for="n in items"
        :key="n.id"
        :to="linkFor(n)"
        class="block px-5 py-3.5 hover:bg-gray-50"
      >
        <div class="flex items-center justify-between gap-3">
          <span
            class="truncate text-sm font-medium"
            :class="n.read_at ? 'text-gray-600' : 'text-gray-900'"
          >{{ titleFor(n) }}</span>
          <span class="shrink-0 text-xs text-gray-400">{{ n.created_at?.slice(0, 10) }}</span>
        </div>
        <p
          v-if="parseData(n).excerpt"
          class="mt-1 truncate text-xs text-gray-500"
        >
          {{ parseData(n).excerpt }}
        </p>
      </NuxtLink>

      <p
        v-if="!items.length"
        class="py-16 text-center text-sm text-gray-400"
      >
        暂无通知
      </p>

      <PublicPagination
        :page="page"
        :total-pages="totalPages"
        @change="goToPage"
      />
    </div>

    <div
      v-else
      class="rounded-xl border bg-white p-10 text-center"
    >
      <p class="text-gray-500">
        登录后查看通知
      </p>
      <NuxtLink
        to="/login?redirect=/notifications"
        class="mt-4 inline-block rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >登录</NuxtLink>
    </div>
  </div>
</template>
