/* i18n: modular language packs. The admin picks the site language
 * (app.locale setting, surfaced through /api/v1/site/config as `locale`);
 * the public frontend reads it and serves UI chrome strings from the
 * matching pack. Admin-authored content (menus, widgets) is user data —
 * the admin writes it in the chosen language directly. */

export type Locale = 'zh-CN' | 'en'

import en from '../locales/en'
import zhCN from '../locales/zh-CN'

const packs: Record<Locale, Record<string, string>> = { 'zh-CN': zhCN, en }

export function useLocale(): { locale: Ref<Locale>, t: (key: string) => string } {
  const { config } = useSiteConfig()
  const locale = computed<Locale>(() => {
    const l = config.value?.locale as Locale | undefined
    return l === 'en' ? 'en' : 'zh-CN'
  })
  const t = (key: string): string => packs[locale.value][key] ?? packs['zh-CN'][key] ?? key
  return { locale, t }
}
