<script setup lang="ts">
definePageMeta({ layout: false })

import { UserPlusIcon } from 'lucide-vue-next'

const auth = useAuthStore()
const route = useRoute()

const username = ref('')
const email = ref('')
const password = ref('')
const inviteCode = ref('')
const inviteRequired = ref(false)
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  try {
    const cfg = await $fetch<{ sections: string[] }>('/api/v1/site/config').catch(() => null)
    // invite requirement is read from the registration settings; probe via
    // the registration config endpoint
    const reg = await $fetch<{ invite_required: boolean }>('/api/v1/auth/register-config').catch(() => null)
    inviteRequired.value = Boolean(reg?.invite_required)
  } catch { /* silent */ }
})

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await $fetch('/api/v1/auth/register', {
      method: 'POST',
      body: {
        username: username.value.trim(),
        email: email.value.trim(),
        password: password.value,
        invite_code: inviteCode.value.trim()
      }
    })
    // Sign in with credentials (BFF sets the shared session cookie)
    await auth.login(email.value.trim(), password.value)
    const target = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    const safe = target.startsWith('/') && !target.startsWith('//') ? target : '/'
    await navigateTo(safe, { replace: true })
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string, message?: string } }
    error.value = err.data?.statusMessage ?? err.data?.message ?? '注册失败'
  } finally {
    loading.value = false
  }
}

useHead({ title: '注册' })
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 p-4">
    <div class="w-full max-w-sm space-y-6">
      <div class="space-y-2 text-center">
        <div class="mx-auto flex h-11 w-11 items-center justify-center rounded-xl bg-blue-600 text-white">
          <UserPlusIcon class="h-5 w-5" />
        </div>
        <h1 class="text-xl font-semibold text-gray-900">
          注册账号
        </h1>
        <p class="text-sm text-gray-500">
          加入社区，参与讨论
        </p>
      </div>

      <form
        class="space-y-4 rounded-xl border bg-white p-6"
        @submit.prevent="submit"
      >
        <div>
          <label class="mb-1 block text-sm text-gray-600">用户名</label>
          <input
            v-model="username"
            type="text"
            required
            maxlength="32"
            class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
          >
        </div>
        <div>
          <label class="mb-1 block text-sm text-gray-600">邮箱</label>
          <input
            v-model="email"
            type="email"
            required
            class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
          >
        </div>
        <div>
          <label class="mb-1 block text-sm text-gray-600">密码（至少 8 位）</label>
          <input
            v-model="password"
            type="password"
            required
            minlength="8"
            class="h-9 w-full rounded-md border border-gray-200 px-3 text-sm focus:border-blue-400 focus:outline-none"
          >
        </div>
        <div v-if="inviteRequired">
          <label class="mb-1 block text-sm text-gray-600">邀请码 <span class="text-red-500">*</span></label>
          <input
            v-model="inviteCode"
            type="text"
            required
            placeholder="INV-XXXXXXXXXX"
            class="h-9 w-full rounded-md border border-gray-200 px-3 font-mono text-sm focus:border-blue-400 focus:outline-none"
          >
          <p class="mt-0.5 text-xs text-gray-400">
            本站为邀请制注册；没有邀请码？<NuxtLink to="/shop" class="text-blue-600 hover:underline">前往商城获取</NuxtLink>
          </p>
        </div>

        <p
          v-if="error"
          class="text-xs text-red-600"
        >
          {{ error }}
        </p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full rounded-lg bg-blue-600 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {{ loading ? '注册中…' : '注册' }}
        </button>
        <p class="text-center text-xs text-gray-400">
          已有账号？
          <NuxtLink
            :to="`/login?redirect=${route.query.redirect ?? '/'}`"
            class="text-blue-600 hover:underline"
          >直接登录</NuxtLink>
        </p>
      </form>
    </div>
  </div>
</template>
