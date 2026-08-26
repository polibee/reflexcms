const statusBadges = {
  active: { label: 'Active', variant: 'success' },
  inactive: { label: 'Inactive', variant: 'secondary' },
  banned: { label: 'Banned', variant: 'warning' }
} as const

export default defineResource({
  name: 'users',
  model: 'User',
  label: 'User',
  labelPlural: 'Users',
  icon: 'users',
  group: 'Access Control',
  sort: 10,
  searchable: ['username', 'email'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('username', 'Username', { sortable: true }),
    textColumn('email', 'Email'),
    numberColumn('role_id', 'Role ID'),
    badgeColumn('status', 'Status', statusBadges),
    dateColumn('created_at', 'Created', { sortable: true })
  ],

  form: () => [
    section('Account', [
      grid(2, [
        textInput('username', 'Username', { required: true }),
        emailInput('email', 'Email', { required: true })
      ]),
      grid(2, [
        passwordInput('password', 'Password', { placeholder: 'Leave empty to keep current' }),
        numberInput('role_id', 'Role ID', { min: 1, placeholder: 'See Access Control → Roles' })
      ]),
      grid(2, [
        selectInput('status', 'Status', [
          { label: 'Active', value: 'active' },
          { label: 'Inactive', value: 'inactive' },
          { label: 'Banned', value: 'banned' }
        ], { defaultValue: 'active' }),
        textInput('website', 'Website')
      ]),
      textarea('bio', 'Bio')
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('username', 'Username'),
    textEntry('email', 'Email'),
    textEntry('role_id', 'Role ID'),
    badgeEntry('status', 'Status', statusBadges),
    dateEntry('created_at', 'Created')
  ]
})
