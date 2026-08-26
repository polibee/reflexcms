import type { NavGroup } from '../core/types'

/** Build permission-filtered navigation groups from the registry. */
export function buildNavigation(): NavGroup[] {
  const allow = useCan()
  return getNavGroups()
    .map(group => ({
      ...group,
      items: group.items.filter(item => allow(item.permission))
    }))
    .filter(group => group.items.length > 0)
}
