<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Public user homepage: overview / topics / replies / favorites tabs,
 * all paginated through /api/v1/users/{id}/... */

interface Profile {
  id: number
  username: string
  avatar: string
  bio: string
  website: string
  profile_card?: string
  created_at: string
  counts: { topics: number, replies: number, favorites: number }
}

interface TopicRow {
  id: number
  title: string
  forum_category_id: number
  view_count: number
  reply_count: number
  created_at: string
}

interface ReplyRow {
  id: number
  topic_id: number
  topic_title: string
  content: string
  floor: number
  created_at: string
}

interface FavRow {
  id: number
  fav_type: 'topic' | 'article'
  fav_id: number
  title: string
  reply_count: number
  comment_count?: number
  article_slug?: string
  favorited_at: string
}

const route = useRoute()
const { user: me } = useAuthUser()

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

const profile = ref<Profile | null>(null)
const loading = ref(true)

type Tab = 'overview' | 'topics' | 'replies' | 'favorites'
const tab = ref<Tab>((route.query.tab as Tab) || 'overview')
const page = ref(1)
const totalPages = ref(1)

const topics = ref<TopicRow[]>([])
const replies = ref<ReplyRow[]>([])
const favorites = ref<FavRow[]>()

async function loadProfile() {
  try {
    profile.value = await rawFetch(`/api/v1/users/${route.params.id}`) as Profile
  } catch { /* silent */ }
  loading.value = false
}

async function loadTab() {
  const uid = route.params.id
  const p = `page=${page.value}&perPage=20`
  try {
    if (tab.value === 'topics') {
      const res = await rawFetch(`/api/v1/users/${uid}/topics?${p}`) as { items: TopicRow[], totalPages: number }
      topics.value = res.items ?? []
      totalPages.value = res.totalPages ?? 1
    } else if (tab.value === 'replies') {
      const res = await rawFetch(`/api/v1/users/${uid}/replies?${p}`) as { items: ReplyRow[], totalPages: number }
      replies.value = res.items ?? []
      totalPages.value = res.totalPages ?? 1
    } else if (tab.value === 'favorites') {
      const res = await rawFetch(`/api/v1/users/${uid}/favorites?${p}`) as { items: FavRow[], totalPages: number }
      favorites.value = res.items ?? []
      totalPages.value = res.totalPages ?? 1
    }  } catch { /* silent */ }
}

function switchTab(t: Tab) {
  if (tab.value === t) return
  tab.value = t
  page.value = 1
  router.replace({ query: t !== 'overview' ? { tab: t } : {} })
  void loadTab()
}

function goToPage(p: number) {
  page.value = p
  void loadTab()
}

const router = useRouter()
onMounted(async () => {
  await loadProfile()
  // overview previews + the active tab's list
  const uid = route.params.id
  try {
    const [t, r] = await Promise.all([
      rawFetch(`/api/v1/users/${uid}/topics?page=1&perPage=5`) as Promise<{ items: TopicRow[] }>,
      rawFetch(`/api/v1/users/${uid}/replies?page=1&perPage=5`) as Promise<{ items: ReplyRow[] }>
    ])
    topics.value = t.items ?? []
    replies.value = r.items ?? []
  } catch { /* silent */ }
  if (tab.value !== 'overview') await loadTab()
  void loadBlockState()
})

function fmtDate(s?: string) {
  return s?.slice(0, 10) ?? ''
}

/* ---------- block / DM ---------- */

const blocked = ref(false)
const blockBusy = ref(false)

async function loadBlockState() {
  if (!me.value || isOwner.value) return
  try {
    const res = await rawFetch('/api/v1/me/blocked') as { items: Array<{ blocked_id: number }> }
    blocked.value = (res.items ?? []).some(b => String(b.blocked_id) === String(route.params.id))
  } catch { /* silent */ }
}

async function toggleBlock() {
  if (blockBusy.value) return
  blockBusy.value = true
  try {
    const res = await rawFetch('/api/v1/me/blocked/toggle', {
      method: 'POST',
      body: { user_id: Number(route.params.id) }
    }) as { blocked: boolean, message?: string }
    blocked.value = res.blocked
    notify(res.message ?? (res.blocked ? '已拉黑' : '已取消拉黑'))
  } catch { /* silent */ } finally {
    blockBusy.value = false
  }
}

/* ---------- self-description markdown card ---------- */

const isOwner = computed(() => me.value?.id === profile.value?.id)
const cardEditing = ref(false)
const cardContent = ref('')
const cardSaving = ref(false)
const cardError = ref('')

