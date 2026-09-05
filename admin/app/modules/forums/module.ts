import ForumResource from './admin/ForumResource'
import BoardModeratorResource from './admin/BoardModeratorResource'

export default defineModule({
  name: 'forums',
  resources: [ForumResource, BoardModeratorResource],
  navGroups: [{
    label: 'Forum',
    sort: 30,
    items: [{ label: '等级与签到设置', to: '/admin/forum-levels', icon: 'trending-up' }]
  }]
})
