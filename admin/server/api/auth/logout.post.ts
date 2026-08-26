import { clearSessionCookie } from '../../utils/auth'
import { forward } from '../../utils/backend'

export default defineEventHandler(async (event) => {
  try {
    await forward(event, 'POST', '/api/auth/logout')
  } catch {
    // backend session may already be gone; the cookie must go regardless
  } finally {
    clearSessionCookie(event)
  }
  return { ok: true }
})
