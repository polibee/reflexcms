import { forward } from '../../../utils/backend'
import { getConfig } from '../../../utils/resourceConfigs'

export default defineEventHandler(async (event) => {
  const resource = getRouterParam(event, 'resource')!
  getConfig(resource)

  const query = getQuery(event)
  return await forward(event, 'GET', `/api/admin/${resource}`, { query })
})
