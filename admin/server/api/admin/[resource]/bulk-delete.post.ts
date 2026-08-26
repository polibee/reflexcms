import { forward } from '../../../utils/backend'
import { getConfig } from '../../../utils/resourceConfigs'

export default defineEventHandler(async (event) => {
  const resource = getRouterParam(event, 'resource')!
  getConfig(resource)

  const body = await readBody<{ ids?: Array<number | string> }>(event)
  const ids = (body?.ids ?? []).map(Number).filter(n => Number.isInteger(n))

  if (ids.length === 0) {
    throw createError({ statusCode: 422, statusMessage: 'No ids provided' })
  }

  return await forward(event, 'POST', `/api/admin/${resource}/bulk-delete`, {
    body: { ids }
  })
})
