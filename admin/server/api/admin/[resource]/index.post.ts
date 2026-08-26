import { forward } from '../../../utils/backend'
import { getConfig } from '../../../utils/resourceConfigs'

export default defineEventHandler(async (event) => {
  const resource = getRouterParam(event, 'resource')!
  getConfig(resource)

  const body = await readBody<Record<string, unknown>>(event)
  return await forward(event, 'POST', `/api/admin/${resource}`, { body: body ?? {} })
})
