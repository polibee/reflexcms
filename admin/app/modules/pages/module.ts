import PageResource from './admin/PageResource'

export default defineModule({
  name: 'pages',
  resources: [PageResource],
  navGroups: [{ label: '页面管理', sort: 45 }]
})
