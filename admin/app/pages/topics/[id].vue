<script setup lang="ts">
definePageMeta({ layout: 'public' })

interface Reply {
  id: number
  user_id: number
  parent_id?: number
  content: string
  floor: number
  like_count: number
  created_at: string
}

interface TopicDetail {
  id: number
  title: string
  content: string
  status: string
  is_pinned: boolean
  is_featured: boolean
  reply_count: number
  like_count: number
  view_count: number
  best_reply_id: number
  replies: Reply[]
}

const route = useRoute()
const topic = ref<TopicDetail | null>(null)
const loading = ref(true)

const { user: me } = useAuthUser()

const rawFetch = $fetch as unknown as (url: string, opts?: { method?: string, body?: unknown }) => Promise<unknown>

/* replies are paginated: page 1 arrives embedded with the topic, deeper
 * pages load from the dedicated endpoint */
const replyPage = ref(1)
const replyTotal = ref(0)
const replyTotalPages = computed(() => Math.max(1, Math.ceil(replyTotal.value / repliesPerPage)))
const repliesPerPage = 10

onMounted(async () => {
  try {
    const data = await rawFetch(`/api/v1/topics/${route.params.id}`) as TopicDetail
    data.replies = data.replies ?? []
    topic.value = data
    replyTotal.value = data.reply_count
    // the API embeds up to 50; trim the client page size for consistency
    if (data.replies.length > repliesPerPage) data.replies = data.replies.slice(0, repliesPerPage)
  } catch { /* silent */ }
  loading.value = false
  void loadFavoriteState()
})

