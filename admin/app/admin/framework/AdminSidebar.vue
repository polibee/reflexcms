<script setup lang="ts">
import { getIcon } from '~/admin/utils/iconMap'
import { cn } from '~/admin/utils/cn'

const panel = getPanel()
const route = useRoute()
const navigation = computed(() => buildNavigation())

function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(`${to}/`)
}
</script>

<template>
  <aside class="flex h-full w-64 shrink-0 flex-col border-r bg-sidebar text-sidebar-foreground">
    <!-- brand -->
    <div class="flex h-14 items-center gap-2.5 border-b px-5">
      <div class="flex h-7 w-7 items-center justify-center rounded-lg bg-primary text-primary-foreground">
        <svg
          viewBox="0 0 24 24"
          class="h-4 w-4"
          fill="none"
          stroke="currentColor"
          stroke-width="2.5"
        >
          <path
            d="M4 17l6-6-6-6M12 19h8"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </div>
      <span class="text-sm font-semibold">{{ panel.branding?.name ?? panel.title }}</span>
    </div>

    <!-- navigation groups -->
    <nav class="flex-1 space-y-5 overflow-y-auto p-3">
      <!-- dashboard shortcut -->
      <NuxtLink
        :to="panel.path"
        class="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-sm transition-colors hover:bg-accent hover:text-accent-foreground"
        :class="cn(route.path === panel.path && 'bg-primary/10 font-medium text-primary dark:bg-primary/15')"
      >
        <component
          :is="getIcon('dashboard')"
          class="h-4 w-4 shrink-0"
        />
        Dashboard
      </NuxtLink>

      <div
        v-for="group in navigation"
        :key="group.label"
      >
        <p class="px-2 pb-1.5 text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
          {{ group.label }}
        </p>
        <div class="space-y-0.5">
          <NuxtLink
            v-for="item in group.items"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-sm transition-colors hover:bg-accent hover:text-accent-foreground"
            :class="cn(isActive(item.to) && 'bg-primary/10 font-medium text-primary dark:bg-primary/15')"
          >
            <component
              :is="getIcon(item.icon)"
              class="h-4 w-4 shrink-0"
            />
            {{ item.label }}
            <span
              v-if="item.badge !== undefined"
              class="ml-auto rounded-full bg-primary px-1.5 py-0.5 text-[10px] font-semibold text-primary-foreground"
            >
              {{ item.badge }}
            </span>
          </NuxtLink>
        </div>
      </div>
    </nav>

    <div class="border-t p-3 text-[11px] text-muted-foreground">
      Nuxt Admin
    </div>
  </aside>
</template>
