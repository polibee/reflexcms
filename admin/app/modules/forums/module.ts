import ForumResource from './admin/ForumResource'

export default defineModule({
  name: 'forums',
  resources: [ForumResource],
  navGroups: [{ label: 'Forum', sort: 30 }]
})
