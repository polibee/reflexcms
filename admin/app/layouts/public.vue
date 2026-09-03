<script setup lang="ts">
const { config, siteName, navItems } = useSiteConfig()

interface SidebarWidget {
  type: string
  title: string
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: any
}
const sidebarWidgets = ref<SidebarWidget[]>([])
const route = useRoute()

const { get } = usePublicApi()

// The layout component survives client-side navigation, so the widget area
// must be reactive: refetch whenever the route crosses the blog/forum split,
// otherwise the previous section's cards linger (the "wrong sidebar" bug).
// User-center pages (profile/notifications/DMs/settings) get NO content
// widgets — categories/tags clouds don't belong there.
const sidebarArea = computed(() => {
  const p = route.path
  if (p.startsWith('/users') || p.startsWith('/notifications') || p.startsWith('/messages') || p.startsWith('/settings') || p.startsWith('/shop')) {
    return ''
  }
  return p.startsWith('/forums') || p.startsWith('/topics') ? 'forum_sidebar' : 'blog_sidebar'
})

const showAside = computed(() => sidebarWidgets.value.length > 0)

async function loadSidebar(area: string) {
  if (!area) {
    sidebarWidgets.value = []
    return
  }
  try {
    const res = await get<{ widgets: SidebarWidget[] }>(`sidebar?area=${area}`)
    sidebarWidgets.value = res.widgets ?? []
  } catch { /* silent */ }
}

watch(sidebarArea, (area) => { void loadSidebar(area) }, { immediate: true })

// Site config is shared state (useSiteConfig); the layout guarantees it is
// loaded so header/footer and page sections always have it. The signed-in
// identity is shared too and re-validated when the tab regains focus, so
// logging in from another tab reflects here without a manual reload.
const { user: authUser, refresh: refreshAuth } = useAuthUser()

/* user-center card counters (topics / replies / comments / favorites /
 * mentions / unread DMs) — fetched whenever the identity changes */
interface MyStats {
  topics: number
  replies: number
  comments: number
  favorites: number
  mentions: number
  messages_unread: number
}
const stats = ref<MyStats | null>(null)

interface MyPoints {
  currency: string
  balance: number
  total_earned: number
  level: { current: number, next: number, percent: number, next_floor: number }
  daily: {
    topic: { earned: number, cap: number, reward: number }
    reply: { earned: number, cap: number, reward: number }
    signin: { done: boolean, reward: number }
  }
}
const points = ref<MyPoints | null>(null)
const signinBusy = ref(false)

async function loadGamification() {
  if (!authUser.value) {
    points.value = null
    return
  }
  try {
    points.value = await get<MyPoints>('me/points')
  } catch { points.value = null }
}

async function doSignin() {
  if (signinBusy.value) return
  signinBusy.value = true
  try {
    await $fetch('/api/v1/me/signin', { method: 'POST' }).catch(() => null)
    await loadGamification()
  } finally {
    signinBusy.value = false
  }
}

watch(authUser, async (u) => {
  if (!u) {
    stats.value = null
    points.value = null
    return
  }
  try {
    stats.value = await get<MyStats>('me/stats')
  } catch { stats.value = null }
  void loadGamification()
}, { immediate: true })

onMounted(async () => {
  if (!config.value) {
    try {
      config.value = await get<SiteConfig>('site/config')
    } catch { /* silent */ }
  }
  void refreshAuth()
  window.addEventListener('focus', on_focus_refresh)
})

function on_focus_refresh() {
  void refreshAuth()
}

onBeforeUnmount(() => {
  window.removeEventListener('focus', on_focus_refresh)
})

useHead({
  titleTemplate: title => title ? `${title} · ${siteName.value}` : siteName.value
})

/* Header/footer navigation comes from the admin 菜单管理 (menus table,
 * locations top_nav / footer_nav). Falls back to the site-config sections
 * when an admin hasn't defined top_nav entries. */
interface MenuItem {
  id: number
  label: string
  url: string
}

