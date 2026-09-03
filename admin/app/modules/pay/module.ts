import ProductResource from './admin/ProductResource'
import OrderResource from './admin/OrderResource'

export default defineModule({
  name: 'pay',
  resources: [ProductResource, OrderResource],
  navGroups: [{ label: '商城', sort: 42 }]
})
