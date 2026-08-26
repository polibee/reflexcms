import type { ActionDef, BadgeStyle, ColumnDefLite, ColumnMeta } from '../../core/types'

interface ColumnOptions {
  sortable?: boolean
  align?: ColumnMeta['align']
  prefix?: string
  suffix?: string
}

export function textColumn(name: string, label: string, opts?: ColumnOptions): ColumnDefLite {
  return { name, label, meta: { kind: 'text', ...opts } }
}

export function numberColumn(name: string, label: string, opts?: ColumnOptions): ColumnDefLite {
  return { name, label, meta: { kind: 'number', sortable: true, ...opts } }
}

export function moneyColumn(name: string, label: string, opts?: ColumnOptions): ColumnDefLite {
  return { name, label, meta: { kind: 'money', sortable: true, prefix: '$', ...opts } }
}

export function dateColumn(name: string, label: string, opts?: ColumnOptions): ColumnDefLite {
  return { name, label, meta: { kind: 'date', ...opts } }
}

export function badgeColumn(
  name: string,
  label: string,
  badges: Record<string | number, BadgeStyle>,
  opts?: ColumnOptions
): ColumnDefLite {
  return { name, label, meta: { kind: 'badge', badges, ...opts } }
}

export function booleanColumn(
  name: string,
  label: string,
  labels = { true: 'Yes', false: 'No' },
  opts?: ColumnOptions
): ColumnDefLite {
  return { name, label, meta: { kind: 'boolean', booleanLabels: labels, ...opts } }
}

/** row actions column (rendered right-aligned, no label) */
export function actionsColumn(actions: ActionDef[]): ColumnDefLite {
  return { name: '__actions', label: '', meta: { kind: 'actions', actions, align: 'right' } }
}
