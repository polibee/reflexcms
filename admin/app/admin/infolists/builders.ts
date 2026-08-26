import type { BadgeStyle, EntryNode, EntryKind } from '../core/types'

interface EntryOptions {
  prefix?: string
  suffix?: string
}

function entry(kind: EntryKind, name: string, label: string, opts?: EntryOptions): EntryNode {
  return { kind, name, label, ...opts }
}

export function textEntry(name: string, label: string, opts?: EntryOptions): EntryNode {
  return entry('text', name, label, opts)
}

export function badgeEntry(
  name: string,
  label: string,
  badges: Record<string | number, BadgeStyle>
): EntryNode {
  return { kind: 'badge', name, label, badges }
}

export function booleanEntry(name: string, label: string): EntryNode {
  return { kind: 'boolean', name, label }
}

export function moneyEntry(name: string, label: string, opts?: EntryOptions): EntryNode {
  return entry('money', name, label, { prefix: '$', ...opts })
}

export function dateEntry(name: string, label: string): EntryNode {
  return entry('date', name, label)
}

export function datetimeEntry(name: string, label: string): EntryNode {
  return entry('datetime', name, label)
}

export function linkEntry(name: string, label: string): EntryNode {
  return entry('link', name, label)
}