const topNavMenus = ref<MenuItem[]>([])
const footerNavMenus = ref<MenuItem[]>([])

async function loadMenus(location: string): Promise<MenuItem[]> {
  try {
    const res = await get<{ items: MenuItem[] }>(`menus?location=${location}`)
    return (res.items ?? []).filter(m => m.label && m.url)
  } catch {
    return []
  }
}

onMounted(async () => {
  const [top, footer] = await Promise.all([loadMenus('top_nav'), loadMenus('footer_nav')])
  topNavMenus.value = top
  footerNavMenus.value = footer
})

const headerNav = computed(() => {
  const items = topNavMenus.value.length
    ? topNavMenus.value.map(m => ({ label: m.label, to: m.url }))
    : navItems.value
  // site.mode gates content sections; nav links to disabled sections hide too
  return items.filter(i => navAllowed(i.to))
})

const footerNav = computed(() => footerNavMenus.value.filter(m => navAllowed(m.url)))

/* map a nav URL to the site.config section it requires */
function navAllowed(url: string): boolean {
  const sections = config.value?.sections
  if (!sections) return true
  if (url.startsWith('/articles')) return sections.includes('articles')
  if (url.startsWith('/forums') || url.startsWith('/topics')) return sections.includes('topics')
  return true
}
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <header class="sticky top-0 z-40 border-b bg-white/80 backdrop-blur">
      <div class="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
        <NuxtLink
          to="/"
          class="text-lg font-bold text-gray-900"
        >{{ siteName }}</NuxtLink>
        <nav class="flex items-center gap-6">
          <NuxtLink
            v-for="item in headerNav"
            :key="item.to"
            :to="item.to"
            class="text-sm text-gray-600 hover:text-gray-900"
          >{{ item.label }}</NuxtLink>
        </nav>
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-4 py-8">
      <div class="flex gap-8">
        <div class="min-w-0 flex-1">
          <slot />
        </div>
        <aside
          v-if="showAside"
          class="hidden w-72 shrink-0 space-y-4 lg:block"
        >
          <div
            v-for="(widget, i) in sidebarWidgets"
            :key="i"
            class="rounded-lg border bg-white p-4"
          >
            <h3 class="mb-3 text-sm font-semibold text-gray-900">
              {{ widget.title }}
            </h3>

            <div
              v-if="widget.type === 'categories'"
              class="space-y-1"
            >
              <NuxtLink
                v-for="cat in widget.data"
                :key="cat.id"
                :to="`/categories/${cat.slug}`"
                class="block rounded px-2 py-1 text-sm text-gray-600 hover:bg-gray-50"
              >{{ cat.name }}</NuxtLink>
            </div>

            <div
              v-else-if="widget.type === 'tags'"
              class="flex flex-wrap gap-1.5"
            >
              <span
                v-for="tag in widget.data"
                :key="tag.id"
                class="rounded-full bg-gray-100 px-2.5 py-0.5 text-xs text-gray-600"
              >{{ tag.name }}</span>
            </div>

            <div
              v-else-if="widget.type?.includes('articles')"
              class="space-y-2"
            >
              <NuxtLink
                v-for="a in widget.data"
                :key="a.id"
                :to="`/articles/${a.slug}`"
                class="block text-sm text-gray-600 hover:text-gray-900"
              >{{ a.title }}</NuxtLink>
            </div>

            <div
              v-else-if="widget.type?.includes('topics')"
              class="space-y-2"
            >
              <NuxtLink
                v-for="t in widget.data"
                :key="t.id"
                :to="`/topics/${t.id}`"
                class="block text-sm text-gray-600 hover:text-gray-900"
              >{{ t.title }}</NuxtLink>
            </div>

            <div
              v-else-if="widget.type === 'boards'"
              class="space-y-1"
            >
              <NuxtLink
                v-for="b in widget.data"
                :key="b.id"
                :to="`/forums/${b.id}`"
                class="block rounded px-2 py-1 text-sm text-gray-600 hover:bg-gray-50"
              >{{ b.name }} ({{ b.topic_count }})</NuxtLink>
            </div>

            <div
              v-else-if="widget.type === 'user_stats'"
              class="space-y-1 text-sm text-gray-600"
            >
              <div>用户: {{ widget.data?.users ?? 0 }}</div>
              <div>帖子: {{ widget.data?.topics ?? 0 }}</div>
            </div>

            <div
              v-else-if="widget.type === 'author_center'"
              class="space-y-3 text-sm"
            >
              <template v-if="authUser">
                <NuxtLink
                  :to="`/users/${authUser.id}`"
                  class="flex items-center gap-2 rounded-lg bg-gray-50 p-2 hover:bg-gray-100"
                >
                  <span class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-500 text-xs font-bold text-white">
                    {{ authUser.name.slice(0, 1).toUpperCase() }}
                  </span>
                  <span class="min-w-0">
                    <span class="block truncate font-medium text-gray-900">{{ authUser.name }}</span>
                    <span class="block text-xs text-gray-400">个人主页</span>
                  </span>
                </NuxtLink>

                <!-- counters: topics / replies / comments / favorites / mentions / DMs -->
                <div class="grid grid-cols-3 gap-1 text-center text-xs">
                  <NuxtLink
                    :to="`/users/${authUser.id}?tab=topics`"
                    class="rounded-md py-1.5 hover:bg-gray-50"
                  >
                    <span class="block font-semibold text-gray-900">{{ stats?.topics ?? 0 }}</span>
                    <span class="text-gray-400">主题帖</span>
                  </NuxtLink>
                  <NuxtLink
                    :to="`/users/${authUser.id}?tab=replies`"
                    class="rounded-md py-1.5 hover:bg-gray-50"
                  >
                    <span class="block font-semibold text-gray-900">{{ stats?.replies ?? 0 }}</span>
                    <span class="text-gray-400">回复</span>
                  </NuxtLink>
                  <NuxtLink
                    :to="`/users/${authUser.id}`"
                    class="rounded-md py-1.5 hover:bg-gray-50"
                  >
                    <span class="block font-semibold text-gray-900">{{ stats?.comments ?? 0 }}</span>
                    <span class="text-gray-400">评论</span>
                  </NuxtLink>
                  <NuxtLink
                    :to="`/users/${authUser.id}?tab=favorites`"
                    class="rounded-md py-1.5 hover:bg-gray-50"
                  >
                    <span class="block font-semibold text-gray-900">{{ stats?.favorites ?? 0 }}</span>
                    <span class="text-gray-400">收藏</span>
                  </NuxtLink>
                  <NuxtLink
                    to="/notifications"
                    class="rounded-md py-1.5 hover:bg-gray-50"
                  >
                    <span
                      class="block font-semibold"
                      :class="stats?.mentions ? 'text-blue-600' : 'text-gray-900'"
                    >{{ stats?.mentions ?? 0 }}</span>
                    <span class="text-gray-400">提及</span>
                  </NuxtLink>
                  <NuxtLink
                    to="/messages"
                    class="relative rounded-md py-1.5 hover:bg-gray-50"
                  >
                    <span
                      class="block font-semibold"
                      :class="stats?.messages_unread ? 'text-blue-600' : 'text-gray-900'"
                    >{{ stats?.messages_unread ?? 0 }}</span>
                    <span class="text-gray-400">私信</span>
                  </NuxtLink>
                </div>

                <!-- level progress + daily tasks -->
                <div
                  v-if="points"
                  class="space-y-2 rounded-lg bg-gray-50 p-2.5"
                >
                  <div>
                    <div class="flex items-center justify-between text-xs">
                      <span class="font-semibold text-gray-900">Lv{{ points.level.current }}</span>
                      <span
                        v-if="points.level.next_floor"
                        class="text-gray-400"
                      >{{ points.total_earned }} / {{ points.level.next_floor }} → Lv{{ points.level.next }}</span>
                      <span
                        v-else
                        class="text-amber-500"
                      >满级 🎉</span>
                    </div>
                    <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-200">
                      <div
                        class="h-full rounded-full bg-blue-500 transition-all"
                        :style="{ width: points.level.percent + '%' }"
                      />
                    </div>
                    <p class="mt-1 text-xs text-gray-400">
                      {{ points.currency }} × {{ points.balance }}
                    </p>
                  </div>
                  <div class="space-y-1 text-xs text-gray-500">
                    <p>发帖 {{ points.daily.topic.earned }}/{{ points.daily.topic.cap }} · 评论 {{ points.daily.reply.earned }}/{{ points.daily.reply.cap }}</p>
                    <button
                      v-if="!points.daily.signin.done"
                      type="button"
                      :disabled="signinBusy"
                      class="w-full rounded-md bg-amber-500 py-1 font-medium text-white hover:bg-amber-600 disabled:opacity-50"
                      @click="doSignin"
                    >
                      签到 +{{ points.daily.signin.reward }} {{ points.currency }}
                    </button>
                    <p
                      v-else
                      class="rounded-md bg-green-50 py-1 text-center text-green-600"
                    >
                      ✓ 今日已签到
                    </p>
                  </div>
                </div>

                <NuxtLink
                  to="/forums/new"
                  class="block rounded-md bg-blue-600 py-2 text-center font-medium text-white hover:bg-blue-700"
                >✏️ 发布帖子</NuxtLink>

                <div class="flex gap-2">
                  <NuxtLink
                    to="/settings"
                    class="flex-1 rounded-md border border-gray-300 py-1.5 text-center text-gray-600 hover:bg-gray-50"
                  >个人设置</NuxtLink>
                  <NuxtLink
                    to="/admin"
                    class="flex-1 rounded-md border border-gray-300 py-1.5 text-center text-gray-600 hover:bg-gray-50"
                  >进入后台</NuxtLink>
                </div>
              </template>
              <template v-else>
                <NuxtLink
                  to="/login"
                  class="block rounded-md bg-blue-600 py-1.5 text-center font-medium text-white hover:bg-blue-700"
                >登录</NuxtLink>
                <p class="pt-1 text-xs leading-relaxed text-gray-400">
                  登录后可发帖、评论、收藏与管理内容
                </p>
              </template>
            </div>

            <div
              v-else-if="widget.type === 'composer'"
              class="space-y-2 text-sm"
            >
              <NuxtLink
                to="/forums/new"
                class="block rounded-md bg-blue-600 py-2 text-center font-medium text-white hover:bg-blue-700"
              >✏️ 发布帖子</NuxtLink>
              <p class="text-xs leading-relaxed text-gray-400">
                分享你的想法、问题或经验
              </p>
            </div>

            <div
              v-else-if="widget.type === 'custom_html'"
              v-html="widget.data?.html"
            />

            <NuxtLink
              v-else-if="widget.type === 'custom_image' && widget.data?.image_url"
              :to="widget.data?.link_url ?? '#'"
            >
              <img
                :src="widget.data.image_url"
                :alt="widget.title ?? ''"
                class="w-full rounded"
              >
            </NuxtLink>

            <NuxtLink
              v-else-if="widget.type === 'custom_text_link' && widget.data?.url"
              :to="widget.data.url"
              class="text-sm text-blue-600 hover:underline"
            >{{ widget.data?.text ?? widget.title }}</NuxtLink>
          </div>
        </aside>
      </div>
    </main>

    <footer class="border-t bg-white py-6 text-center text-xs text-gray-400">
      <nav
        v-if="footerNav.length"
        class="mb-3 flex flex-wrap items-center justify-center gap-5"
      >
        <NuxtLink
          v-for="m in footerNav"
          :key="m.id"
          :to="m.url"
          class="text-gray-500 hover:text-gray-900"
        >{{ m.label }}</NuxtLink>
      </nav>
      <span v-if="config?.seo?.icp">{{ config.seo.icp }} · </span>
      {{ siteName }}
    </footer>
  </div>
</template>
