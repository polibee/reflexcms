import { describe, expect, it } from 'vitest'
import { can } from '../../app/admin/permissions'

const admin = { permissions: ['*'] }
const editor = { permissions: ['posts.view', 'posts.create', 'posts.edit', 'orders.view'] }
const nobody = { permissions: [] }

describe('can() permission matcher', () => {
  it('allows actions that declare no required permission', () => {
    expect(can(undefined, null)).toBe(true)
    expect(can(undefined, nobody)).toBe(true)
  })

  it('denies everything when there is no user', () => {
    expect(can('users.view')).toBe(false)
    expect(can('users.view', null)).toBe(false)
  })

  it('matches exact permissions', () => {
    expect(can('posts.view', editor)).toBe(true)
    expect(can('users.delete', editor)).toBe(false)
  })

  it('grants the global wildcard', () => {
    expect(can('users.delete', admin)).toBe(true)
    expect(can('anything.at.all', admin)).toBe(true)
  })

  it('supports resource-level wildcards without leaking to other resources', () => {
    const postsAdmin = { permissions: ['posts.*'] }
    expect(can('posts.view', postsAdmin)).toBe(true)
    expect(can('posts.delete', postsAdmin)).toBe(true)
    expect(can('users.view', postsAdmin)).toBe(false)
  })

  it('denies when no grant matches', () => {
    expect(can('orders.create', editor)).toBe(false)
    expect(can('posts.view', nobody)).toBe(false)
  })
})
