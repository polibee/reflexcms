import { forward } from '../../../utils/backend'

export default defineEventHandler(async (event) => {
  const body = await readBody<Record<string, unknown>>(event)
  return await forward(event, 'POST', '/api/admin/mail/send', { body: body ?? {} })
})
