import WidgetResource from './admin/WidgetResource'
import MenuResource from './admin/MenuResource'
import CarouselResource from './admin/CarouselResource'
import InviteResource from './admin/InviteResource'

export default defineModule({
  name: 'layout',
  resources: [WidgetResource, MenuResource, CarouselResource, InviteResource],
  navGroups: [{ label: '布局管理', sort: 40 }]
})
