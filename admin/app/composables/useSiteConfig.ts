interface SiteSection {
  label: string
  to: string
}

export interface SiteConfig {
  name: string
  mode: string
  home?: 'cms' | 'forum' | 'both'
  locale?: string
  seo: { description: string, icp: string }
  sections: string[]
}

export function useSiteConfig() {
  const config = useState<SiteConfig | null>('site-config', () => null)
  const siteName = computed(() => config.value?.name ?? 'ReflexCMS')
  const navItems = computed<SiteSection[]>(() => {
    const sections = config.value?.sections ?? []
    const items: SiteSection[] = [{ label: '首页', to: '/' }]
    if (sections.includes('articles')) items.push({ label: '文章', to: '/articles' })
    if (sections.includes('topics')) items.push({ label: '论坛', to: '/forums' })
    if (sections.includes('shop')) items.push({ label: '商城', to: '/shop' })
    return items
  })
  return { config, siteName, navItems }
}
