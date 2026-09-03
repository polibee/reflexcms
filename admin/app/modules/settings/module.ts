import SettingResource from './admin/SettingResource'

export default defineModule({
  name: 'settings',
  resources: [SettingResource],
  navGroups: [{
    label: 'Platform',
    sort: 10,
    items: [{ label: '邮件群发', to: '/admin/mailer', icon: 'mail' }]
  }]
})
