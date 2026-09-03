/* =============================================================
 * Client-side resource metadata only. Field validation is authoritative on
 * the Goravel backend (ResourceSpec); the two declarations are kept in sync
 * by scripts/contract-test.mjs (plan §D3).
 * ============================================================= */

export interface ResourceConfig {
  label: string
  searchable: string[]
  permissionPrefix: string
}

const configs: Record<string, ResourceConfig> = {
  users: {
    label: 'User',
    searchable: ['username', 'email'],
    permissionPrefix: 'users'
  },
  roles: {
    label: 'Role',
    searchable: ['name', 'display_name'],
    permissionPrefix: 'roles'
  },
  articles: {
    label: 'Article',
    searchable: ['title', 'summary'],
    permissionPrefix: 'articles'
  },
  categories: {
    label: 'Category',
    searchable: ['name', 'slug'],
    permissionPrefix: 'categories'
  },
  tags: {
    label: 'Tag',
    searchable: ['name', 'slug'],
    permissionPrefix: 'tags'
  },
  forums: {
    label: 'Board',
    searchable: ['name', 'slug'],
    permissionPrefix: 'forums'
  },
  topics: {
    label: 'Topic',
    searchable: ['title'],
    permissionPrefix: 'topics'
  },
  replies: {
    label: 'Reply',
    searchable: ['content'],
    permissionPrefix: 'replies'
  },
  comments: {
    label: 'Comment',
    searchable: ['content'],
    permissionPrefix: 'comments'
  },
  notifications: {
    label: 'Notification',
    searchable: ['type'],
    permissionPrefix: 'notifications'
  },
  settings: {
    label: 'Setting',
    searchable: ['key'],
    permissionPrefix: 'settings'
  },
  operation_logs: {
    label: 'Operation Log',
    searchable: ['action', 'resource'],
    permissionPrefix: 'audit'
  },
  widgets: {
    label: 'Widget',
    searchable: ['title', 'widget_type'],
    permissionPrefix: 'widgets'
  },
  menus: {
    label: 'Menu',
    searchable: ['label', 'url'],
    permissionPrefix: 'menus'
  },
  carousels: {
    label: 'Carousel',
    searchable: ['title'],
    permissionPrefix: 'carousels'
  },
  pages: {
    label: 'Page',
    searchable: ['title', 'slug'],
    permissionPrefix: 'pages'
  },
  products: {
    label: 'Product',
    searchable: ['title'],
    permissionPrefix: 'products'
  },
  orders: {
    label: 'Order',
    searchable: ['order_no', 'title'],
    permissionPrefix: 'orders'
  },
  invites: {
    label: 'Invite',
    searchable: ['code'],
    permissionPrefix: 'invites'
  },
  board_moderators: {
    label: 'Board Moderator',
    searchable: [],
    permissionPrefix: 'board_moderators'
  }
}

export function getConfig(resource: string): ResourceConfig {
  const config = configs[resource]
  if (!config) {
    throw createError({ statusCode: 404, statusMessage: `Unknown resource "${resource}"` })
  }
  return config
}

export function hasConfig(resource: string): boolean {
  return resource in configs
}
