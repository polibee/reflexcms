import type { AdminResource } from './types'

/** Identity helper giving full type inference to resource definitions. */
export function defineResource(resource: AdminResource): AdminResource {
  return resource
}
