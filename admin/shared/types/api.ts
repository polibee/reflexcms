// Shared API contracts between app (client) and server (Nitro)

export interface ListQuery {
  q?: string
  page?: number
  perPage?: number
  sortBy?: string
  sortDir?: 'asc' | 'desc'
  filters?: Record<string, string>
}

export interface Paginated<T> {
  items: T[]
  total: number
  page: number
  perPage: number
  totalPages: number
}

// Roles are data-driven in the backend roles table (plan §O4); any slug may
// appear here. Actual capability always comes from `permissions`.
export type AdminRole = string

export interface AuthUser {
  id: number
  name: string
  email: string
  role: AdminRole
  permissions: string[]
}