const cardEditorRef = ref<HTMLTextAreaElement>()

const mdButtons = [
  { label: 'B', title: '粗体', before: '**', after: '**', cls: 'font-bold' },
  { label: 'I', title: '斜体', before: '*', after: '*', cls: 'italic' },
  { label: 'H2', title: '标题', before: '## ', after: '', cls: 'font-semibold' },
  { label: '🔗', title: '链接', before: '[', after: '](https://)', cls: '' },
  { label: '</>', title: '行内代码', before: '`', after: '`', cls: 'font-mono text-xs' },
  { label: '📋', title: '代码块', before: '```\n', after: '\n```', cls: 'font-mono text-xs' },
  { label: '📝', title: '引用', before: '> ', after: '', cls: '' },
  { label: '•', title: '列表', before: '- ', after: '', cls: '' }
]

function wrapCardSel(before: string, after = before) {
  const el = cardEditorRef.value
  if (!el) return
  const s = el.selectionStart
  const e = el.selectionEnd
  const sel = el.value.slice(s, e)
  el.value = el.value.slice(0, s) + before + sel + after + el.value.slice(e)
  el.selectionStart = s + before.length
  el.selectionEnd = s + before.length + sel.length
  el.focus()
  cardContent.value = el.value
}

function startCardEdit() {
  cardContent.value = profile.value?.profile_card ?? ''
  cardError.value = ''
  cardEditing.value = true
}

async function saveCard() {
  cardSaving.value = true
  cardError.value = ''
  try {
    await rawFetch('/api/v1/me/profile-card', {
      method: 'PUT',
      body: { content: cardContent.value }
    })
    if (profile.value) profile.value.profile_card = cardContent.value.trim()
    cardEditing.value = false
  } catch (e: unknown) {
    const err = e as { data?: { message?: string } }
    cardError.value = err.data?.message ?? '保存失败'
  } finally {
    cardSaving.value = false
  }
}

