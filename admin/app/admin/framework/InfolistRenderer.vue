<script setup lang="ts">
import type { EntryNode } from '~/admin/core/types'
import { formatDate, formatDateTime, formatMoney } from '~/admin/utils/format'

const props = defineProps<{
  entries: EntryNode[]
  record: Record<string, unknown>
}>()

function display(entry: EntryNode): string {
  const raw = props.record[entry.name]
  switch (entry.kind) {
    case 'money': return formatMoney(raw, entry.prefix)
    case 'date': return formatDate(raw)
    case 'datetime': return formatDateTime(raw)
    case 'boolean': return raw ? 'Yes' : 'No'
    case 'badge': {
      const style = entry.badges?.[raw as string | number]
      return style ? style.label : String(raw ?? '-')
    }
    default: return `${entry.prefix ?? ''}${String(raw ?? '-')}${entry.suffix ?? ''}`
  }
}

function badgeVariantFor(entry: EntryNode): string | undefined {
  const raw = props.record[entry.name]
  return entry.badges?.[raw as string | number]?.variant
}
</script>

<template>
  <dl class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-5">
    <div
      v-for="entry in entries"
      :key="entry.name"
      class="space-y-1"
    >
      <dt class="text-xs font-medium uppercase tracking-wide text-muted-foreground">
        {{ entry.label }}
      </dt>
      <dd class="text-sm">
        <UiBadge
          v-if="entry.kind === 'badge'"
          :variant="(badgeVariantFor(entry) as never)"
        >
          {{ display(entry) }}
        </UiBadge>
        <template v-else-if="entry.kind === 'link'">
          <a
            v-if="record[entry.name]"
            :href="String(record[entry.name])"
            target="_blank"
            rel="noopener"
            class="text-primary underline-offset-4 hover:underline break-all"
          >{{ display(entry) }}</a>
          <span
            v-else
            class="text-muted-foreground"
          >-</span>
        </template>
        <span v-else>{{ display(entry) }}</span>
      </dd>
    </div>
  </dl>
</template>
