export default defineResource({
  name: 'board_moderators',
  model: 'BoardModerator',
  label: '版主管理',
  labelPlural: '版主管理',
  icon: 'shield-check',
  group: 'Forum',
  sort: 20,
  searchable: [],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('board_id', '板块 ID', { sortable: true }),
    textColumn('user_id', '用户 ID', { sortable: true }),
    dateColumn('created_at', 'Created')
  ],

  form: () => [
    section('版主任命', [
      relationInput('board_id', '板块', { resource: 'forums', labelKey: 'name' }, { required: true }),
      relationInput('user_id', '用户', { resource: 'users', labelKey: 'username' }, { required: true }),
      selectInput('action_hint', '操作说明', [
        { label: '被任命的用户可在对应板块执行 置顶/精华/关闭/设最佳回复', value: 'info' }
      ], { helpText: '版主仅能在自己负责的板块内执行管理操作' })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('board_id', '板块'),
    textEntry('user_id', '用户'),
    textEntry('created_at', 'Created')
  ]
})
