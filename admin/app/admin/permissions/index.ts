/* =============================================================
 * Permission system: User -> Role -> Permission (wildcard aware)
 * e.g. "users.view", "posts.*", "*" grants everything.
 * ============================================================= */

export function can(permission: string | undefined, user?: { permissions: string[] } | null): boolean {
  if (!permission) return true
  if (!user) return false
  return user.permissions.some((granted) => {
    if (granted === '*' || granted === permission) return true
    if (granted.endsWith('.*')) {
      const prefix = granted.slice(0, -1) // keep trailing dot: "posts."
      return permission.startsWith(prefix)
    }
    return false
  })
}

export function useCan(): (permission?: string) => boolean {
  const auth = useAuthStore()
  return (permission?: string) => can(permission, auth.user)
}
