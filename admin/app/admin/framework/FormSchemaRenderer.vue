<script setup lang="ts">
import type { FieldNode, SchemaNode } from '~/admin/core/types'
import { cn } from '~/admin/utils/cn'

const props = withDefaults(defineProps<{
  schema: SchemaNode[]
  inGrid?: boolean
}>(), { inGrid: false })

function spanClass(node: SchemaNode): string {
  if (node.type !== 'field') return ''
  switch (node.colSpan) {
    case 2: return 'col-span-2'
    case 3: return 'col-span-3'
    case 4: return 'col-span-4'
    default: return 'col-span-1'
  }
}

function gridClass(columns: number): string {
  if (columns === 1) return 'grid grid-cols-1 gap-4'
  if (columns === 2) return 'grid grid-cols-1 md:grid-cols-2 gap-4'
  if (columns === 3) return 'grid grid-cols-1 md:grid-cols-3 gap-4'
  return 'grid grid-cols-1 md:grid-cols-4 gap-4'
}
</script>

<template>
  <div :class="cn(!inGrid && 'space-y-6')">
    <template
      v-for="(node, i) in props.schema"
      :key="i"
    >
      <!-- section layout -->
      <section
        v-if="node.type === 'section'"
        class="space-y-4"
      >
        <div v-if="node.title || node.description">
          <h3
            v-if="node.title"
            class="text-sm font-semibold"
          >
            {{ node.title }}
          </h3>
          <p
            v-if="node.description"
            class="text-xs text-muted-foreground mt-0.5"
          >
            {{ node.description }}
          </p>
        </div>
        <FormSchemaRenderer :schema="node.children" />
      </section>

      <!-- grid layout -->
      <div
        v-else-if="node.type === 'grid'"
        :class="gridClass(node.columns)"
      >
        <FormSchemaRenderer
          :schema="(node as { children: SchemaNode[] }).children"
          in-grid
        />
      </div>

      <!-- field leaf -->
      <div
        v-else
        :class="props.inGrid ? spanClass(node as FieldNode) : ''"
      >
        <FormField :node="node as FieldNode" />
      </div>
    </template>
  </div>
</template>
