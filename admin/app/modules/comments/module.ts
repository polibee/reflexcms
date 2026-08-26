import CommentResource from './admin/CommentResource'

export default defineModule({
  name: 'comments',
  resources: [CommentResource],
  navGroups: [{ label: 'Content', sort: 20 }]
})
