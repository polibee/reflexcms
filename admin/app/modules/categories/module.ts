import CategoryResource from './admin/CategoryResource'

export default defineModule({
  name: 'categories',
  resources: [CategoryResource],
  navGroups: [{ label: 'Content', sort: 20 }]
})
