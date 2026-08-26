import { z } from 'zod'
import { describe, expect, it } from 'vitest'
import { schemaToZod } from '../../app/admin/forms/schemaToZod'

describe('schemaToZod compiler', () => {
  it('collects fields from nested sections and grids into one flat shape', () => {
    const schema = schemaToZod([
      {
        type: 'section',
        title: 'Account',
        children: [
          { type: 'grid', columns: 2, children: [
            { type: 'field', kind: 'text', name: 'name', label: 'Name', required: true },
            { type: 'field', kind: 'email', name: 'email', label: 'Email', required: true }
          ] }
        ]
      },
      { type: 'field', kind: 'switch', name: 'active', label: 'Active' }
    ])

    const result = schema.safeParse({ name: 'Ada', email: 'ada@example.com' })
    expect(result.success).toBe(true)
    if (result.success) {
      expect(result.data.active).toBe(false) // switch defaults to false
    }
  })

  it('rejects missing required text with the field label in the message', () => {
    const schema = schemaToZod([
      { type: 'field', kind: 'text', name: 'title', label: 'Title', required: true }
    ])
    const result = schema.safeParse({ title: '' })
    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues[0]?.message).toContain('Title is required')
    }
  })

  it('validates email format and required-ness independently', () => {
    const schema = schemaToZod([
      { type: 'field', kind: 'email', name: 'email', label: 'Email', required: true }
    ])
    expect(schema.safeParse({ email: '' }).success).toBe(false)
    expect(schema.safeParse({ email: 'not-an-email' }).success).toBe(false)
    expect(schema.safeParse({ email: 'a@b.dev' }).success).toBe(true)
  })

  it('coerces number input and keeps it optional when not required', () => {
    const schema = schemaToZod([
      { type: 'field', kind: 'number', name: 'price', label: 'Price' }
    ])
    expect(schema.safeParse({ price: '' }).success).toBe(true)
    expect(schema.safeParse({ price: null }).success).toBe(true)
    const ok = schema.safeParse({ price: '42' })
    expect(ok.success).toBe(true)
    if (ok.success) expect(ok.data.price).toBe(42)
    expect(schema.safeParse({ price: 'abc' }).success).toBe(false)
  })

  it('requires numbers when flagged', () => {
    const schema = schemaToZod([
      { type: 'field', kind: 'number', name: 'amount', label: 'Amount', required: true }
    ])
    expect(schema.safeParse({}).success).toBe(false)
    expect(schema.safeParse({ amount: '' }).success).toBe(false)
    const ok = schema.safeParse({ amount: '7' })
    expect(ok.success).toBe(true)
  })

  it('rejects empty required select/relation but accepts valid values', () => {
    const schema = schemaToZod([
      { type: 'field', kind: 'select', name: 'role', label: 'Role', required: true },
      { type: 'field', kind: 'relation', name: 'authorId', label: 'Author', required: true }
    ])
    expect(schema.safeParse({ role: '', authorId: null }).success).toBe(false)
    expect(schema.safeParse({ role: 'admin', authorId: 3 }).success).toBe(true)
  })

  it('appends extra zod rules through the pipe', () => {
    const schema = schemaToZod([
      {
        type: 'field',
        kind: 'text',
        name: 'code',
        label: 'Code',
        rules: [z.string().max(4, 'Too long')]
      }
    ])
    const fail = schema.safeParse({ code: 'ABCDEFGH' })
    expect(fail.success).toBe(false)
    if (!fail.success) {
      expect(fail.error.issues[0]?.message).toContain('Too long')
    }
    expect(schema.safeParse({ code: 'AB' }).success).toBe(true)
  })

  it('supports boolean switches explicitly set to true', () => {
    const schema = schemaToZod([
      { type: 'field', kind: 'checkbox', name: 'agree', label: 'Agree' }
    ])
    const result = schema.safeParse({ agree: true })
    expect(result.success).toBe(true)
    if (result.success) expect(result.data.agree).toBe(true)
  })
})
