export const useUiStore = defineStore('admin-ui', () => {
  // persisted UI preferences only - never server data
  const theme = useCookie<'light' | 'dark'>('admin-theme', { default: () => 'light' })
  const sidebarOpen = useCookie<boolean>('admin-sidebar-open', { default: () => true })

  function toggleTheme(): void {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
  }

  function toggleSidebar(): void {
    sidebarOpen.value = !sidebarOpen.value
  }

  return { theme, sidebarOpen, toggleTheme, toggleSidebar }
})
