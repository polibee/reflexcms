import NotificationResource from './admin/NotificationResource'

export default defineModule({
  name: 'notifications',
  resources: [NotificationResource],
  navGroups: [{ label: 'Platform', sort: 10 }]
})
