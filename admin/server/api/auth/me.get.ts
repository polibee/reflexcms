import type { AuthUser } from '#shared/types/api'
import { forward } from '../../utils/backend'
import { requireSessionToken } from '../../utils/auth'

export default defineEventHandler(async (event) => {
  requireSessionToken(event)
  return await forward<AuthUser>(event, 'GET', '/api/auth/me')
})
