<script setup lang="ts">
import { getCoreRowModel, useVueTable, type ColumnDef } from '@tanstack/vue-table'
import {
  ChevronDownIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronUpIcon,
  ChevronsUpDownIcon,
  DownloadIcon,
  MoreHorizontalIcon,
  SearchIcon
} from 'lucide-vue-next'
import {
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuRoot,
  DropdownMenuTrigger
} from 'reka-ui'
import type { ActionDef, BadgeVariant, ColumnDefLite, ResolvedResource } from '~/admin/core/types'
import type { ResourceTable, RowRecord } from '~/admin/tables/useResourceTable'
import { exportCsv, formatDate, formatMoney, formatNumber } from '~/admin/utils/format'
import { getIcon } from '~/admin/utils/iconMap'

const props = withDefaults(defineProps<{
  resource: ResolvedResource
  columns: ColumnDefLite[]
  state: ResourceTable
  bulkActions?: ActionDef[]
}>(), { bulkActions: () => [] })

const runner = useActionRunner()

/* ---------- TanStack core row model over lite columns ---------- */

interface LiteMeta { lite: ColumnDefLite }

const tanColumns = computed<ColumnDef<RowRecord>[]>(() =>
  props.columns.map((lite) => {
    if (lite.name === '__actions') {
      return {
        id: lite.name,
        header: '',
        enableSorting: false,
        meta: { lite }
      } as ColumnDef<RowRecord>
    }
    return {
      id: lite.name,
      accessorFn: (row: RowRecord) => row[lite.name],
      header: lite.label,
      enableSorting: !!lite.meta.sortable,
      meta: { lite }
    } as ColumnDef<RowRecord>
  })
)

const table = useVueTable<RowRecord>({
  get data() { return props.state.items.value },
  get columns() { return tanColumns.value },
  getCoreRowModel: getCoreRowModel(),
  manualPagination: true,
  manualSorting: true
})

function liteOf(id: string): ColumnDefLite | undefined {
  return (table.getColumn(id)?.columnDef.meta as LiteMeta | undefined)?.lite
}

/* ---------- cell rendering ---------- */

function cellText(lite: ColumnDefLite, value: unknown): string {
  const m = lite.meta
  switch (m.kind) {
    case 'money': return formatMoney(value, m.prefix)
    case 'number': return formatNumber(value)
    case 'date': return formatDate(value)
    case 'boolean': return value ? (m.booleanLabels?.true ?? 'Yes') : (m.booleanLabels?.false ?? 'No')
    case 'badge': return m.badges?.[value as string | number]?.label ?? String(value ?? '-')
    default: return `${m.prefix ?? ''}${String(value ?? '-')}${m.suffix ?? ''}`
  }
}

function badgeVariant(lite: ColumnDefLite, value: unknown): BadgeVariant | undefined {
  return lite.meta.badges?.[value as string | number]?.variant
}

function visibleRowActions(lite: ColumnDefLite, row: RowRecord): ActionDef[] {
  return (lite.meta.actions ?? []).filter(a => !a.visible || a.visible(row))
}

function runRowAction(action: ActionDef, row: RowRecord): void {
  runner.run(action, { record: row, resource: props.resource })
}

function runBulkAction(action: ActionDef): void {
  runner.run(action, { ids: props.state.selectedIds.value, resource: props.resource })
}

const hasSelectionColumn = computed(() => props.bulkActions.length > 0)

/* search box binds through the composable's mutation API */
const searchQuery = computed<string>({
  get: () => props.state.q.value,
  set: v => props.state.setQuery(v)
})

function onExport(): void {
  exportCsv(props.resource.name, props.state.items.value, props.columns)
}
</script>

