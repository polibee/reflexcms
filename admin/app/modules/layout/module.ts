import WidgetResource from './admin/WidgetResource'
import MenuResource from './admin/MenuResource'
import CarouselResource from './admin/CarouselResource'

export default defineModule({
  name: 'layout',
  resources: [WidgetResource, MenuResource, CarouselResource],
  navGroups: [{ label: '布局管理', sort: 40 }]
})
