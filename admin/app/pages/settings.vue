<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Self-service account settings: profile (username/bio/website) and
 * security (password change). Both write through the session token. */

const { user: me, refresh: refreshAuth } = useAuthUser()

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

/* profile form */
const username = ref('')
const bio = ref('')
const website = ref('')
const signature = ref('')
const email = ref('')
const profileMsg = ref('')
const profileErr = ref('')
const savingProfile = ref(false)

async function loadProfile() {
  try {
    const p = await rawFetch('/api/v1/me/profile') as { username: string, email: string, bio: string, website: string, signature?: string }
    username.value = p.username ?? ''
    email.value = p.email ?? ''
    bio.value = p.bio ?? ''
    website.value = p.website ?? ''
    signature.value = p.signature ?? ''
  } catch { /* silent */ }
}

async function saveProfile() {
  savingProfile.value = true
  profileMsg.value = ''
  profileErr.value = ''
  try {
    const res = await rawFetch('/api/v1/me/profile', {
      method: 'PUT',
      body: { username: username.value, bio: bio.value, website: website.value, signature: signature.value }
    }) as { message?: string }
    profileMsg.value = res.message ?? '资料已更新'
    void refreshAuth()
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string, message?: string } }
    profileErr.value = err.data?.statusMessage ?? err.data?.message ?? '保存失败'
  } finally {
    savingProfile.value = false
  }
}

/* password form */
const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const pwMsg = ref('')
const pwErr = ref('')
const savingPw = ref(false)

/* profile card (markdown) form */
const cardContent = ref('')
const cardMsg = ref('')
const cardErr = ref('')
const savingCard = ref(false)
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

async function loadCard() {
  try {
    const p = await rawFetch('/api/v1/me/profile') as { profile_card?: string }
    cardContent.value = p.profile_card ?? ''
  } catch { /* silent */ }
}

async function saveCard() {
  savingCard.value = true
  cardMsg.value = ''
  cardErr.value = ''
  try {
    const res = await rawFetch('/api/v1/me/profile-card', {
      method: 'PUT',
      body: { content: cardContent.value }
    }) as { message?: string }
    cardMsg.value = res.message ?? '卡片已保存'
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string, message?: string } }
    cardErr.value = err.data?.statusMessage ?? err.data?.message ?? '保存失败'
  } finally {
    savingCard.value = false
  }
}

