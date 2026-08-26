export function formatMoney(value: unknown, prefix = '$'): string {
  const n = Number(value ?? 0)
  if (Number.isNaN(n)) return '-'
  return `${prefix}${n.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

export function formatNumber(value: unknown): string {
  const n = Number(value)
  if (Number.isNaN(n)) return '-'
  return n.toLocaleString('en-US')
}

export function formatDate(value: unknown): string {
  if (!value) return '-'
  const d = new Date(String(value))
  if (Number.isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

export function formatDateTime(value: unknown): string {
  if (!value) return '-'
  const d = new Date(String(value))
  if (Number.isNaN(d.getTime())) return '-'
  return d.toLocaleString('en-US', {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

/** client-side CSV export of the current table page */
export function exportCsv(
  filename: string,
  rows: Array<Record<string, unknown>>,
  columns: Array<{ name: string, label: string }>
): void {
  if (import.meta.server || rows.length === 0) return
  const cols = columns.filter(c => c.name !== '__actions')
  const escapeCell = (v: unknown) => `"${String(v ?? '').replaceAll('"', '""')}"`
  const lines = [
    cols.map(c => escapeCell(c.label)).join(','),
    ...rows.map(row => cols.map(c => escapeCell(row[c.name])).join(','))
  ]
  const blob = new Blob([lines.join('\n')], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${filename}-${new Date().toISOString().slice(0, 10)}.csv`
  a.click()
  URL.revokeObjectURL(url)
}
