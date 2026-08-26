import { PERMISSION_GROUPS } from './permissionCatalog'

const ALL_OPTIONS = PERMISSION_GROUPS.flatMap(g =>
  g.permissions.map(p => ({ label: `${p.label}（${p.name}）`, value: p.name }))
)

const PRESETS: Record<string, { label: string, icon: string, perms: string[] }> = {
  content: {
    label: '内容组', icon: 'file-text',
    perms: ['articles.*', 'categories.*', 'tags.*', 'comments.view', 'users.view']
  },
  forum: {
    label: '论坛组', icon: 'messages-square',
    perms: [
      'forums.view',
      'topics.view', 'topics.edit', 'topics.delete', 'topics.pin', 'topics.feature',
      'replies.view', 'replies.delete',
      'comments.view', 'comments.moderate', 'users.view'
    ]
  },
  viewer: {
    label: '只读', icon: 'eye',
    perms: [
      'users.view', 'roles.view',
      'articles.view', 'categories.view', 'tags.view',
      'forums.view', 'topics.view', 'replies.view',
      'comments.view', 'notifications.view', 'settings.view', 'audit.view'
    ]
  }
}

export default defineResource({
  name: 'roles',
  model: 'Role',
  label: 'Role',
  labelPlural: 'Roles',
  icon: 'lock',
  group: 'Access Control',
  sort: 20,
  searchable: ['name', 'display_name'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('name', 'Name', { sortable: true }),
    textColumn('display_name', 'Display Name'),
    numberColumn('sort', 'Sort'),
    dateColumn('created_at', 'Created', { sortable: true })
  ],

  rowActions: Object.entries(PRESETS).map(([key, preset]) =>
    defineAction({
      name: `preset-${key}`,
      label: `套用${preset.label}预设`,
      icon: preset.icon,
      permission: 'roles.edit',
      visible: record => record.name !== 'super-admin',
      confirm: {
        title: `将权限重置为「${preset.label}」预设？现有权限将被覆盖。`,
        confirmLabel: '套用'
      },
      handler: async ({ record }) => {
        await $fetch(`/api/admin/roles/${record!.id}`, {
          method: 'PUT',
          body: { permissions: [...preset.perms] }
        })
        notify(`已套用「${preset.label}」预设`)
        emitAdminEvent('roles:refresh')
      }
    })
  ),

  form: () => [
    section('Role', [
      grid(2, [
        textInput('name', 'Name', { required: true, placeholder: 'editor' }),
        numberInput('sort', 'Sort Order', { min: 0, defaultValue: 100 })
      ]),
      textInput('display_name', 'Display Name', { required: true }),
      checkboxGroup('permissions', '权限（按组勾选）', ALL_OPTIONS, {
        helpText: '勾选即授予；后端校验未知条目'
      })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('name', 'Name'),
    textEntry('display_name', 'Display Name'),
    textEntry('sort', 'Sort Order'),
    textEntry('permissions', 'Permissions'),
    dateEntry('created_at', 'Created')
  ]
})
