<script setup lang="ts">
import { EyeIcon, EyeOffIcon, LogInIcon } from 'lucide-vue-next'

definePageMeta({ layout: false })

const auth = useAuthStore()
const route = useRoute()

const email = ref('admin@reflexcms.dev')
const password = ref('ReflexCMS@2026')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

async function submit(): Promise<void> {
  loading.value = true
  error.value = ''
  try {
    await auth.login(email.value, password.value)
    const target = typeof route.query.redirect === 'string' ? route.query.redirect : ''
    // only allow in-app relative redirects (open-redirect guard)
    const safe = target.startsWith('/') && !target.startsWith('//') ? target : '/admin'
    await navigateTo(safe, { replace: true })
  } catch (e: unknown) {
    error.value = (e as { data?: { message?: string } })?.data?.message ?? 'Invalid credentials'
  } finally {
    loading.value = false
  }
}

function fill(account: string): void {
  email.value = account
  password.value = 'ReflexCMS@2026'
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-muted/40 p-4 dark:bg-background">
    <div class="w-full max-w-sm space-y-6">
      <div class="space-y-2 text-center">
        <div class="mx-auto flex h-11 w-11 items-center justify-center rounded-xl bg-primary text-primary-foreground">
          <svg
            viewBox="0 0 24 24"
            class="h-5 w-5"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
          >
            <path
              d="M4 17l6-6-6-6M12 19h8"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </div>
        <h1 class="text-xl font-semibold">
          Nuxt Admin
        </h1>
        <p class="text-sm text-muted-foreground">
          Sign in to your panel
        </p>
      </div>

      <UiCard class="p-6">
        <form
          class="space-y-4"
          @submit.prevent="submit"
        >
          <div class="space-y-1.5">
            <UiLabel for="email">
              Email
            </UiLabel>
            <UiInput
              id="email"
              v-model="email"
              type="email"
              placeholder="you@example.com"
              required
            />
          </div>
          <div class="space-y-1.5">
            <UiLabel for="password">
              Password
            </UiLabel>
            <div class="relative">
              <UiInput
                id="password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="••••••••"
                required
                class="pr-9"
              />
              <button
                type="button"
                class="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                @click="showPassword = !showPassword"
              >
                <EyeIcon
                  v-if="showPassword"
                  class="h-4 w-4"
                />
                <EyeOffIcon
                  v-else
                  class="h-4 w-4"
                />
              </button>
            </div>
          </div>

          <p
            v-if="error"
            class="text-xs text-destructive"
          >
            {{ error }}
          </p>

          <UiButton
            type="submit"
            class="w-full"
            :disabled="loading"
          >
            <LogInIcon /> {{ loading ? 'Signing in…' : 'Sign In' }}
          </UiButton>
        </form>
      </UiCard>

      <!-- demo accounts -->
      <UiCard class="p-4">
        <p class="mb-2 text-xs font-medium text-muted-foreground">
          演示账号（密码 ReflexCMS@2026）
        </p>
        <div class="grid gap-1 text-xs">
          <button
            class="flex justify-between rounded px-2 py-1 hover:bg-accent"
            type="button"
            @click="fill('admin@reflexcms.dev')"
          >
            <span>admin@reflexcms.dev</span><span class="text-muted-foreground">超级管理员</span>
          </button>
        </div>
      </UiCard>
    </div>
  </div>
</template>
