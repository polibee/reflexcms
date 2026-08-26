import type { ActionContext, ActionDef } from '../core/types'
import { notifyError } from '../notifications/notify'

/**
 * Central action pipeline: permission check -> optional modal form ->
 * optional confirmation -> handler execution -> error notification.
 * State is consumed by the single <ActionHost> mounted in the panel layout.
 */
export const useActionRunner = defineStore('admin-action-runner', () => {
  const confirmTarget = ref<{ action: ActionDef, ctx: ActionContext } | null>(null)
  const formTarget = ref<{ action: ActionDef, ctx: ActionContext } | null>(null)
  const busy = ref(false)

  async function run(action: ActionDef, ctx: ActionContext): Promise<void> {
    if (action.permission && !can(action.permission, useAuthStore().user)) return
    if (action.form) {
      formTarget.value = { action, ctx }
      return
    }
    if (action.confirm) {
      confirmTarget.value = { action, ctx }
      return
    }
    await execute(action, ctx)
  }

  async function execute(action: ActionDef, ctx: ActionContext, values?: Record<string, unknown>): Promise<void> {
    busy.value = true
    try {
      await action.handler?.({ ...ctx, values })
      close()
    } catch (e: unknown) {
      notifyError(`Action "${action.label}" failed`, (e as Error)?.message)
    } finally {
      busy.value = false
    }
  }

  function close(): void {
    confirmTarget.value = null
    formTarget.value = null
  }

  return { confirmTarget, formTarget, busy, run, execute, close }
})
