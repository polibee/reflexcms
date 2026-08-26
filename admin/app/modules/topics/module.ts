import TopicResource from './admin/TopicResource'

export default defineModule({
  name: 'topics',
  resources: [TopicResource],
  navGroups: [{ label: 'Forum', sort: 30 }]
})
