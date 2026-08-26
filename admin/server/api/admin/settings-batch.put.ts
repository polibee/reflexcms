import { forward } from '../../utils/backend'
import { requireSessionToken } from '../../utils/auth'

export default defineEventHandler(async (event) => {
  requireSessionToken(event)
  const body = await readBody<Record<string, Record<string, unknown>>>(event)
  return await forward(event, 'PUT', '/api/admin/settings-batch', { body: body ?? {} })
})
