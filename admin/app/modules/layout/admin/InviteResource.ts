export default defineResource({
  name: 'invites',
  model: 'InviteCode',
  label: '邀请码',
  labelPlural: '邀请码',
  icon: 'ticket',
  group: '布局管理',
  sort: 20,
  searchable: ['code'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('code', '邀请码'),
    numberColumn('max_uses', '可用次数', { sortable: true }),
    numberColumn('used_count', '已用次数', { sortable: true }),
    dateColumn('expires_at', '过期时间'),
    dateColumn('created_at', 'Created')
  ],

  form: () => [
    section('邀请码', [
      textInput('code', '邀请码', {
        placeholder: '留空自动生成（INV-XXXXXXXXXX）',
        colSpan: 2
      }),
      grid(2, [
        numberInput('max_uses', '可用次数', { required: true, min: 1, defaultValue: 1 }),
        dateInput('expires_at', '过期时间（可选）')
      ])
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('code', '邀请码'),
    textEntry('max_uses', '可用次数'),
    textEntry('used_count', '已用次数'),
    textEntry('expires_at', '过期时间'),
    textEntry('created_at', 'Created')
  ]
})