async function goToReplyPage(p: number) {
  replyPage.value = p
  try {
    const res = await rawFetch(`/api/v1/topics/${route.params.id}/replies?page=${p}&perPage=${repliesPerPage}`) as { items: Reply[], total: number }
    if (topic.value) topic.value.replies = res.items ?? []
    replyTotal.value = res.total ?? replyTotal.value
  } catch { /* silent */ }
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

/* favorite toggle */
const favorited = ref(false)

async function loadFavoriteState() {
  if (!me.value) {
    favorited.value = false
    return
  }
  try {
    const res = await rawFetch(`/api/v1/favorites/state?type=topic&id=${route.params.id}`) as { favorited: boolean }
    favorited.value = res.favorited
  } catch { /* silent */ }
}

async function toggleFavorite() {
  if (!me.value) return
  try {
    const res = await rawFetch('/api/v1/favorites/toggle', {
      method: 'POST',
      body: { type: 'topic', id: Number(route.params.id) }
    }) as { favorited: boolean }
    favorited.value = res.favorited
  } catch { /* silent */ }
}

useHead(() => ({ title: topic.value?.title ?? '帖子' }))

/* ---------- reply composer ---------- */

const replyContent = ref('')
const submitting = ref(false)
const replyError = ref('')

const editorRef = ref<HTMLTextAreaElement>()

const mdButtons = [
  { label: 'B', title: '粗体', before: '**', after: '**', cls: 'font-bold' },
  { label: 'I', title: '斜体', before: '*', after: '*', cls: 'italic' },
  { label: 'H2', title: '标题', before: '## ', after: '', cls: 'font-semibold' },
  { label: '🔗', title: '链接', before: '[', after: '](url)', cls: '' },
  { label: '</>', title: '行内代码', before: '`', after: '`', cls: 'font-mono text-xs' },
  { label: '📋', title: '代码块', before: '```\n', after: '\n```', cls: 'font-mono text-xs' },
  { label: '📝', title: '引用', before: '> ', after: '', cls: '' },
  { label: '•', title: '列表', before: '- ', after: '', cls: '' }
]

function wrapSel(before: string, after = before) {
  const el = editorRef.value
  if (!el) return
  const s = el.selectionStart
  const e = el.selectionEnd
  const sel = el.value.slice(s, e)
  el.value = el.value.slice(0, s) + before + sel + after + el.value.slice(e)
  el.selectionStart = s + before.length
  el.selectionEnd = s + before.length + sel.length
  el.focus()
  replyContent.value = el.value
}

async function submitReply() {
  if (!me.value || submitting.value) return
  const content = replyContent.value.trim()
  if (!content) {
    replyError.value = '回复内容不能为空'
    return
  }
  submitting.value = true
  replyError.value = ''
  try {
    const reply = await rawPostReply(content)
    if (topic.value && replyPage.value === 1) {
      topic.value.replies.push(reply)
      // keep the client page window at repliesPerPage
      if (topic.value.replies.length > repliesPerPage) topic.value.replies.shift()
    }
    if (topic.value) topic.value.reply_count++
    replyTotal.value++
    replyContent.value = ''
  } catch (e: unknown) {
    const err = e as { statusCode?: number, data?: { statusMessage?: string, message?: string } }
    if (err.statusCode === 401) {
      replyError.value = '登录状态已过期，请重新登录'
    } else {
      replyError.value = err.data?.statusMessage ?? err.data?.message ?? '回复失败，请重试'
    }
  } finally {
    submitting.value = false
  }
}

function rawPostReply(content: string): Promise<Reply> {
  return rawFetch(`/api/v1/topics/${route.params.id}/replies`, {
    method: 'POST',
    body: { content }
  }) as Promise<Reply>
}

/* ---------- shared ---------- */

const avatarColors = [
  'bg-blue-500', 'bg-green-500', 'bg-purple-500', 'bg-orange-500',
  'bg-pink-500', 'bg-teal-500', 'bg-indigo-500'
]
function avatarColor(userId: number) {
  return avatarColors[userId % avatarColors.length]
}

/* block a user from the reply card; their replies vanish after reload */
const blockBusyId = ref(0)

async function blockUser(userId: number) {
  if (!me.value || blockBusyId.value) return
  blockBusyId.value = userId
  try {
    await rawFetch('/api/v1/me/blocked/toggle', {
      method: 'POST',
      body: { user_id: userId }
    })
    // refetch current reply page — blocked author's replies are filtered server-side
    await goToReplyPage(replyPage.value)
  } catch { /* silent */ } finally {
    blockBusyId.value = 0
  }
}
</script>

<template>
  <div class="mx-auto max-w-4xl">
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>
    <template v-else-if="topic">
      <!-- Topic header -->
      <div class="rounded-xl border bg-white p-6">
        <div class="flex items-center gap-2">
          <span
            v-if="topic.is_pinned"
            class="rounded bg-red-100 px-1.5 py-0.5 text-xs font-medium text-red-600"
          >置顶</span>
          <span
            v-if="topic.is_featured"
            class="rounded bg-blue-100 px-1.5 py-0.5 text-xs font-medium text-blue-600"
          >精华</span>
          <span
            v-if="topic.status === 'closed'"
            class="rounded bg-gray-100 px-1.5 py-0.5 text-xs font-medium text-gray-500"
          >已关闭</span>
        </div>
        <h1 class="mt-2 text-2xl font-bold text-gray-900">
          {{ topic.title }}
        </h1>
        <div class="mt-3 flex items-center gap-4 text-sm text-gray-400">
          <span>👁 {{ topic.view_count }}</span>
          <span>💬 {{ topic.reply_count }}</span>
          <span>👍 {{ topic.like_count }}</span>
          <button
            v-if="me"
            type="button"
            class="ml-auto rounded-md px-3 py-1 text-xs font-medium"
            :class="favorited
              ? 'bg-amber-100 text-amber-700 hover:bg-amber-200'
              : 'border border-gray-300 text-gray-500 hover:bg-gray-50'"
            @click="toggleFavorite"
          >
            {{ favorited ? '★ 已收藏' : '☆ 收藏' }}
          </button>
        </div>
        <div
          class="mt-4 border-t pt-4 text-base leading-relaxed text-gray-700"
          v-html="renderMarkdown(topic.content)"
        />
      </div>

      <!-- Replies -->
      <div class="mt-6">
        <h2 class="mb-4 text-lg font-semibold text-gray-900">
          {{ topic.reply_count }} 条回复
        </h2>

        <div class="space-y-4">
          <div
            v-for="reply in topic.replies"
            :key="reply.id"
            class="rounded-xl border bg-white p-5"
            :class="reply.id === topic.best_reply_id
              ? 'border-green-300 bg-green-50/50 ring-1 ring-green-200'
              : ''"
          >
            <div class="flex items-start gap-3">
              <!-- Avatar + name link to the user's public homepage -->
              <NuxtLink
                :to="`/users/${reply.user_id}`"
                class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-sm font-bold text-white"
                :class="avatarColor(reply.user_id)"
                :title="`用户 ${reply.user_id} 的主页`"
              >
                U{{ reply.user_id }}
              </NuxtLink>

              <!-- Content -->
              <div class="min-w-0 flex-1">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <NuxtLink
                      :to="`/users/${reply.user_id}`"
                      class="text-sm font-medium text-gray-900 hover:text-blue-600 hover:underline"
                    >用户 {{ reply.user_id }}</NuxtLink>
                    <span
                      v-if="reply.id === topic.best_reply_id"
                      class="rounded bg-green-100 px-1.5 py-0.5 text-xs font-medium text-green-700"
                    >✓ 最佳回复</span>
                    <button
                      v-if="me && reply.user_id !== me.id"
                      type="button"
                      title="拉黑该用户，不再显示其回复"
                      class="text-xs text-gray-300 hover:text-red-500"
                      :disabled="blockBusyId === reply.user_id"
                      @click="blockUser(reply.user_id)"
                    >
                      ⛔
                    </button>
                  </div>
                  <span class="text-xs text-gray-400">#{{ reply.floor }}</span>
                </div>
                <div
                  class="mt-2 text-sm leading-relaxed text-gray-700"
                  v-html="renderMarkdown(reply.content)"
                />
                <div class="mt-3 flex items-center gap-4 text-xs text-gray-400">
                  <time>{{ reply.created_at?.slice(0, 10) }}</time>
                  <span>👍 {{ reply.like_count }}</span>
                </div>
              </div>
            </div>
          </div>

          <p
            v-if="!topic.replies?.length"
            class="py-10 text-center text-gray-400"
          >
            暂无回复
          </p>
        </div>

        <PublicPagination
          :page="replyPage"
          :total-pages="replyTotalPages"
          @change="goToReplyPage"
        />
      </div>

      <!-- Reply composer (below the reply list) -->
      <div class="mt-6 rounded-xl border bg-white p-5">
        <h2 class="text-base font-semibold text-gray-900">
          发表回复
        </h2>

        <template v-if="topic.status === 'open'">
          <!-- logged in: markdown editor -->
          <template v-if="me">
            <p class="mt-1 text-xs text-gray-400">
              以 <span class="font-medium text-gray-600">{{ me.name }}</span> 身份回复 · 支持 Markdown
            </p>
            <div class="mt-3 mb-2 flex flex-wrap gap-1">
              <button
                v-for="btn in mdButtons"
                :key="btn.title"
                type="button"
                :title="btn.title"
                class="rounded border border-gray-200 px-2 py-0.5 text-xs text-gray-600 hover:bg-gray-50"
                :class="btn.cls"
                @click.prevent="wrapSel(btn.before, btn.after)"
              >
                {{ btn.label }}
              </button>
            </div>
            <textarea
              ref="editorRef"
              v-model="replyContent"
              :rows="5"
              placeholder="写下你的回复…（支持 **粗体**、`代码`、[链接](url) 等 Markdown 语法）"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm leading-relaxed focus:border-blue-400 focus:outline-none focus:ring-1 focus:ring-blue-400"
            />
            <div class="mt-2 flex items-center justify-between">
              <span class="text-xs text-gray-400">{{ replyContent.length }}/5000</span>
              <button
                type="button"
                :disabled="submitting || !replyContent.trim()"
                class="rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
                @click="submitReply"
              >
                {{ submitting ? '发布中…' : '发布回复' }}
              </button>
            </div>
            <p
              v-if="replyError"
              class="mt-2 text-sm text-red-600"
            >
              {{ replyError }}
            </p>
          </template>

          <!-- anonymous: login prompt -->
          <div
            v-else
            class="mt-3 flex items-center justify-between rounded-lg bg-gray-50 px-4 py-3"
          >
            <p class="text-sm text-gray-500">
              登录后即可参与回复
            </p>
            <NuxtLink
              :to="`/login?redirect=/topics/${topic.id}`"
              class="rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700"
            >登录</NuxtLink>
          </div>
        </template>

        <p
          v-else
          class="mt-3 rounded-lg bg-gray-50 px-4 py-3 text-sm text-gray-500"
        >
          帖子已关闭，无法回复
        </p>
      </div>
    </template>
    <p
      v-else
      class="py-20 text-center text-gray-400"
    >
      帖子不存在
    </p>
  </div>
</template>
