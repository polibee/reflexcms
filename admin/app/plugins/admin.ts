import adminPanel from '~/admin/panels/admin.panel'
import dashboardModule from '~/modules/dashboard/module'
import articlesModule from '~/modules/articles/module'
import categoriesModule from '~/modules/categories/module'
import tagsModule from '~/modules/tags/module'
import commentsModule from '~/modules/comments/module'
import forumsModule from '~/modules/forums/module'
import topicsModule from '~/modules/topics/module'
import repliesModule from '~/modules/replies/module'
import usersModule from '~/modules/users/module'
import rolesModule from '~/modules/roles/module'
import notificationsModule from '~/modules/notifications/module'
import settingsModule from '~/modules/settings/module'
import operationLogsModule from '~/modules/operation-logs/module'
import layoutModule from '~/modules/layout/module'
import pagesModule from '~/modules/pages/module'

/**
 * Application composition root.
 * Register panels and business modules here - exactly like
 * Laravel registers service providers.
 *
 * Single-product extraction: remove the module imports + registration lines
 * of the domains you do not deploy (docs/模块化架构与产品组合.md §5).
 */
export default defineNuxtPlugin(() => {
  setPanel(adminPanel)

  const modules = [
    dashboardModule,
    articlesModule,
    categoriesModule,
    tagsModule,
    commentsModule,
    forumsModule,
    topicsModule,
    repliesModule,
    layoutModule,
    pagesModule,
    usersModule,
    rolesModule,
    notificationsModule,
    settingsModule,
    operationLogsModule
  ]
  for (const module of modules) {
    registerModule(module)
  }
})
