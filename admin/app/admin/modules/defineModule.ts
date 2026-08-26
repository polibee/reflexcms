import type { ModuleDef } from '../core/types'

/**
 * Identity helper for business modules.
 * A module bundles resources + widgets + nav groups into one
 * registrable unit - the Nuxt equivalent of a Laravel Service Provider.
 */
export function defineModule(module: ModuleDef): ModuleDef {
  return module
}
