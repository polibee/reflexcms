import { z, type ZodTypeAny } from 'zod'
import type { FieldNode, SchemaNode } from '../core/types'

/* =============================================================
 * Schema -> Zod compiler. Field definitions stay declarative;
 * validation rules are derived automatically from the schema.
 * ============================================================= */

/** kinds whose base rule is a plain z.string(), safe for .min() */
const STRING_BASE_KINDS: FieldNode['kind'][] = ['text', 'password', 'textarea', 'date', 'upload', 'richtext', 'wysiwyg', 'tiptap']

function baseRuleFor(node: FieldNode): ZodTypeAny {
  switch (node.kind) {
    case 'email':
      return z.string().email(`${node.label} must be a valid email`)
    case 'number':
      return node.required
        ? z.coerce.number({ message: `${node.label} is required` })
        : z.preprocess(
            v => (v === '' || v === null || v === undefined ? undefined : Number(v)),
            z.number().optional()
          )
    case 'switch':
    case 'checkbox':
      return z.boolean().default(false)
    case 'checkboxgroup':
      return z.array(z.string()).default([])
    case 'select':
      return z.union([z.string(), z.number(), z.null()]).transform(v => v ?? '')
    case 'relation':
      return z.union([z.string(), z.number(), z.null()]).transform(v => (v === null ? '' : v))
    default:
      return z.string()
  }
}

function compileField(node: FieldNode): ZodTypeAny {
  const extras = node.rules ?? []
  let rule = baseRuleFor(node)

  if (node.required && node.kind === 'email') {
    rule = z.string().min(1, `${node.label} is required`).email(`${node.label} must be a valid email`)
  } else if (node.required && STRING_BASE_KINDS.includes(node.kind)) {
    rule = (rule as z.ZodString).min(1, `${node.label} is required`)
  } else if (node.required && node.kind === 'number') {
    rule = z.preprocess(
      v => (v === '' || v === null || v === undefined ? undefined : Number(v)),
      z.number({ message: `${node.label} is required` })
    )
  } else if (node.required && (node.kind === 'select' || node.kind === 'relation')) {
    rule = z.union([z.string(), z.number()])
      .refine(v => v !== '' && v !== null && v !== undefined, `${node.label} is required`)
  } else if (!node.required) {
    // Non-required fields may be absent entirely — conditionally-rendered
    // fields (visibleIf) never mount, leaving their keys undefined.
    rule = rule.optional()
  }

  for (const extra of extras) {
    rule = rule.pipe(extra)
  }
  return rule
}

export function schemaToZod(schema: SchemaNode[]): z.ZodObject<Record<string, ZodTypeAny>> {
  const shape: Record<string, ZodTypeAny> = {}
  // Layout nodes nest arbitrarily (section → grid → field), so walk the
  // whole tree — a one-level walk silently drops nested fields from both
  // validation AND parsed output.
  collectFields(schema, shape)
  return z.object(shape)
}

function collectFields(nodes: SchemaNode[], shape: Record<string, ZodTypeAny>): void {
  for (const node of nodes) {
    if (node.type === 'field') {
      shape[node.name] = compileField(node)
    } else if ('children' in node && Array.isArray(node.children)) {
      collectFields(node.children, shape)
    }
  }
}
