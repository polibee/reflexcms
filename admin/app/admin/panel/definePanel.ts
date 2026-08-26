import type { PanelConfig } from '../core/types'

/** Panels are isolated admin spaces: /admin today, /partner tomorrow. */
export function definePanel(config: PanelConfig): PanelConfig {
  return config
}
