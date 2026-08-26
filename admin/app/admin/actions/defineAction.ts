import type { ActionDef } from '../core/types'

/** Identity helper giving full type inference to action definitions. */
export function defineAction(action: ActionDef): ActionDef {
  return action
}
