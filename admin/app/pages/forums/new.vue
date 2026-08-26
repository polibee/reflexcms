<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Topic composer: board picker + markdown editor. Requires a session;
 * anonymous visitors get a login prompt with redirect back. */

interface Board {
  id: number
  name: string
}

const router = useRouter()
const { user: me } = useAuthUser()

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

const boards = ref<Board[]>([])
const boardId = ref<number | ''>('')
const title = ref('')
const content = ref('')
const submitting = ref(false)
const errorMsg = ref('')

const editorRef = ref<HTMLTextAreaElement>()

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
  content.value = el.value
}

async function submitTopic() {
  if (submitting.value) return
  errorMsg.value = ''
  if (!boardId.value) {
    errorMsg.value = '请选择板块'
    return
  }
  if (!title.value.trim()) {
    errorMsg.value = '标题不能为空'
    return
  }
  if (!content.value.trim()) {
    errorMsg.value = '正文不能为空'
    return
  }
  submitting.value = true
  try {
    const res = await rawFetch('/api/v1/topics', {
      method: 'POST',
      body: {
        title: title.value.trim(),
        content: content.value,
        forum_category_id: boardId.value
      }
    }) as { id: number }
    await navigateTo(`/topics/${res.id}`)
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string, message?: string } }
    errorMsg.value = err.data?.statusMessage ?? err.data?.message ?? '发布失败，请重试'
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const res = await rawFetch('/api/v1/forums') as { items: Board[] }
    boards.value = res.items ?? []
  } catch { /* silent */ }
  // preselect board when arriving from a board page (?board=ID)
  const q = Number(router.currentRoute.value.query.board)
  if (q > 0) boardId.value = q
})

useHead({ title: '发布帖子' })
</script>

<template>
  <div class="mx-auto max-w-3xl">
    <h1 class="mb-6 text-2xl font-bold text-gray-900">
      发布帖子
    </h1>

    <template v-if="me">
      <div class="space-y-4 rounded-xl border bg-white p-6">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm text-gray-600">选择板块 <span class="text-red-500">*</span></label>
            <select
              v-model="boardId"
              class="h-9 w-full rounded-md border border-gray-200 bg-white px-3 text-sm focus:border-blue-400 focus:outline-none"
            >
              <option
                value=""
                disabled
              >
                请选择板块
              </option>
              <option
                v-for="b in boards"
                :key="b.id"
                :value="b.id"
              >
                {{ b.name }}
              </option>
            </select>
          </div>
          <div>
            <label class="mb-1 block text-sm text-gray-600">标题 <span class="text-red-500">*</span></label>
            <input
              v-model="title"
              type="text"
              maxlength="150"
              placeholder="一句话说清你的问题或主题"
              class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
            >
            <p class="mt-0.5 text-xs text-gray-400">
              {{ title.length }}/150
            </p>
          </div>
        </div>

        <div>
          <label class="mb-2 block text-sm text-gray-600">正文（支持 Markdown） <span class="text-red-500">*</span></label>
          <div class="mb-2 flex flex-wrap gap-1">
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
            v-model="content"
            :rows="12"
            placeholder="详细描述你的内容…支持 **粗体**、`代码`、代码块、列表等 Markdown 语法"
            class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm leading-relaxed focus:border-blue-400 focus:outline-none focus:ring-1 focus:ring-blue-400"
          />
        </div>

        <p
          v-if="errorMsg"
          class="text-sm text-red-600"
        >
          {{ errorMsg }}
        </p>

        <div class="flex items-center justify-end gap-3 border-t pt-4">
          <button
            type="button"
            class="rounded-lg border border-gray-300 px-4 py-1.5 text-sm text-gray-600 hover:bg-gray-50"
            @click="router.back()"
          >
            取消
          </button>
          <button
            type="button"
            :disabled="submitting"
            class="rounded-lg bg-blue-600 px-5 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            @click="submitTopic"
          >
            {{ submitting ? '发布中…' : '发布帖子' }}
          </button>
        </div>
      </div>
    </template>

    <div
      v-else
      class="rounded-xl border bg-white p-10 text-center"
    >
      <p class="text-gray-500">
        登录后即可发布帖子
      </p>
      <NuxtLink
        to="/login?redirect=/forums/new"
        class="mt-4 inline-block rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >登录</NuxtLink>
    </div>
  </div>
</template>
