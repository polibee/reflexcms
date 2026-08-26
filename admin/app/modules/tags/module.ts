import TagResource from './admin/TagResource'

export default defineModule({
  name: 'tags',
  resources: [TagResource],
  navGroups: [{ label: 'Content', sort: 20 }]
})
