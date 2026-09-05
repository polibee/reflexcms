/* Public-site UI singletons: toast notifications and confirm dialogs,
 * replacing window.alert/window.confirm with modern in-app surfaces.
 * Host component: FrontToastHost (mounted once in the public layout). */

export interface FrontToast {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

export interface FrontConfirmOptions {
  title: string
  message?: string
  confirmLabel?: string
  cancelLabel?: string
  danger?: boolean
}

interface FrontConfirmState extends FrontConfirmOptions {
  resolve: (value: boolean) => void
}

const toastSeq = ref(0)
const toasts = ref<FrontToast[]>([])
const confirmState = ref<FrontConfirmState | null>(null)

export function useFrontUi() {
  function toast(message: string, type: FrontToast['type'] = 'success') {
    const id = ++toastSeq.value
    toasts.value.push({ id, message, type })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, 3200)
  }

  function confirmDialog(opts: FrontConfirmOptions): Promise<boolean> {
    return new Promise((resolve) => {
      confirmState.value = { ...opts, resolve }
    })
  }

  function settleConfirm(value: boolean) {
    confirmState.value?.resolve(value)
    confirmState.value = null
  }

  return { toasts, confirmState, toast, confirmDialog, settleConfirm }
}
