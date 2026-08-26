// Reverse proxy for /api/v1/** → Goravel backend.
// Uses h3's built-in proxyRequest for safe HTTP forwarding. The httpOnly
// admin_session cookie never reaches browser JS, but is translated into an
// Authorization bearer header here so v1 endpoints can authenticate the
// caller (e.g. posting replies) without exposing the raw token.

const BACKEND_HOST = 'http://127.0.0.1:9000'

export default defineEventHandler(async (event) => {
  const reqPath = event.path
  if (!reqPath.startsWith('/api/v1/') || reqPath.includes('..')) {
    throw createError({ statusCode: 403, statusMessage: 'Forbidden' })
  }

  const token = getCookie(event, 'admin_session') ?? ''
  const headers = token ? { authorization: `Bearer ${token}` } : undefined

  return proxyRequest(event, BACKEND_HOST + reqPath, { headers })
})
