/* Shared signed-in state for the public frontend. One /api/v1/me probe is
 * shared app-wide via useState; refresh() re-validates (login, logout,
 * window focus) so every consumer — sidebar user card, reply composer,
 * comment form — sees the same identity without duplicate requests. */

export interface PublicUser {
  id: number
  name: string
  email: string
}

export function useAuthUser() {
  const user = useState<PublicUser | null>('auth-user', () => null)
  let pending: Promise<void> | null = null

  async function refresh(): Promise<void> {
    if (pending) return pending
    pending = (async () => {
      try {
        const res = await $fetch<{ user: PublicUser | null }>('/api/v1/me')
        user.value = res.user ?? null
      } catch {
        user.value = null
      }
    })()
    try {
      await pending
    } finally {
      pending = null
    }
  }

  return { user, refresh }
}
