import RoleResource from './admin/RoleResource'

export default defineModule({
  name: 'roles',
  resources: [RoleResource],
  navGroups: [{ label: 'Access Control', sort: 10 }]
})
