<script setup lang="ts">
import { cn } from '~/admin/utils/cn'

const panel = getPanel()
const allow = useCan()

const widgets = computed(() =>
  getWidgets().filter(w => !w.permission || allow(w.permission))
)

function spanClass(span: number | undefined): string {
  switch (span) {
    case 1: return 'xl:col-span-1'
    case 2: return 'xl:col-span-2'
    case 3: return 'xl:col-span-3'
    case 4: return 'xl:col-span-4'
    default: return 'xl:col-span-4'
  }
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">
        Dashboard
      </h1>
      <p class="text-sm text-muted-foreground">
        Welcome to {{ panel.branding?.name ?? 'your admin panel' }}
      </p>
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
      <div
        v-for="widget in widgets"
        :key="widget.name"
        :class="cn(spanClass(widget.span))"
      >
        <component :is="widget.component" />
      </div>
    </div>

    <p
      v-if="widgets.length === 0"
      class="py-12 text-center text-sm text-muted-foreground"
    >
      No widgets available for your role.
    </p>
  </div>
</template>
