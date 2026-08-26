/* =============================================================
 * Session helpers. The session itself lives in the backend (opaque token
 * hashed into admin_sessions); we only manage the httpOnly cookie that
 * carries the raw token to our proxy layer (plan §5.1).
 * ============================================================= */

type H3Evt = Parameters<typeof getCookie>[0]

export function getSessionToken(event: H3Evt): string {
  return getCookie(event, 'admin_session') ?? ''
}

/** Cheap local gate: real authentication/authorization happens backend-side
 *  and its 401/403 pass straight through the proxy (defence at the boundary). */
export function requireSessionToken(event: H3Evt): string {
  const token = getSessionToken(event)
  if (!token) {
    throw createError({ statusCode: 401, statusMessage: 'Unauthorized' })
  }
  return token
}

export function setSessionCookie(event: H3Evt, token: string): void {
  setCookie(event, 'admin_session', token, {
    httpOnly: true,
    sameSite: 'lax',
    secure: Boolean(useRuntimeConfig().cookieSecure),
    path: '/',
    maxAge: 60 * 60 * 8
  })
}

export function clearSessionCookie(event: H3Evt): void {
  deleteCookie(event, 'admin_session', { path: '/' })
}
