import type { Paginated } from '#shared/types/api'
import type { ColumnDefLite, ResolvedResource } from '../core/types'
import { notifyError } from '../notifications/notify'

export interface RowRecord extends Record<string, unknown> {
  id: number | string
}

/**
 * Server-driven table engine: search / sort / paginate / selection / bulk.
 * One source of truth consumed by <DataTable>.
 */
export function useResourceTable(resource: ResolvedResource, options?: { defaultSort?: string }) {
  const panel = getPanel()

  const q = ref('')
  const page = ref(1)
  const perPage = ref(panel.perPage ?? 10)
  const sortBy = ref(options?.defaultSort ?? '')
  const sortDir = ref<'asc' | 'desc'>('asc')

  const items = ref<RowRecord[]>([])
  const total = ref(0)
  const totalPages = ref(1)
  const loading = ref(false)

  const selection = ref(new Set<string | number>())
  const selectedIds = computed(() => [...selection.value])

  async function fetch(): Promise<void> {
    loading.value = true
    try {
      const res = await $fetch<Paginated<RowRecord>>(`/api/admin/${resource.name}`, {
        query: {
          q: q.value || undefined,
          page: page.value,
          perPage: perPage.value,
          sortBy: sortBy.value || undefined,
          sortDir: sortDir.value
        }
      })
      items.value = res.items
      total.value = res.total
      totalPages.value = Math.max(res.totalPages, 1)
    } catch (e: unknown) {
      notifyError(`Failed to load ${resource.labelPlural ?? resource.label}`, (e as Error).message)
    } finally {
      loading.value = false
    }
  }

  // debounce search input
  let timer: ReturnType<typeof setTimeout> | null = null
  watch(q, () => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      page.value = 1
      fetch()
    }, 300)
  })

  watch([page, perPage, sortBy, sortDir], fetch)
  onMounted(fetch)

  function toggleSort(column: ColumnDefLite): void {
    if (!column.meta.sortable) return
    if (sortBy.value === column.name) {
      sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortBy.value = column.name
      sortDir.value = 'asc'
    }
  }

  /* ---------------- selection ---------------- */

  function isSelected(id: string | number): boolean {
    return selection.value.has(id)
  }

  function toggleSelect(row: RowRecord): void {
    const next = new Set(selection.value)
    if (next.has(row.id)) next.delete(row.id)
    else next.add(row.id)
    selection.value = next
  }

  const allPageSelected = computed(
    () => items.value.length > 0 && items.value.every(r => selection.value.has(r.id))
  )

  function toggleSelectAll(): void {
    const next = new Set(selection.value)
    if (allPageSelected.value) items.value.forEach(r => next.delete(r.id))
    else items.value.forEach(r => next.add(r.id))
    selection.value = next
  }

  function clearSelection(): void {
    selection.value = new Set()
  }

  /* mutation API - components call these instead of writing refs directly */
  function setQuery(term: string): void {
    q.value = term
  }

  function setPage(p: number): void {
    page.value = Math.min(Math.max(p, 1), totalPages.value)
  }

  function nextPage(): void {
    setPage(page.value + 1)
  }

  function prevPage(): void {
    setPage(page.value - 1)
  }

  function setPerPage(n: number): void {
    perPage.value = n
    page.value = 1
  }

  async function refresh(): Promise<void> {
    await fetch()
    clearSelection()
  }

  return {
    q, page, perPage, sortBy, sortDir,
    items, total, totalPages, loading,
    selection, selectedIds,
    fetch, refresh,
    toggleSort, toggleSelect, toggleSelectAll, isSelected, allPageSelected, clearSelection,
    setQuery, setPage, nextPage, prevPage, setPerPage
  }
}

export type ResourceTable = ReturnType<typeof useResourceTable>
