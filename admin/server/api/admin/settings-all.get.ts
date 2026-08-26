import { forward } from '../../utils/backend'
import { requireSessionToken } from '../../utils/auth'

export default defineEventHandler(async (event) => {
  requireSessionToken(event)
  return await forward(event, 'GET', '/api/admin/settings-all')
})
