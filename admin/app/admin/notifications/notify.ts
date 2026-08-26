interface Toast {
  id: number
  title: string
  description?: string
  variant: 'success' | 'error' | 'info'
}

export const useToasts = defineStore('admin-toasts', () => {
  const items = ref<Toast[]>([])
  let seq = 0

  function push(t: Omit<Toast, 'id'>): void {
    const id = ++seq
    items.value.push({ ...t, id })
    if (import.meta.client) {
      setTimeout(() => dismiss(id), t.variant === 'error' ? 6000 : 3500)
    }
  }

  function dismiss(id: number): void {
    items.value = items.value.filter(t => t.id !== id)
  }

  return { items, push, dismiss }
})

export function notify(title: string, description?: string, variant: Toast['variant'] = 'success'): void {
  useToasts().push({ title, description, variant })
}

export function notifyError(title: string, description?: string): void {
  notify(title, description, 'error')
}
