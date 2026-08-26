import { forward } from '../../../utils/backend'
import { getConfig } from '../../../utils/resourceConfigs'

export default defineEventHandler(async (event) => {
  const resource = getRouterParam(event, 'resource')!
  getConfig(resource)
  const id = getRouterParam(event, 'id')!

  return await forward(event, 'DELETE', `/api/admin/${resource}/${id}`)
})