async function savePassword() {
  pwMsg.value = ''
  pwErr.value = ''
  if (newPassword.value.length < 8) {
    pwErr.value = '新密码至少 8 位'
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    pwErr.value = '两次输入的新密码不一致'
    return
  }
  savingPw.value = true
  try {
    const res = await rawFetch('/api/v1/me/password', {
      method: 'PUT',
      body: { current_password: currentPassword.value, new_password: newPassword.value }
    }) as { message?: string }
    pwMsg.value = res.message ?? '密码已修改'
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (e: unknown) {
    const err = e as { data?: { message?: string } }
    pwErr.value = err.data?.message ?? '修改失败'
  } finally {
    savingPw.value = false
  }
}

onMounted(() => {
  void refreshAuth().then(() => {
    if (me.value) {
      void loadProfile()
      void loadCard()
    }
  })
})

useHead({ title: '个人设置' })
</script>

<template>
  <div class="mx-auto max-w-3xl space-y-6">
    <h1 class="text-2xl font-bold text-gray-900">
      个人设置
    </h1>

    <template v-if="me">
      <!-- profile -->
      <div class="rounded-xl border bg-white p-6">
        <h2 class="text-base font-semibold text-gray-900">
          个人信息
        </h2>
        <div class="mt-4 space-y-4">
          <div>
            <label class="mb-1 block text-sm text-gray-600">用户名</label>
            <input
              v-model="username"
              type="text"
              maxlength="32"
              class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
            >
          </div>
          <div>
            <label class="mb-1 block text-sm text-gray-600">邮箱</label>
            <input
              :value="email"
              type="email"
              disabled
              class="h-9 w-full cursor-not-allowed rounded-md border border-gray-200 bg-gray-50 px-3 text-sm text-gray-400"
            >
            <p class="mt-0.5 text-xs text-gray-400">
              邮箱暂不支持自助修改
            </p>
          </div>
          <div>
            <label class="mb-1 block text-sm text-gray-600">个人简介</label>
            <textarea
              v-model="bio"
              :rows="3"
              maxlength="500"
              class="w-full rounded-md border border-gray-200 px-3 py-2 text-sm focus:border-blue-400 focus:outline-none"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm text-gray-600">个人网站</label>
            <input
              v-model="website"
              type="url"
              placeholder="https://example.com"
              class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
            >
          </div>
          <div class="md:col-span-2">
            <label class="mb-1 block text-sm text-gray-600">回复签名</label>
            <p class="mb-1 text-xs text-gray-400">
              展示在你的每条回复卡片下方；支持 Markdown 文本与链接，不支持图片和脚本
            </p>
            <input
              v-model="signature"
              type="text"
              maxlength="200"
              placeholder="如：专注 Go 后端 · [我的博客](https://example.com)"
              class="w-full rounded-md border border-gray-200 px-3 py-2 text-sm focus:border-blue-400 focus:outline-none"
            >
            <p
              v-if="signature"
              class="mt-2 rounded-md bg-gray-50 px-3 py-2 text-xs text-gray-500"
            >
              预览：<span v-html="renderMarkdown(signature)" />
            </p>
          </div>
          <div class="flex items-center justify-between">
            <span
              v-if="profileMsg"
              class="text-xs text-green-600"
            >{{ profileMsg }}</span>
            <span
              v-else-if="profileErr"
              class="text-xs text-red-600"
            >{{ profileErr }}</span>
            <button
              type="button"
              :disabled="savingProfile"
              class="ml-auto rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              @click="saveProfile"
            >
              {{ savingProfile ? '保存中…' : '保存资料' }}
            </button>
          </div>
        </div>
      </div>

      <!-- profile card (markdown) -->
      <div class="rounded-xl border bg-white p-6">
        <h2 class="text-base font-semibold text-gray-900">
          个人卡片
        </h2>
        <p class="mt-1 text-xs text-gray-400">
          用 Markdown 自由编写展示在个人主页的卡片，保存后所有访客可见；清空则隐藏卡片
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
          maxlength="10000"
          placeholder="## 关于我&#10;&#10;自我介绍、技能、作品链接…"
          class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm leading-relaxed focus:border-blue-400 focus:outline-none focus:ring-1 focus:ring-blue-400"
        />
        <div class="mt-2 flex items-center justify-between">
          <span class="text-xs text-gray-400">{{ cardContent.length }}/10000</span>
          <div class="flex items-center gap-3">
            <span
              v-if="cardMsg"
              class="text-xs text-green-600"
            >{{ cardMsg }}</span>
            <span
              v-else-if="cardErr"
              class="text-xs text-red-600"
            >{{ cardErr }}</span>
            <button
              type="button"
              :disabled="savingCard"
              class="rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              @click="saveCard"
            >
              {{ savingCard ? '保存中…' : '保存卡片' }}
            </button>
          </div>
        </div>
      </div>

      <!-- security -->
      <div class="rounded-xl border bg-white p-6">
        <h2 class="text-base font-semibold text-gray-900">
          账号安全
        </h2>
        <div class="mt-4 space-y-4">
          <div>
            <label class="mb-1 block text-sm text-gray-600">当前密码</label>
            <input
              v-model="currentPassword"
              type="password"
              autocomplete="current-password"
              class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
            >
          </div>
          <div class="grid gap-4 md:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm text-gray-600">新密码（至少 8 位）</label>
              <input
                v-model="newPassword"
                type="password"
                autocomplete="new-password"
                class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
              >
            </div>
            <div>
              <label class="mb-1 block text-sm text-gray-600">确认新密码</label>
              <input
                v-model="confirmPassword"
                type="password"
                autocomplete="new-password"
                class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
              >
            </div>
          </div>
          <div class="flex items-center justify-between">
            <span
              v-if="pwMsg"
              class="text-xs text-green-600"
            >{{ pwMsg }}</span>
            <span
              v-else-if="pwErr"
              class="text-xs text-red-600"
            >{{ pwErr }}</span>
            <button
              type="button"
              :disabled="savingPw || !currentPassword || !newPassword"
              class="ml-auto rounded-lg bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              @click="savePassword"
            >
              {{ savingPw ? '提交中…' : '修改密码' }}
            </button>
          </div>
        </div>
      </div>
    </template>

    <div
      v-else
      class="rounded-xl border bg-white p-10 text-center"
    >
      <p class="text-gray-500">
        请先登录后访问个人设置
      </p>
      <NuxtLink
        to="/login?redirect=/settings"
        class="mt-4 inline-block rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >登录</NuxtLink>
    </div>
  </div>
</template>
