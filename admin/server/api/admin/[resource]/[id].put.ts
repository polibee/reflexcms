import { forward } from '../../../utils/backend'
import { getConfig } from '../../../utils/resourceConfigs'

export default defineEventHandler(async (event) => {
  const resource = getRouterParam(event, 'resource')!
  getConfig(resource)
  const id = getRouterParam(event, 'id')!

  const body = await readBody<Record<string, unknown>>(event)
  return await forward(event, 'PUT', `/api/admin/${resource}/${id}`, { body: body ?? {} })
})
