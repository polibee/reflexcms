import type {
  AdminResource,
  ModuleDef,
  NavGroup,
  NavItem,
  PanelConfig,
  ResolvedResource,
  WidgetDef
} from './types'

/* =============================================================
 * AdminRegistry - the service-container of the framework.
 * Modules register resources/widgets; Panel + pages consume them.
 * ============================================================= */

interface RegistryState {
  panel: PanelConfig
  resources: Map<string, AdminResource>
  widgets: Map<string, WidgetDef>
  modules: Map<string, ModuleDef>
  navGroups: Map<string, number>
}

const state: RegistryState = {
  panel: {
    id: 'admin',
    path: '/admin',
    title: 'Admin',
    auth: true,
    perPage: 10
  },
  resources: new Map(),
  widgets: new Map(),
  modules: new Map(),
  navGroups: new Map()
}

export function setPanel(panel: Partial<PanelConfig>): void {
  state.panel = { ...state.panel, ...panel }
}

export function registerResource(resource: AdminResource): void {
  if (state.resources.has(resource.name)) {
    console.warn(`[admin] duplicate resource "${resource.name}" ignored`)
    return
  }
  state.resources.set(resource.name, resource)
  const group = resource.group || 'General'
  if (!state.navGroups.has(group)) state.navGroups.set(group, 100)
}

export function registerWidget(widget: WidgetDef): void {
  if (state.widgets.has(widget.name)) return
  state.widgets.set(widget.name, widget)
}

export function registerModule(module: ModuleDef): void {
  if (state.modules.has(module.name)) return
  state.modules.set(module.name, module)
  for (const g of module.navGroups ?? []) {
    if (g.sort !== undefined) state.navGroups.set(g.label, g.sort)
  }
  for (const r of module.resources ?? []) registerResource(r)
  for (const w of module.widgets ?? []) registerWidget(w)
}

function resolve(resource: AdminResource): ResolvedResource {
  return {
    ...resource,
    labelPlural: resource.labelPlural ?? resource.label,
    permissionPrefix: resource.permissionPrefix ?? resource.name
  }
}

/* ---------------- read accessors ---------------- */

export function getPanel(): PanelConfig {
  return state.panel
}

export function getResources(): ResolvedResource[] {
  return [...state.resources.values()]
    .map(resolve)
    .sort((a, b) => (a.sort ?? 100) - (b.sort ?? 100))
}

export function getResource(name: string): ResolvedResource | undefined {
  const found = state.resources.get(name)
  return found ? resolve(found) : undefined
}

export function getWidgets(): WidgetDef[] {
  return [...state.widgets.values()].sort(
    (a, b) => (a.order ?? 100) - (b.order ?? 100)
  )
}

export function getNavGroups(): NavGroup[] {
  const groups = new Map<string, NavItem[]>()
  for (const resource of getResources()) {
    const label = resource.group || 'General'
    if (!groups.has(label)) groups.set(label, [])
    groups.get(label)!.push({
      label: resource.labelPlural!,
      to: `${getPanel().path}/${resource.name}`,
      icon: resource.icon,
      permission: `${resolve(resource).permissionPrefix}.view`
    })
  }
  return [...groups.entries()]
    .map(([label, items]) => ({
      label,
      sort: state.navGroups.get(label) ?? 100,
      items
    }))
    .sort((a, b) => a.sort - b.sort)
}
