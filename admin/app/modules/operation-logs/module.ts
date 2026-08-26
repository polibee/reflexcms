import OperationLogResource from './admin/OperationLogResource'

export default defineModule({
  name: 'operation-logs',
  resources: [OperationLogResource],
  navGroups: [{ label: 'Platform', sort: 30 }]
})
