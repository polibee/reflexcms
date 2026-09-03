<script setup lang="ts">
definePageMeta({ layout: 'default' })

/* Admin mailer: send one email or bulk-paste/upload a recipient list.
 * Delivery goes through the backend, which uses the email module's
 * configured driver (log / resend / smtp). */

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

const to = ref('')
const subject = ref('')
const body = ref('')
const sending = ref(false)
const result = ref('')

const fileInput = ref<HTMLInputElement>()

function addRecipients(text: string) {
  to.value = to.value
    ? to.value + '\n' + text
    : text
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    addRecipients(String(reader.result ?? ''))
  }
  reader.readAsText(file)
  input.value = ''
}

async function send() {
  sending.value = true
  result.value = ''
  try {
    const res = await rawFetch('/api/admin/mail/send', {
      method: 'POST',
      body: { to: to.value.split('\n'), csv_text: '', subject: subject.value, body: body.value }
    }) as { sent: string[], failed: string[], total: number }
    result.value = `发送完成：成功 ${res.sent.length} / 共 ${res.total}${res.failed.length ? '，失败：' + res.failed.join(', ') : ''}`
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string, message?: string } }
    result.value = err.data?.statusMessage ?? err.data?.message ?? '发送失败'
  } finally {
    sending.value = false
  }
}

useHead({ title: '邮件群发' })
</script>

<template>
  <div class="mx-auto max-w-3xl space-y-6 p-6">
    <h1 class="text-2xl font-bold text-gray-900">
      邮件发送
    </h1>

    <div class="space-y-4 rounded-xl border bg-white p-6">
      <div>
        <label class="mb-1 block text-sm text-gray-600">收件人（每行一个邮箱）</label>
        <textarea
          v-model="to"
          :rows="5"
          placeholder="user1@example.com&#10;user2@example.com"
          class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
        />
        <div class="mt-1 flex items-center gap-3 text-xs text-gray-400">
          <span>支持手动输入、逗号分隔或从 txt/csv 文件导入（单次最多 500 个）</span>
          <button
            type="button"
            class="rounded border border-input px-2 py-0.5 hover:bg-accent"
            @click="fileInput?.click()"
          >
            从文件导入
          </button>
          <input
            ref="fileInput"
            type="file"
            accept=".txt,.csv"
            class="hidden"
            @change="onFileChange"
          >
        </div>
      </div>

      <div>
        <label class="mb-1 block text-sm text-gray-600">邮件主题</label>
        <input
          v-model="subject"
          type="text"
          class="h-9 w-full rounded-md border border-input bg-background px-3 text-sm"
        >
      </div>

      <div>
        <label class="mb-1 block text-sm text-gray-600">邮件正文</label>
        <textarea
          v-model="body"
          :rows="8"
          placeholder="支持 HTML 内容"
          class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
        />
      </div>

      <div class="flex items-center justify-between">
        <span
          v-if="result"
          class="text-xs"
          :class="result.startsWith('发送完成') ? 'text-green-600' : 'text-red-600'"
        >{{ result }}</span>
        <button
          type="button"
          :disabled="sending || !subject.trim() || !body.trim() || !to.trim()"
          class="ml-auto rounded-lg bg-blue-600 px-5 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
          @click="send"
        >
          {{ sending ? '发送中…' : '发送邮件' }}
        </button>
      </div>
    </div>

    <p class="text-xs text-gray-400">
      发信驱动与发件人信息在「站点设置 → 邮件配置」中维护；开发环境（driver=log）只写日志不真实发信。
    </p>
  </div>
</template>
