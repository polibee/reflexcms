import type { FieldNode, FieldOption, FieldType, RelationConfig } from '../../core/types'

interface FieldOptions {
  placeholder?: string
  required?: boolean
  disabled?: boolean
  helpText?: string
  defaultValue?: unknown
  options?: FieldOption[]
  min?: number
  max?: number
  step?: number
  rows?: number
  colSpan?: 1 | 2 | 3 | 4
  relation?: RelationConfig
  /** render only when this predicate over the current form values holds */
  visibleIf?: (values: Record<string, unknown>) => boolean
}

function make(kind: FieldType, name: string, label: string, opts: FieldOptions = {}): FieldNode {
  return { type: 'field', kind, name, label, ...opts }
}

export function textInput(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('text', name, label, opts)
}

export function emailInput(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('email', name, label, opts)
}

export function passwordInput(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('password', name, label, opts)
}

export function numberInput(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('number', name, label, opts)
}

export function textarea(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('textarea', name, label, opts)
}

export function selectInput(
  name: string,
  label: string,
  options: FieldOption[],
  opts?: Omit<FieldOptions, 'options'>
): FieldNode {
  return make('select', name, label, { ...opts, options })
}

export function switchInput(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('switch', name, label, opts)
}

export function checkboxInput(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('checkbox', name, label, opts)
}

export function dateInput(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('date', name, label, opts)
}

/** multi-select rendered as a checkbox group; binds string[] */
export function checkboxGroup(
  name: string,
  label: string,
  options: FieldOption[],
  opts?: Omit<FieldOptions, 'options'>
): FieldNode {
  return make('checkboxgroup', name, label, { ...opts, options })
}

/** file upload control; uploads to /api/admin/uploads and binds the URL */
export function uploadInput(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('upload', name, label, opts)
}

/** relation picker: loads options from another resource's API */
export function relationInput(
  name: string,
  label: string,
  relation: { resource: string, labelKey: string },
  opts?: FieldOptions
): FieldNode {
  return make('relation', name, label, { ...opts, relation })
}

/** rich text editor (CKEditor 5 classic build) */
export function richtext(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('richtext', name, label, opts)
}

/** WYSIWYG rich text (contenteditable); stores HTML. */
export function wysiwyg(name: string, label: string, opts?: FieldOptions): FieldNode {
  return make('wysiwyg', name, label, opts)
}
