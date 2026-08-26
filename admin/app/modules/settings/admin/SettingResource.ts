export default defineResource({
  name: 'settings',
  model: 'Setting',
  label: '站点设置',
  labelPlural: '站点设置',
  icon: 'settings',
  group: 'Platform',
  sort: 20,
  searchable: [],

  table: () => [
    textColumn('key', '键'),
    textColumn('value', '值'),
    textColumn('group', '分组'),
    dateColumn('updated_at', '更新时间', { sortable: true })
  ],

  form: () => [
    section('Setting', [
      textInput('key', 'Key', { required: true }),
      textarea('value', 'Value', { rows: 3 }),
      selectInput('group', 'Group', [
        { label: 'Site', value: 'site' },
        { label: 'Registration', value: 'registration' },
        { label: 'Email', value: 'email' },
        { label: 'SEO', value: 'seo' },
        { label: 'General', value: 'general' }
      ])
    ])
  ],

  infolist: () => [
    textEntry('key', 'Key'),
    textEntry('value', 'Value'),
    textEntry('group', 'Group')
  ]
})
