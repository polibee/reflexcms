import type { Component } from 'vue'
import type { ZodTypeAny } from 'zod'

/* =============================================================
 * Admin Framework - Core Type System
 * Design principle: everything is declarative data (schema),
 * rendered by generic engines (Form/Table/Infolist).
 * ============================================================= */

/* ---------------- Schema: Fields ---------------- */

export type FieldType
  = | 'text'
    | 'email'
    | 'password'
    | 'number'
    | 'textarea'
    | 'richtext'
    | 'tiptap'
    | 'wysiwyg'
    | 'select'
    | 'switch'
    | 'checkbox'
    | 'checkboxgroup'
    | 'date'
    | 'relation'
    | 'upload'

export interface FieldOption {
  label: string
  value: string | number
}

export interface RelationConfig {
  /** target resource slug, e.g. "users" */
  resource: string
  /** field of related record used as option label */
  labelKey: string
}

export interface FieldNode {
  type: 'field'
  kind: FieldType
  name: string
  label: string
  placeholder?: string
  required?: boolean
  disabled?: boolean
  helpText?: string
  defaultValue?: unknown
  /** select options */
  options?: FieldOption[]
  /** number constraints */
  min?: number
  max?: number
  step?: number
  /** textarea rows */
  rows?: number
  /** relation picker */
  relation?: RelationConfig
  /** extra zod rules appended to base rule */
  rules?: ZodTypeAny[]
  /** grid column span */
  colSpan?: 1 | 2 | 3 | 4
  /** file input accept attribute (upload kind) */
  accept?: string
  /** render only when this predicate over the current form values holds */
  visibleIf?: (values: Record<string, unknown>) => boolean
}

/* ---------------- Schema: Layouts ---------------- */

export interface SectionNode {
  type: 'section'
  title?: string
  description?: string
  children: SchemaNode[]
}

export interface GridNode {
  type: 'grid'
  columns: 1 | 2 | 3 | 4
  children: SchemaNode[]
}

export type LayoutNode = SectionNode | GridNode
export type SchemaNode = FieldNode | LayoutNode

/* ---------------- Table Columns ---------------- */

export type BadgeVariant = 'default' | 'success' | 'warning' | 'destructive' | 'secondary'

export interface BadgeStyle {
  label: string
  variant: BadgeVariant
}

export interface RowActionContext {
  record: Record<string, unknown>
  resource: ResolvedResource
}

export interface BulkActionContext {
  ids: Array<string | number>
  resource: ResolvedResource
}

/** unified execution context passed to action handlers.
 *  resource is always provided by the runner; record/ids depend on origin. */
export interface ActionContext extends Partial<Omit<RowActionContext, 'resource'>>, Partial<Omit<BulkActionContext, 'resource'>> {
  resource: ResolvedResource
  /** values from the action's modal form, if it defines one */
  values?: Record<string, unknown>
}

export interface ActionDef {
  name: string
  label: string
  icon?: string
  variant?: 'default' | 'outline' | 'destructive' | 'ghost'
  permission?: string
  confirm?: {
    title: string
    description?: string
    confirmLabel?: string
  }
  /** optional modal form schema resolved before handler runs */
  form?: () => SchemaNode[]
  handler?: (ctx: ActionContext) => Promise<void> | void
  visible?: (record: Record<string, unknown>) => boolean
}

export interface ColumnMeta {
  kind: 'text' | 'number' | 'money' | 'date' | 'badge' | 'boolean' | 'actions'
  sortable?: boolean
  align?: 'left' | 'center' | 'right'
  badges?: Record<string | number, BadgeStyle>
  booleanLabels?: { true: string, false: string }
  /** kind === 'actions': row actions */
  actions?: ActionDef[]
  prefix?: string
  suffix?: string
}

export interface ColumnDefLite {
  name: string
  label: string
  meta: ColumnMeta
}

/* ---------------- Infolist Entries ---------------- */

export type EntryKind = 'text' | 'badge' | 'boolean' | 'money' | 'date' | 'datetime' | 'link'

export interface EntryNode {
  kind: EntryKind
  name: string
  label: string
  badges?: Record<string | number, BadgeStyle>
  prefix?: string
  suffix?: string
}

/* ---------------- Resource ---------------- */

export interface AdminResource {
  /** url slug + api segment, e.g. "users" */
  name: string
  model?: string
  label: string
  labelPlural?: string
  /** lucide icon key, see utils/iconMap */
  icon?: string
  /** navigation group label */
  group?: string
  sort?: number
  /** permission prefix, defaults to name */
  permissionPrefix?: string
  searchable?: string[]
  form?: () => SchemaNode[]
  table?: () => ColumnDefLite[]
  infolist?: () => EntryNode[]
  rowActions?: ActionDef[]
  bulkActions?: ActionDef[]
  /** API record → form values (e.g. unpack a JSON config column into fields) */
  transformIn?: (record: Record<string, unknown>) => Record<string, unknown>
  /** form values → API payload (e.g. pack fields back into a JSON column) */
  transformOut?: (values: Record<string, unknown>) => Record<string, unknown>
  /** initial values for the CREATE form (e.g. sort = current max + 1) */
  defaultValues?: () => Record<string, unknown> | Promise<Record<string, unknown>>
}

export type ResolvedResource = Omit<AdminResource, 'labelPlural' | 'permissionPrefix'>
  & Required<Pick<AdminResource, 'labelPlural' | 'permissionPrefix'>>

/* ---------------- Widgets ---------------- */

export interface WidgetDef {
  name: string
  label?: string
  description?: string
  /** grid span out of 4 */
  span?: 1 | 2 | 3 | 4
  order?: number
  permission?: string
  component: Component
}

/* ---------------- Navigation ---------------- */

export interface NavItem {
  label: string
  to: string
  icon?: string
  badge?: string | number
  permission?: string
}

export interface NavGroup {
  label: string
  sort?: number
  items: NavItem[]
}

/* ---------------- Panel ---------------- */

export interface PanelConfig {
  id: string
  path: string
  title: string
  branding?: { name: string, logo?: string }
  auth?: boolean
  /** default items per page for tables */
  perPage?: number
}

/* ---------------- Module ---------------- */

export interface ModuleDef {
  name: string
  resources?: AdminResource[]
  widgets?: WidgetDef[]
  navGroups?: Array<{ label: string, sort?: number, items?: NavItem[] }>
}
