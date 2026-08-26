import type { AuthUser } from '#shared/types/api'

export const useAuthStore = defineStore('admin-auth', () => {
  const user = ref<AuthUser | null>(null)

  async function fetchMe(): Promise<AuthUser | null> {
    try {
      const headers = import.meta.server ? useRequestHeaders(['cookie']) : undefined
      user.value = await $fetch<AuthUser>('/api/auth/me', { headers })
    } catch {
      user.value = null
    }
    return user.value
  }

  async function login(email: string, password: string): Promise<void> {
    await $fetch('/api/auth/login', { method: 'POST', body: { email, password } })
    await fetchMe()
  }

  async function logout(): Promise<void> {
    try {
      await $fetch('/api/auth/logout', { method: 'POST' })
    } finally {
      user.value = null
      await navigateTo('/login')
    }
  }

  return { user, login, logout, fetchMe }
})
