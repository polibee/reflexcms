import SettingResource from './admin/SettingResource'

export default defineModule({
  name: 'settings',
  resources: [SettingResource],
  navGroups: [{ label: 'Platform', sort: 10 }]
})
