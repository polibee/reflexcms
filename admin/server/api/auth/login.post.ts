import type { AuthUser } from '#shared/types/api'
import { backendUrl } from '../../utils/backend'
import { setSessionCookie } from '../../utils/auth'

export default defineEventHandler(async (event) => {
  const body = await readBody<{ email?: string, password?: string }>(event)

  if (!body?.email || !body?.password) {
    throw createError({ statusCode: 400, statusMessage: 'Email and password are required' })
  }

  const res = await $fetch<{ token: string, user: AuthUser }>(
    `${backendUrl()}/api/auth/login`,
    { method: 'POST', body: { email: body.email, password: body.password } }
  ).catch((err) => {
    const e = err as { statusCode?: number, data?: { statusMessage?: string } }
    throw createError({
      statusCode: e.statusCode ?? 502,
      statusMessage: e.data?.statusMessage ?? 'Backend unavailable'
    })
  })

  setSessionCookie(event, res.token)
  return { user: res.user }
})
