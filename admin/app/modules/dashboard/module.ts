import WidgetStats from './widgets/WidgetStats.vue'
import WidgetRecentArticles from './widgets/WidgetRecentArticles.vue'
import WidgetPublishTrend from './widgets/WidgetPublishTrend.vue'

export default defineModule({
  name: 'dashboard',
  navGroups: [{ label: 'General', sort: 0 }],
  widgets: [
    { name: 'stats-overview', span: 4, order: 1, component: WidgetStats },
    { name: 'recent-articles', label: 'Recent Articles', span: 2, order: 2, component: WidgetRecentArticles },
    { name: 'publish-trend', label: 'Publish Trend (10 weeks)', span: 2, order: 3, component: WidgetPublishTrend }
  ]
})
