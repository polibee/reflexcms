import UserResource from './admin/UserResource'

export default defineModule({
  name: 'users',
  resources: [UserResource],
  navGroups: [{ label: 'Access Control', sort: 10 }]
})
