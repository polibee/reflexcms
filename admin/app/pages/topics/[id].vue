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
  author_name?: string
  author_signature?: string
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
  can_moderate?: boolean
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
const { toast, confirmDialog } = useFrontUi()

async function blockUser(userId: number) {
  if (!me.value || blockBusyId.value) return
  const ok = await confirmDialog({
    title: '拉黑该用户？',
    message: '拉黑后将不再看到对方发表的回复，可随时在个人设置中解除。',
    confirmLabel: '拉黑',
    danger: true
  })
  if (!ok) return
  blockBusyId.value = userId
  try {
    await rawFetch('/api/v1/me/blocked/toggle', {
      method: 'POST',
      body: { user_id: userId }
    })
    // refetch current reply page — blocked author's replies are filtered server-side
    await goToReplyPage(replyPage.value)
    toast('已更新屏蔽状态')
  } catch {
    toast('操作失败，请稍后再试', 'error')
  } finally {
    blockBusyId.value = 0
  }
}

/* ---------- moderator actions (board moderators / super-admin) ---------- */

const modBusy = ref(false)

async function moderate(action: string, replyId?: number) {
  if (!topic.value || modBusy.value) return
  modBusy.value = true
  try {
    await rawFetch(`/api/v1/topics/${route.params.id}/moderate`, {
      method: 'POST',
      body: { action, reply_id: replyId ?? 0 }
    })
    const fresh = await rawFetch(`/api/v1/topics/${route.params.id}`) as TopicDetail
    fresh.replies = (fresh.replies ?? []).slice(0, repliesPerPage)
    topic.value = fresh
    replyTotal.value = fresh.reply_count
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string } }
    toast(err.data?.statusMessage ?? '操作失败', 'error')
  } finally {
    modBusy.value = false
  }
}

/* mute / ban from the reply card (moderators mute; admins can also ban) */
const muteBusyId = ref(0)
const muteMenuId = ref(0)

const MUTE_PRESETS = [
  { label: '禁言 1 小时', hours: 1 },
  { label: '禁言 6 小时', hours: 6 },
  { label: '禁言 12 小时', hours: 12 },
  { label: '禁言 1 天', hours: 24 },
  { label: '禁言 3 天', hours: 72 },
  { label: '禁言 7 天', hours: 168 },
  { label: '禁言 30 天', hours: 720 },
  { label: '永久禁言', hours: 875999 }
]

async function muteUser(userId: number, hours: number) {
  if (muteBusyId.value) return
  muteBusyId.value = userId
  muteMenuId.value = 0
  try {
    const res = await rawFetch(`/api/v1/users/${userId}/mute`, {
      method: 'POST',
      body: { duration_hours: hours }
    }) as { message?: string }
    toast(res.message ?? '已禁言')
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string } }
    notifyFront(err.data?.statusMessage ?? '操作失败')
  } finally {
    muteBusyId.value = 0
  }
}

async function banUser(userId: number, days: number) {
  if (muteBusyId.value) return
  const ok = await confirmDialog({
    title: '封禁该账号？',
    message: days >= 3650 ? '封禁后该用户将无法再登录本站。' : `该用户将被禁止登录 ${days} 天。`,
    confirmLabel: '封禁',
    danger: true
  })
  if (!ok) return
  muteBusyId.value = userId
  try {
    const res = await rawFetch(`/api/v1/users/${userId}/ban`, {
      method: 'POST',
      body: { days }
    }) as { message?: string }
    toast(res.message ?? '已封禁')
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string } }
    toast(err.data?.statusMessage ?? '操作失败', 'error')
  } finally {
    muteBusyId.value = 0
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

        <!-- moderator action bar -->
        <div
          v-if="topic.can_moderate"
          class="mt-3 flex flex-wrap items-center gap-2 rounded-lg bg-blue-50 px-3 py-2"
        >
          <span class="text-xs font-medium text-blue-700">版主操作：</span>
          <button
            v-for="act in [
              { key: topic.is_pinned ? 'unpin' : 'pin', label: topic.is_pinned ? '取消置顶' : '置顶' },
              { key: topic.is_featured ? 'unfeature' : 'feature', label: topic.is_featured ? '取消精华' : '精华' },
              { key: topic.status === 'open' ? 'close' : 'open', label: topic.status === 'open' ? '关闭' : '重新开放' }
            ]"
            :key="act.key"
            type="button"
            :disabled="modBusy"
            class="rounded-md border border-blue-200 bg-white px-2.5 py-1 text-xs text-blue-700 hover:bg-blue-50 disabled:opacity-50"
            @click="moderate(act.key)"
          >
            {{ act.label }}
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
                    <!-- moderator: set best reply / mute -->
                    <button
                      v-if="topic.can_moderate && reply.id !== topic.best_reply_id"
                      type="button"
                      title="设为最佳回复"
                      class="text-xs text-gray-300 hover:text-green-600"
                      :disabled="modBusy"
                      @click="moderate('best_reply', reply.id)"
                    >
                      ⭐
                    </button>
                    <button
                      v-if="topic.can_moderate && reply.user_id !== me?.id"
                      type="button"
                      title="禁言"
                      class="text-xs text-gray-300 hover:text-orange-500"
                      :disabled="muteBusyId === reply.user_id"
                      @click="muteMenuId = muteMenuId === reply.user_id ? 0 : reply.user_id"
                    >
                      🔇
                    </button>
                    <button
                      v-if="me?.role === 'super-admin' && reply.user_id !== me.id"
                      type="button"
                      title="封禁 7 天（仅超级管理员）"
                      class="text-xs text-gray-300 hover:text-red-600"
                      :disabled="muteBusyId === reply.user_id"
                      @click="banUser(reply.user_id, 7)"
                    >
                      ⛔️
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
                <p
                  v-if="reply.author_signature"
                  class="mt-2 border-t border-dashed border-gray-100 pt-2 text-xs text-gray-400"
                >
                  <span v-html="renderMarkdown(reply.author_signature)" />
                </p>
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

        <!-- mute duration menu -->
        <div
          v-if="muteMenuId"
          class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm"
          @click="muteMenuId = 0"
        >
          <div
            class="absolute left-1/2 top-1/2 w-60 -translate-x-1/2 -translate-y-1/2 rounded-xl border border-gray-100 bg-white p-3 shadow-2xl"
            @click.stop
          >
            <p class="mb-2 flex items-center gap-2 px-1 text-sm font-semibold text-gray-900">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                class="h-4 w-4 text-orange-500"
                aria-hidden="true"
              >
                <circle
                  cx="12"
                  cy="12"
                  r="9"
                /><path d="M12 7v5l3 2" />
              </svg>
              选择禁言时长
            </p>
            <button
              v-for="preset in MUTE_PRESETS"
              :key="preset.hours"
              type="button"
              class="block w-full rounded-lg px-3 py-1.5 text-left text-sm text-gray-700 transition-colors hover:bg-orange-50 hover:text-orange-700 disabled:opacity-50"
              :disabled="muteBusyId !== 0"
              @click="muteUser(muteMenuId, preset.hours)"
            >
              {{ preset.label }}
            </button>
          </div>
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
