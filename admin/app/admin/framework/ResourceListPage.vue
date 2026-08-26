<script setup lang="ts">
import type { ActionDef, ResolvedResource } from '~/admin/core/types'
import { PlusIcon } from 'lucide-vue-next'

const props = defineProps<{ resource: ResolvedResource }>()

const allow = useCan()
const panel = getPanel()

const prefix = computed(() => props.resource.permissionPrefix)
const basePath = computed(() => `${panel.path}/${props.resource.name}`)
const canCreate = computed(() => allow(`${prefix.value}.create`))
const canEdit = computed(() => allow(`${prefix.value}.edit`))
const canDelete = computed(() => allow(`${prefix.value}.delete`))

const state = useResourceTable(props.resource)
const columns = computed(() => props.resource.table?.() ?? [])

/* module actions (publish / mark-shipped) trigger refresh via the event bus */
const offRefresh = onAdminEvent(`${props.resource.name}:refresh`, () => {
  state.fetch()
})
onUnmounted(offRefresh)

/* ---------------- default actions ---------------- */

const deleteAction = computed<ActionDef>(() => ({
  name: 'delete',
  label: 'Delete',
  icon: 'trash',
  variant: 'destructive',
  permission: `${prefix.value}.delete`,
  confirm: {
    title: `Delete this ${props.resource.label.toLowerCase()}?`,
    description: 'This action cannot be undone.',
    confirmLabel: 'Delete'
  },
  handler: async ({ record, resource }) => {
    await $fetch(`/api/admin/${resource.name}/${record!.id}`, { method: 'DELETE' })
    notify(`${resource.label} deleted`)
    await state.refresh()
  }
}))

const editAction = computed<ActionDef>(() => ({
  name: 'edit',
  label: 'Edit',
  icon: 'pencil',
  permission: `${prefix.value}.edit`,
  handler: async ({ record }) => {
    await navigateTo(`${basePath.value}/${record!.id}/edit`)
  }
}))

const rowActions = computed<ActionDef[]>(() => [
  ...(props.resource.rowActions ?? []),
  ...(canEdit.value ? [editAction.value] : []),
  ...(canDelete.value ? [deleteAction.value] : [])
])

const bulkDeleteAction = computed<ActionDef>(() => ({
  name: 'bulk-delete',
  label: 'Delete Selected',
  icon: 'trash',
  variant: 'destructive',
  permission: `${prefix.value}.delete`,
  confirm: {
    title: `Delete ${state.selection.value.size || 'selected'} selected ${props.resource.labelPlural.toLowerCase()}?`,
    description: 'This action cannot be undone.',
    confirmLabel: 'Delete'
  },
  handler: async ({ ids, resource }) => {
    await $fetch(`/api/admin/${resource.name}/bulk-delete`, {
      method: 'POST',
      body: { ids }
    })
    notify(`${ids!.length} ${resource.labelPlural.toLowerCase()} deleted`)
    await state.refresh()
  }
}))

const bulkActions = computed<ActionDef[]>(() => [
  ...(props.resource.bulkActions ?? []),
  ...(!props.resource.bulkActions?.some(a => a.name === 'bulk-delete') && canDelete.value
    ? [bulkDeleteAction.value]
    : [])
])

const tableColumns = computed(() => {
  const base = columns.value.filter(c => c.meta.kind !== 'actions')
  const declaredActions = columns.value.find(c => c.meta.kind === 'actions')?.meta.actions ?? []
  return [
    ...base,
    {
      name: '__actions',
      label: '',
      meta: { kind: 'actions' as const, align: 'right' as const, actions: [...declaredActions, ...rowActions.value] }
    }
  ]
})

function openCreate(): void {
  navigateTo(`${basePath.value}/create`)
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">
          {{ resource.labelPlural }}
        </h1>
        <p class="text-sm text-muted-foreground">
          Manage {{ resource.labelPlural.toLowerCase() }}
        </p>
      </div>
      <UiButton
        v-if="canCreate"
        @click="openCreate"
      >
        <PlusIcon /> New {{ resource.label }}
      </UiButton>
    </div>

    <DataTable
      :resource="resource"
      :columns="tableColumns"
      :state="state"
      :bulk-actions="bulkActions"
    >
      <template #filters>
        <slot name="filters" />
      </template>
    </DataTable>
  </div>
</template>
