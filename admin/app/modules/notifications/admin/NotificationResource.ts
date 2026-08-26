export default defineResource({
  name: 'notifications',
  model: 'Notification',
  label: 'Notification',
  labelPlural: 'Notifications',
  icon: 'bell',
  group: 'Platform',
  sort: 10,
  searchable: ['type'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    numberColumn('user_id', 'User'),
    textColumn('type', 'Type'),
    textColumn('data', 'Data'),
    booleanColumn('read_at', 'Read'),
    dateColumn('created_at', 'Created', { sortable: true })
  ],

  form: () => [],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('user_id', 'User ID'),
    textEntry('type', 'Type'),
    textEntry('data', 'Data'),
    dateEntry('created_at', 'Created')
  ],

  rowActions: [
    defineAction({
      name: 'mark-read',
      label: 'Mark All Read',
      icon: 'check',
      permission: 'notifications.edit',
      handler: async () => {
        await $fetch('/api/admin/notifications/read-all', { method: 'POST' })
        notify('All notifications marked as read')
        emitAdminEvent('notifications:refresh')
      }
    })
  ]
})
