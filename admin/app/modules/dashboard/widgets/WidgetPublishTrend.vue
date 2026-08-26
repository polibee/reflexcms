<script setup lang="ts">
interface SeriesPoint {
  label: string
  start?: string
  count: number
}

const series = ref<SeriesPoint[]>([])

onMounted(async () => {
  try {
    const stats = await $fetch<{ publishSeries: SeriesPoint[] }>('/api/admin/stats')
    series.value = stats.publishSeries ?? []
  } catch {
    series.value = []
  }
})

const maxCount = computed(() => Math.max(1, ...series.value.map(p => p.count)))
</script>

<template>
  <div class="space-y-2">
    <p
      v-if="series.length === 0"
      class="text-sm text-muted-foreground"
    >
      No data.
    </p>
    <div
      v-for="point in series"
      :key="point.label"
      class="flex items-center gap-3"
    >
      <span class="w-8 shrink-0 text-xs text-muted-foreground">{{ point.label }}</span>
      <div class="h-2 flex-1 overflow-hidden rounded-full bg-muted">
        <div
          class="h-full rounded-full bg-primary transition-all"
          :style="{ width: `${(point.count / maxCount) * 100}%` }"
        />
      </div>
      <span class="w-6 shrink-0 text-right text-xs tabular-nums">{{ point.count }}</span>
    </div>
  </div>
</template>
