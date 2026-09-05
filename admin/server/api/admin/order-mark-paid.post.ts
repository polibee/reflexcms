// Flat path on purpose: a static /api/admin/orders/... segment would shadow
// the dynamic [resource] routes in Nitro (same class of conflict as the gin
// static/wildcard clash).
import { forward } from '../../utils/backend'

export default defineEventHandler(async (event) => {
  const body = await readBody<{ order_id?: number }>(event)
  if (!body?.order_id) {
    throw createError({ statusCode: 422, statusMessage: 'order_id is required' })
  }
  return await forward(event, 'POST', `/api/admin/orders/${body.order_id}/mark-paid`, { body: {} })
})
