/* =============================================================
 * BFF proxy core: every admin request is forwarded path-faithful to the
 * Goravel backend (plan §2.1 / ADR-002). The opaque session token lives in
 * our httpOnly cookie and travels to the backend as a bearer token.
 * ============================================================= */

export const SESSION_COOKIE = 'admin_session'

type H3Evt = Parameters<typeof getCookie>[0]

interface BackendErrorShape {
  statusCode?: number
  statusMessage?: string
}

interface FetchErrorLike extends Error {
  statusCode?: number
  statusMessage?: string
  data?: BackendErrorShape | unknown
}

export function backendUrl(): string {
  return useRuntimeConfig().backendUrl as string
}

/**
 * Forward a request to the backend, attaching the caller's session token.
 * Non-2xx responses are rethrown as h3 errors with the backend's original
 * status code and message, so 401/403/404/422 semantics survive the hop.
 */
// Typed $fetch explodes into deep route-matching generics for dynamic URLs;
// an untyped alias keeps vue-tsc sane while behaviour stays identical.
const rawFetch = $fetch as unknown as (
  url: string,
  opts?: Record<string, unknown>
) => Promise<unknown>

export async function forward<T>(
  event: H3Evt,
  method: 'GET' | 'POST' | 'PUT' | 'DELETE',
  path: string,
  opts?: { body?: Record<string, unknown>, query?: Record<string, unknown> }
): Promise<T> {
  const token = getCookie(event, SESSION_COOKIE)

  try {
    return await rawFetch(`${backendUrl()}${path}`, {
      method,
      headers: token ? { authorization: `Bearer ${token}` } : {},
      ...(opts?.body !== undefined ? { body: opts.body } : {}),
      ...(opts?.query !== undefined ? { query: opts.query } : {})
    }) as T
  } catch (err) {
    const e = err as FetchErrorLike
    const data = (e.data ?? {}) as BackendErrorShape
    const status = e.statusCode ?? data.statusCode ?? 502
    const message = data.statusMessage ?? e.statusMessage ?? 'Backend unavailable'
    throw createError({ statusCode: status, statusMessage: message })
  }
}
