export default defineResource({
  name: 'operation_logs',
  model: 'AdminOperationLog',
  label: 'Operation Log',
  labelPlural: 'Operation Logs',
  icon: 'scroll-text',
  group: 'Platform',
  sort: 30,
  searchable: ['action', 'resource'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    numberColumn('user_id', 'User'),
    textColumn('action', 'Action'),
    textColumn('resource', 'Resource'),
    textColumn('resource_id', 'Target'),
    numberColumn('status_code', 'Status'),
    dateColumn('created_at', 'At', { sortable: true })
  ],

  form: () => [],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('user_id', 'User ID'),
    textEntry('action', 'Action'),
    textEntry('resource', 'Resource'),
    textEntry('resource_id', 'Target ID'),
    textEntry('status_code', 'Status Code'),
    textEntry('ip', 'IP'),
    textEntry('user_agent', 'User Agent'),
    textEntry('payload', 'Payload Snapshot'),
    datetimeEntry('created_at', 'At')
  ]
})