function snippet(s: string, n = 90): string {
  const line = s.replace(/[#*`>[\]()~-]/g, '').replace(/\s+/g, ' ').trim()
  return line.length > n ? line.slice(0, n) + '…' : line
}

useHead(() => ({ title: profile.value ? `${profile.value.username} 的主页` : '用户主页' }))
</script>

<template>
  <div class="mx-auto max-w-4xl">
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <template v-else-if="profile">
      <!-- profile header -->
      <div class="rounded-xl border bg-white p-6">
        <div class="flex items-start gap-4">
          <span class="flex h-16 w-16 shrink-0 items-center justify-center rounded-full bg-blue-500 text-xl font-bold text-white">
            {{ profile.username.slice(0, 1).toUpperCase() }}
          </span>
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-3">
              <h1 class="text-xl font-bold text-gray-900">
                {{ profile.username }}
              </h1>
              <NuxtLink
                v-if="isOwner"
                to="/settings"
                class="rounded-md border border-gray-300 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50"
              >编辑资料</NuxtLink>
              <template v-else-if="me">
                <NuxtLink
                  :to="`/messages?to=${profile.username}`"
                  class="rounded-md bg-blue-600 px-3 py-1 text-xs font-medium text-white hover:bg-blue-700"
                >✉️ 发私信</NuxtLink>
                <button
                  type="button"
                  :disabled="blockBusy"
                  class="rounded-md border px-3 py-1 text-xs hover:bg-gray-50"
                  :class="blocked
                    ? 'border-gray-300 text-gray-600'
                    : 'border-red-200 text-red-600 hover:bg-red-50'"
                  @click="toggleBlock"
                >
                  {{ blocked ? '已拉黑 · 取消' : '⛔ 拉黑' }}
                </button>
              </template>
            </div>
            <p
              v-if="profile.bio"
              class="mt-1 text-sm text-gray-600"
            >
              {{ profile.bio }}
            </p>
            <p class="mt-2 flex flex-wrap items-center gap-4 text-xs text-gray-400">
              <span>{{ fmtDate(profile.created_at) }} 加入</span>
              <a
                v-if="profile.website"
                :href="profile.website"
                target="_blank"
                rel="noopener noreferrer"
                class="text-blue-600 hover:underline"
              >{{ profile.website }}</a>
            </p>
          </div>
        </div>

        <!-- stats -->
        <div class="mt-5 grid grid-cols-3 gap-3 border-t pt-4 text-center">
          <div>
            <div class="text-lg font-bold text-gray-900">{{ profile.counts.topics }}</div>
            <div class="text-xs text-gray-400">主题帖</div>
          </div>
          <div>
            <div class="text-lg font-bold text-gray-900">{{ profile.counts.replies }}</div>
            <div class="text-xs text-gray-400">回复</div>
          </div>
          <div>
            <div class="text-lg font-bold text-gray-900">{{ profile.counts.favorites }}</div>
            <div class="text-xs text-gray-400">收藏</div>
          </div>
        </div>
      </div>

      <!-- self-description markdown card -->
      <div
        v-if="isOwner || profile.profile_card"
        class="mt-6 rounded-xl border bg-white p-5"
      >
        <div class="flex items-center justify-between">
          <h2 class="text-base font-semibold text-gray-900">
            个人卡片
          </h2>
          <button
            v-if="isOwner && !cardEditing"
            type="button"
            class="rounded-md border border-gray-300 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50"
            @click="startCardEdit"
          >
            {{ profile.profile_card ? '编辑卡片' : '创建卡片' }}
          </button>
        </div>

        <!-- owner editing: markdown toolbar + textarea -->
        <template v-if="cardEditing">
          <p class="mt-1 text-xs text-gray-400">
            支持 Markdown，保存后所有访客可见
          </p>
          <div class="mt-3 mb-2 flex flex-wrap gap-1">
            <button
              v-for="btn in mdButtons"
              :key="btn.title"
              type="button"
              :title="btn.title"
              class="rounded border border-gray-200 px-2 py-0.5 text-xs text-gray-600 hover:bg-gray-50"
              :class="btn.cls"
              @click.prevent="wrapCardSel(btn.before, btn.after)"
            >
              {{ btn.label }}
            </button>
          </div>
          <textarea
            ref="cardEditorRef"
            v-model="cardContent"
            :rows="8"
            placeholder="用 Markdown 自由编写你的个人卡片：自我介绍、技能、作品链接…"
            class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm leading-relaxed focus:border-blue-400 focus:outline-none focus:ring-1 focus:ring-blue-400"
          />
          <p
            v-if="cardError"
            class="mt-2 text-sm text-red-600"
          >
            {{ cardError }}
          </p>
          <div class="mt-3 flex items-center gap-3">
            <span class="text-xs text-gray-400">{{ cardContent.length }}/10000</span>
            <div class="ml-auto flex gap-2">
              <button
                type="button"
                class="rounded-lg border border-gray-300 px-4 py-1.5 text-sm text-gray-600 hover:bg-gray-50"
                @click="cardEditing = false"
              >
                取消
              </button>
              <button
                type="button"
                :disabled="cardSaving"
                class="rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
                @click="saveCard"
              >
                {{ cardSaving ? '保存中…' : '保存卡片' }}
              </button>
            </div>
          </div>
        </template>

        <!-- visitors (and the owner previewing): rendered markdown -->
        <div
          v-else-if="profile.profile_card"
          class="mt-3 text-sm leading-relaxed text-gray-700"
          v-html="renderMarkdown(profile.profile_card)"
        />
      </div>

      <!-- tabs -->
      <div class="mt-6 flex rounded-lg border border-gray-200 bg-white p-1">
        <button
          v-for="t in [
            { key: 'overview', label: '概览' },
            { key: 'topics', label: '主题帖' },
            { key: 'replies', label: '回复' },
            { key: 'favorites', label: '收藏' }
          ]"
          :key="t.key"
          type="button"
          class="flex-1 rounded-md px-3 py-1.5 text-sm font-medium"
          :class="tab === t.key ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-50'"
          @click="switchTab(t.key as Tab)"
        >
          {{ t.label }}
        </button>
      </div>

      <!-- overview -->
      <div
        v-if="tab === 'overview'"
        class="mt-4 space-y-4"
      >
        <div class="rounded-xl border bg-white p-5">
          <h2 class="mb-3 text-sm font-semibold text-gray-900">
            最近主题帖
          </h2>
          <div
            v-if="topics.length"
            class="space-y-2"
          >
            <NuxtLink
              v-for="t in topics"
              :key="t.id"
              :to="`/topics/${t.id}`"
              class="block rounded-md px-2 py-1.5 text-sm hover:bg-gray-50"
            >
              <span class="font-medium text-gray-800">{{ t.title }}</span>
              <span class="ml-2 text-xs text-gray-400">💬 {{ t.reply_count }} · 👁 {{ t.view_count }}</span>
            </NuxtLink>
          </div>
          <p
            v-else
            class="py-4 text-center text-sm text-gray-400"
          >
            暂无主题帖
          </p>
        </div>
        <div class="rounded-xl border bg-white p-5">
          <h2 class="mb-3 text-sm font-semibold text-gray-900">
            最近回复
          </h2>
          <div
            v-if="replies.length"
            class="space-y-3"
          >
            <NuxtLink
              v-for="r in replies"
              :key="r.id"
              :to="`/topics/${r.topic_id}`"
              class="block rounded-md px-2 py-1.5 hover:bg-gray-50"
            >
              <div class="truncate text-sm font-medium text-gray-800">
                {{ r.topic_title }}
              </div>
              <div class="truncate text-xs text-gray-500">
                {{ snippet(r.content, 60) }}
              </div>
            </NuxtLink>
          </div>
          <p
            v-else
            class="py-4 text-center text-sm text-gray-400"
          >
            暂无回复
          </p>
        </div>
      </div>

      <!-- topics -->
      <div
        v-else-if="tab === 'topics'"
        class="mt-4 divide-y divide-gray-100 rounded-xl border bg-white"
      >
        <NuxtLink
          v-for="t in topics"
          :key="t.id"
          :to="`/topics/${t.id}`"
          class="flex items-center justify-between px-5 py-3.5 hover:bg-gray-50"
        >
          <span class="min-w-0 truncate font-medium text-gray-900">{{ t.title }}</span>
          <span class="ml-3 shrink-0 text-xs text-gray-400">💬 {{ t.reply_count }} · 👁 {{ t.view_count }} · {{ fmtDate(t.created_at) }}</span>
        </NuxtLink>
        <p
          v-if="!topics.length"
          class="py-12 text-center text-sm text-gray-400"
        >
          暂无主题帖
        </p>
        <PublicPagination
          :page="page"
          :total-pages="totalPages"
          @change="goToPage"
        />
      </div>

      <!-- replies -->
      <div
        v-else-if="tab === 'replies'"
        class="mt-4 divide-y divide-gray-100 rounded-xl border bg-white"
      >
        <NuxtLink
          v-for="r in replies"
          :key="r.id"
          :to="`/topics/${r.topic_id}`"
          class="block px-5 py-3.5 hover:bg-gray-50"
        >
          <div class="flex items-center justify-between gap-3">
            <span class="truncate text-sm font-medium text-gray-900">{{ r.topic_title }}</span>
            <span class="shrink-0 text-xs text-gray-400">#{{ r.floor }} · {{ fmtDate(r.created_at) }}</span>
          </div>
          <p class="mt-1 truncate text-sm text-gray-500">
            {{ snippet(r.content) }}
          </p>
        </NuxtLink>
        <p
          v-if="!replies.length"
          class="py-12 text-center text-sm text-gray-400"
        >
          暂无回复
        </p>
        <PublicPagination
          :page="page"
          :total-pages="totalPages"
          @change="goToPage"
        />
      </div>

      <!-- favorites -->
      <div
        v-else
        class="mt-4 divide-y divide-gray-100 rounded-xl border bg-white"
      >
        <NuxtLink
          v-for="f in favorites ?? []"
          :key="f.id"
          :to="f.fav_type === 'article' ? `/articles/${f.article_slug}` : `/topics/${f.fav_id}`"
          class="flex items-center justify-between px-5 py-3.5 hover:bg-gray-50"
        >
          <span class="flex min-w-0 items-center gap-2">
            <span
              class="shrink-0 rounded px-1.5 py-0.5 text-xs"
              :class="f.fav_type === 'article' ? 'bg-indigo-100 text-indigo-600' : 'bg-blue-100 text-blue-600'"
            >{{ f.fav_type === 'article' ? '文章' : '帖子' }}</span>
            <span class="min-w-0 truncate font-medium text-gray-900">{{ f.title }}</span>
          </span>
          <span class="ml-3 shrink-0 text-xs text-gray-400">收藏于 {{ fmtDate(f.favorited_at) }}</span>
        </NuxtLink>
        <p
          v-if="!(favorites ?? []).length"
          class="py-12 text-center text-sm text-gray-400"
        >
          暂无收藏
        </p>
        <PublicPagination
          :page="page"
          :total-pages="totalPages"
          @change="goToPage"
        />
      </div>
    </template>

    <p
      v-else
      class="py-20 text-center text-gray-400"
    >
      用户不存在
    </p>
  </div>
</template>
