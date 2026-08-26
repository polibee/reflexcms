import ArticleResource from './admin/ArticleResource'

export default defineModule({
  name: 'articles',
  resources: [ArticleResource],
  navGroups: [{ label: 'Content', sort: 20 }]
})
