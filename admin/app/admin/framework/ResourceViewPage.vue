<script setup lang="ts">
import type { EntryNode, ResolvedResource } from '~/admin/core/types'
import { ArrowLeftIcon, PencilIcon } from 'lucide-vue-next'

const props = defineProps<{
  resource: ResolvedResource
  id: string
}>()

const panel = getPanel()
const basePath = computed(() => `${panel.path}/${props.resource.name}`)
const allow = useCan()
const canEdit = computed(() => allow(`${props.resource.permissionPrefix}.edit`))

const record = ref<Record<string, unknown> | null>(null)
const loading = ref(true)
const notFound = ref(false)

/* fallback infolist derived from table columns when no infolist is defined */
const entries = computed<EntryNode[]>(() => {
  if (props.resource.infolist) return props.resource.infolist()
  const map: Record<string, EntryNode['kind']> = {
    text: 'text', number: 'text', money: 'money', date: 'date', badge: 'badge', boolean: 'boolean'
  }
  return (props.resource.table?.() ?? [])
    .filter(c => c.meta.kind !== 'actions' && c.name !== '__actions')
    .map(c => ({ kind: map[c.meta.kind] ?? 'text', name: c.name, label: c.label }))
})

onMounted(async () => {
  try {
    record.value = await $fetch<Record<string, unknown>>(`/api/admin/${props.resource.name}/${props.id}`)
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">
          View {{ resource.label }}
        </h1>
        <NuxtLink
          :to="basePath"
          class="text-sm text-muted-foreground hover:text-foreground"
        >
          ← Back to {{ resource.labelPlural.toLowerCase() }}
        </NuxtLink>
      </div>
      <div class="flex items-center gap-2">
        <UiButton
          variant="outline"
          @click="navigateTo(basePath)"
        >
          <ArrowLeftIcon /> Back
        </UiButton>
        <UiButton
          v-if="canEdit"
          @click="navigateTo(`${basePath}/${id}/edit`)"
        >
          <PencilIcon /> Edit
        </UiButton>
      </div>
    </div>

    <UiCard class="p-6">
      <div
        v-if="loading"
        class="py-12 text-center text-muted-foreground"
      >
        Loading…
      </div>
      <div
        v-else-if="notFound || !record"
        class="py-12 text-center text-muted-foreground"
      >
        {{ resource.label }} not found.
      </div>
      <InfolistRenderer
        v-else
        :entries="entries"
        :record="record"
      />
    </UiCard>
  </div>
</template>
