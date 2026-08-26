export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuthStore()

  if (!to.path.startsWith('/admin')) return

  if (!auth.user) await auth.fetchMe()
  if (!auth.user) {
    return navigateTo(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
  }
})
