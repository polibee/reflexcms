<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Direct messages: inbox / sent tabs plus a compose form (send by
 * username). Reading the inbox marks rows read via the read endpoint. */

interface MessageRow {
  id: number
  sender_id: number
  recipient_id: number
  body: string
  read_at: string | null
  created_at: string
  from_name?: string
  to_name?: string
}

const { user: me } = useAuthUser()
const route = useRoute()
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

type Box = 'inbox' | 'sent'
const box = ref<Box>(route.query.box === 'sent' ? 'sent' : 'inbox')
const items = ref<MessageRow[]>([])
const page = ref(1)
const totalPages = ref(1)
const loading = ref(true)

/* compose */
const toUsername = ref('')
const body = ref('')
const sendMsg = ref('')
const sendErr = ref('')
const sending = ref(false)
const editorRef = ref<HTMLTextAreaElement>()

const mdButtons = [
  { label: 'B', title: '粗体', before: '**', after: '**', cls: 'font-bold' },
  { label: 'I', title: '斜体', before: '*', after: '*', cls: 'italic' },
  { label: 'H2', title: '标题', before: '## ', after: '', cls: 'font-semibold' },
  { label: '🔗', title: '链接', before: '[', after: '](https://)', cls: '' },
  { label: '</>', title: '行内代码', before: '`', after: '`', cls: 'font-mono text-xs' },
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
  body.value = el.value
}

async function fetchList() {
  loading.value = true
  try {
    const res = await rawFetch(`/api/v1/me/messages?box=${box.value}&page=${page.value}&perPage=15`) as { items: MessageRow[], totalPages: number }
    items.value = res.items ?? []
    totalPages.value = res.totalPages ?? 1
  } catch { /* silent */ }
  loading.value = false
}

function switchBox(b: Box) {
  if (box.value === b) return
  box.value = b
  page.value = 1
  router.replace({ query: b === 'sent' ? { box: 'sent' } : {} })
  void fetchList()
}

const router = useRouter()

function goToPage(p: number) {
  page.value = p
  void fetchList()
}

async function markRead(m: MessageRow) {
  if (m.read_at) return
  try {
    await rawFetch(`/api/v1/me/messages/${m.id}/read`, { method: 'POST' })
    m.read_at = new Date().toISOString()
  } catch { /* silent */ }
}

async function send() {
  sending.value = true
  sendMsg.value = ''
  sendErr.value = ''
  try {
    const res = await rawFetch('/api/v1/me/messages', {
      method: 'POST',
      body: { to_username: toUsername.value.trim(), body: body.value }
    }) as { message?: string }
    sendMsg.value = res.message ?? '已发送'
    toUsername.value = ''
    body.value = ''
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string, message?: string } }
    sendErr.value = err.data?.statusMessage ?? err.data?.message ?? '发送失败'
  } finally {
    sending.value = false
  }
}

onMounted(async () => {
  // ?to=username prefill (from profile page 发私信 button)
  const to = route.query.to
  if (typeof to === 'string' && to) toUsername.value = to
  // unconditional: anonymous callers just get a 401 that we swallow
  await fetchList()
})

useHead({ title: '私信' })
</script>

<template>
  <div class="mx-auto max-w-3xl">
    <div class="mb-5 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">
        私信
      </h1>
      <div class="flex rounded-lg border border-gray-200 bg-white p-1">
        <button
          type="button"
          class="rounded-md px-4 py-1.5 text-sm font-medium"
          :class="box === 'inbox' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-50'"
          @click="switchBox('inbox')"
        >
          收件箱
        </button>
        <button
          type="button"
          class="rounded-md px-4 py-1.5 text-sm font-medium"
          :class="box === 'sent' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-50'"
          @click="switchBox('sent')"
        >
          已发送
        </button>
      </div>
    </div>

    <template v-if="me">
      <!-- compose -->
      <div class="mb-5 rounded-xl border bg-white p-5">
        <h2 class="text-sm font-semibold text-gray-900">
          写私信
        </h2>
        <input
          v-model="toUsername"
          type="text"
          placeholder="对方用户名"
          class="mt-3 h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
        >
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
          v-model="body"
          :rows="5"
          maxlength="2000"
          placeholder="私信内容…支持 Markdown 语法（最多 2000 字）"
          class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm leading-relaxed focus:border-blue-400 focus:outline-none focus:ring-1 focus:ring-blue-400"
        />
        <div class="mt-2 flex items-center justify-end gap-3">
          <span
            v-if="sendMsg"
            class="text-xs text-green-600"
          >{{ sendMsg }}</span>
          <span
            v-else-if="sendErr"
            class="text-xs text-red-600"
          >{{ sendErr }}</span>
          <button
            type="button"
            :disabled="sending || !toUsername.trim() || !body.trim()"
            class="rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            @click="send"
          >
            {{ sending ? '发送中…' : '发送' }}
          </button>
        </div>
      </div>

      <!-- list -->
      <div class="divide-y divide-gray-100 rounded-xl border bg-white">
        <div
          v-for="m in items"
          :key="m.id"
          class="block px-5 py-3.5 hover:bg-gray-50"
          @click="box === 'inbox' && markRead(m)"
        >
          <div class="flex items-center justify-between gap-3">
            <span
              class="truncate text-sm font-medium"
              :class="!m.read_at && box === 'inbox' ? 'text-blue-600' : 'text-gray-900'"
            >
              {{ box === 'inbox' ? `来自 ${m.from_name ?? '匿名'}` : `发给 ${m.to_name ?? '?'}` }}
              <span
                v-if="!m.read_at && box === 'inbox'"
                class="ml-1 rounded bg-blue-100 px-1.5 py-0.5 text-xs"
              >未读</span>
            </span>
            <span class="shrink-0 text-xs text-gray-400">{{ m.created_at?.slice(0, 16).replace('T', ' ') }}</span>
          </div>
          <p
            class="mt-1 text-sm leading-relaxed text-gray-600"
            v-html="renderMarkdown(m.body)"
          />
        </div>

        <p
          v-if="!items.length"
          class="py-16 text-center text-sm text-gray-400"
        >
          {{ box === 'inbox' ? '收件箱为空' : '还没有发送过私信' }}
        </p>

        <PublicPagination
          :page="page"
          :total-pages="totalPages"
          @change="goToPage"
        />
      </div>
    </template>

    <div
      v-else
      class="rounded-xl border bg-white p-10 text-center"
    >
      <p class="text-gray-500">
        登录后使用私信
      </p>
      <NuxtLink
        to="/login?redirect=/messages"
        class="mt-4 inline-block rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >登录</NuxtLink>
    </div>
  </div>
</template>
