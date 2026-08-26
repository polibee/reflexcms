// Static mirror of the backend permission catalog
// (backend modules/access/services/catalog.go). The backend remains the
// authoritative validator — unknown entries are rejected with 422 — and the
// contract suite keeps both sides in sync.

export interface PermissionItem {
  name: string
  label: string
}

export interface PermissionGroup {
  group: string
  permissions: PermissionItem[]
}

export const PERMISSION_GROUPS: PermissionGroup[] = [
  {
    group: '内容',
    permissions: [
      { name: 'articles.view', label: '文章查看' },
      { name: 'articles.create', label: '文章创建' },
      { name: 'articles.edit', label: '文章编辑' },
      { name: 'articles.delete', label: '文章删除' },
      { name: 'articles.publish', label: '文章发布/定时' },
      { name: 'categories.view', label: '分类查看' },
      { name: 'categories.create', label: '分类创建' },
      { name: 'categories.edit', label: '分类编辑' },
      { name: 'categories.delete', label: '分类删除' },
      { name: 'tags.view', label: '标签查看' },
      { name: 'tags.create', label: '标签创建' },
      { name: 'tags.edit', label: '标签编辑' },
      { name: 'tags.delete', label: '标签删除' },
      { name: 'comments.view', label: '评论查看' },
      { name: 'comments.moderate', label: '评论审核' }
    ]
  },
  {
    group: '论坛',
    permissions: [
      { name: 'forums.view', label: '板块查看' },
      { name: 'forums.create', label: '板块创建' },
      { name: 'forums.edit', label: '板块编辑' },
      { name: 'forums.delete', label: '板块删除' },
      { name: 'topics.view', label: '帖子查看' },
      { name: 'topics.edit', label: '帖子编辑' },
      { name: 'topics.delete', label: '帖子删除' },
      { name: 'topics.pin', label: '帖子置顶' },
      { name: 'topics.feature', label: '帖子加精' },
      { name: 'replies.view', label: '回复查看' },
      { name: 'replies.create', label: '回复创建' },
      { name: 'replies.delete', label: '回复删除' }
    ]
  },
  {
    group: '平台',
    permissions: [
      { name: 'users.view', label: '用户查看' },
      { name: 'users.create', label: '用户创建' },
      { name: 'users.edit', label: '用户编辑' },
      { name: 'users.delete', label: '用户删除' },
      { name: 'roles.view', label: '角色查看' },
      { name: 'roles.create', label: '角色创建' },
      { name: 'roles.edit', label: '角色编辑' },
      { name: 'roles.delete', label: '角色删除' },
      { name: 'notifications.view', label: '通知查看' },
      { name: 'notifications.edit', label: '通知操作' },
      { name: 'settings.view', label: '设置查看' },
      { name: 'settings.edit', label: '设置修改' },
      { name: 'audit.view', label: '审计日志查看' }
    ]
  }
]

export const ALL_PERMISSIONS: PermissionItem[] = PERMISSION_GROUPS.flatMap(g => g.permissions)

/** Form field names cannot contain dots (path separators), so they are
 *  encoded as perm__articles__view — matching the backend Transform output. */
export const encodePermField = (perm: string) => `perm__${perm.replaceAll('.', '__')}`

export const decodePermField = (field: string) =>
  field.startsWith('perm__') ? field.slice('perm__'.length).replaceAll('__', '.') : field
