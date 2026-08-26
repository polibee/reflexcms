import ReplyResource from './admin/ReplyResource'

export default defineModule({
  name: 'replies',
  resources: [ReplyResource],
  navGroups: [{ label: 'Forum', sort: 30 }]
})
