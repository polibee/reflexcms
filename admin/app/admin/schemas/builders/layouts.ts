import type { GridNode, SchemaNode, SectionNode } from '../../core/types'

type SectionOptions = Partial<Pick<SectionNode, 'title' | 'description'>>

/** section('Title', [...]) or section({ title, description }, [...]) */
export function section(options: SectionOptions | string, children: SchemaNode[]): SectionNode {
  const opts = typeof options === 'string' ? { title: options } : options
  return { type: 'section', title: opts.title, description: opts.description, children }
}

export function grid(columns: 1 | 2 | 3 | 4, children: SchemaNode[]): GridNode {
  return { type: 'grid', columns, children }
}