<template>
  <div class="space-y-3">
    <!-- toolbar -->
    <div class="flex flex-wrap items-center gap-2">
      <div class="relative">
        <SearchIcon class="absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <UiInput
          v-model="searchQuery"
          :placeholder="`Search ${resource.labelPlural?.toLowerCase()}…`"
          class="w-64 pl-8"
        />
      </div>
      <slot name="filters" />

      <div class="ml-auto flex items-center gap-2">
        <slot name="actions" />
        <UiButton
          variant="outline"
          size="sm"
          :disabled="state.items.value.length === 0"
          @click="onExport"
        >
          <DownloadIcon /> Export
        </UiButton>
      </div>
    </div>

    <!-- bulk actions bar -->
    <div
      v-if="hasSelectionColumn && state.selection.value.size > 0"
      class="flex items-center gap-2 rounded-lg border bg-muted/50 px-3 py-2 text-sm"
    >
      <span class="font-medium">{{ state.selection.value.size }} selected</span>
      <UiButton
        v-for="action in bulkActions"
        :key="action.name"
        size="sm"
        :variant="action.variant ?? 'outline'"
        @click="runBulkAction(action)"
      >
        <component :is="getIcon(action.icon)" />
        {{ action.label }}
      </UiButton>
      <UiButton
        size="sm"
        variant="ghost"
        class="ml-auto"
        @click="state.clearSelection()"
      >
        Clear
      </UiButton>
    </div>

    <!-- table -->
    <div class="rounded-xl border bg-card">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b bg-muted/40 text-left text-xs uppercase tracking-wide text-muted-foreground">
              <th
                v-if="hasSelectionColumn"
                class="w-10 px-4 py-3"
              >
                <UiCheckbox
                  :model-value="state.allPageSelected.value"
                  @update:model-value="state.toggleSelectAll()"
                />
              </th>
              <th
                v-for="header in table.getHeaderGroups()[0]?.headers ?? []"
                :key="header.id"
                class="whitespace-nowrap px-4 py-3 font-medium"
              >
                <button
                  v-if="header.column.getCanSort() && liteOf(header.id)?.meta.sortable"
                  class="inline-flex items-center gap-1 hover:text-foreground"
                  @click="state.toggleSort(liteOf(header.id)!)"
                >
                  {{ header.column.columnDef.header }}
                  <ChevronsUpDownIcon
                    v-if="state.sortBy.value !== header.id"
                    class="h-3.5 w-3.5"
                  />
                  <ChevronUpIcon
                    v-else-if="state.sortDir.value === 'asc'"
                    class="h-3.5 w-3.5"
                  />
                  <ChevronDownIcon
                    v-else
                    class="h-3.5 w-3.5"
                  />
                </button>
                <span v-else>{{ header.column.columnDef.header }}</span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="state.loading.value && state.items.value.length === 0">
              <td
                :colspan="props.columns.length + (hasSelectionColumn ? 1 : 0)"
                class="px-4 py-12 text-center text-muted-foreground"
              >
                Loading…
              </td>
            </tr>
            <tr v-else-if="state.items.value.length === 0">
              <td
                :colspan="props.columns.length + (hasSelectionColumn ? 1 : 0)"
                class="px-4 py-12 text-center text-muted-foreground"
              >
                No {{ resource.labelPlural?.toLowerCase() }} found.
              </td>
            </tr>

            <tr
              v-for="row in table.getRowModel().rows"
              v-else
              :key="row.id"
              class="border-b transition-colors last:border-0 hover:bg-muted/30"
              :class="state.isSelected(row.original.id) ? 'bg-muted/40' : ''"
            >
              <td
                v-if="hasSelectionColumn"
                class="px-4 py-3"
              >
                <UiCheckbox
                  :model-value="state.isSelected(row.original.id)"
                  @update:model-value="state.toggleSelect(row.original)"
                />
              </td>
              <td
                v-for="cell in row.getVisibleCells()"
                :key="cell.id"
                class="px-4 py-3"
                :class="(cell.column.columnDef.meta as LiteMeta)?.lite?.name === '__actions' ? 'text-right' : ''"
              >
                <template v-if="(cell.column.columnDef.meta as LiteMeta)?.lite?.name === '__actions'">
                  <DropdownMenuRoot v-if="visibleRowActions((cell.column.columnDef.meta as LiteMeta).lite!, row.original).length > 0">
                    <DropdownMenuTrigger as-child>
                      <UiButton
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8"
                      >
                        <MoreHorizontalIcon class="h-4 w-4" />
                      </UiButton>
                    </DropdownMenuTrigger>
                    <DropdownMenuPortal>
                      <DropdownMenuContent
                        align="end"
                        class="z-50 min-w-36 rounded-md border bg-popover p-1 shadow-md"
                      >
                        <DropdownMenuItem
                          v-for="action in visibleRowActions((cell.column.columnDef.meta as LiteMeta).lite!, row.original)"
                          :key="action.name"
                          class="flex cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent"
                          @select="runRowAction(action, row.original)"
                        >
                          <component
                            :is="getIcon(action.icon)"
                            class="h-4 w-4"
                            :class="action.variant === 'destructive' ? 'text-destructive' : ''"
                          />
                          <span :class="action.variant === 'destructive' ? 'text-destructive' : ''">{{ action.label }}</span>
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenuPortal>
                  </DropdownMenuRoot>
                </template>

                <UiBadge
                  v-else-if="(cell.column.columnDef.meta as LiteMeta)?.lite?.meta.kind === 'badge'"
                  :variant="(badgeVariant((cell.column.columnDef.meta as LiteMeta).lite!, cell.getValue()) as never)"
                >
                  {{ cellText((cell.column.columnDef.meta as LiteMeta).lite!, cell.getValue()) }}
                </UiBadge>

                <span
                  v-else
                  :class="(cell.column.columnDef.meta as LiteMeta)?.lite?.meta.align === 'right' ? 'tabular-nums' : ''"
                >{{ cellText((cell.column.columnDef.meta as LiteMeta).lite!, cell.getValue()) }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- pagination -->
      <div class="flex flex-wrap items-center justify-between gap-3 border-t px-4 py-3 text-sm text-muted-foreground">
        <div>
          Showing
          <template v-if="state.total.value > 0">
            {{ (state.page.value - 1) * state.perPage.value + 1 }}–{{ Math.min(state.page.value * state.perPage.value, state.total.value) }}
            of <span class="font-medium text-foreground">{{ state.total.value }}</span>
          </template>
          <template v-else>
            0
          </template>
          results
        </div>
        <div class="flex items-center gap-3">
          <UiSelect
            :model-value="state.perPage.value"
            :options="[10, 25, 50].map(n => ({ label: `${n} / page`, value: n }))"
            class="h-8 w-28 text-xs"
            @update:model-value="state.setPerPage(Number($event))"
          />
          <div class="flex items-center gap-1">
            <UiButton
              variant="outline"
              size="icon"
              class="h-8 w-8"
              :disabled="state.page.value <= 1"
              @click="state.prevPage()"
            >
              <ChevronLeftIcon class="h-4 w-4" />
            </UiButton>
            <span class="min-w-16 text-center text-xs">Page {{ state.page.value }} / {{ state.totalPages.value }}</span>
            <UiButton
              variant="outline"
              size="icon"
              class="h-8 w-8"
              :disabled="state.page.value >= state.totalPages.value"
              @click="state.nextPage()"
            >
              <ChevronRightIcon class="h-4 w-4" />
            </UiButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
