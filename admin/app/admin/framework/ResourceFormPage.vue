<script setup lang="ts">
import type { ResolvedResource } from '~/admin/core/types'
import { ArrowLeftIcon } from 'lucide-vue-next'

const props = defineProps<{
  resource: ResolvedResource
  mode: 'create' | 'edit'
  id?: string
}>()

const panel = getPanel()
const basePath = computed(() => `${panel.path}/${props.resource.name}`)
const schema = computed(() => props.resource.form?.() ?? [])

const { form, submit, submitting } = useFormSchema({
  schema: () => schema.value,
  initialValues: {},
  onSubmit: async (values) => {
    // resource-level hook: pack form values into the API payload shape
    const payload = props.resource.transformOut
      ? props.resource.transformOut(values)
      : values
    if (props.mode === 'create') {
      await $fetch(`/api/admin/${props.resource.name}`, { method: 'POST', body: payload })
      notify(`${props.resource.label} created`)
    } else {
      await $fetch(`/api/admin/${props.resource.name}/${props.id}`, {
        method: 'PUT',
        body: payload
      })
      notify(`${props.resource.label} updated`)
    }
    await navigateTo(basePath.value)
  }
})

/* edit mode: hydrate the form once the record arrives; create mode picks up
 * the resource's defaultValues hook (e.g. sort = current max + 1) */
const loading = ref(props.mode === 'edit')

onMounted(async () => {
  if (props.mode === 'create') {
    if (props.resource.defaultValues) {
      try {
        form.setValues(await props.resource.defaultValues() as never)
      } catch { /* defaults are best-effort */ }
    }
    return
  }
  if (!props.id) return
  try {
    let record = await $fetch<Record<string, unknown>>(`/api/admin/${props.resource.name}/${props.id}`)
    if (props.resource.transformIn) record = props.resource.transformIn(record)
    form.setValues(record as never)
  } catch (e: unknown) {
    notifyError(`Failed to load ${props.resource.label}`, (e as Error).message)
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
          {{ mode === 'create' ? `New ${resource.label}` : `Edit ${resource.label}` }}
        </h1>
        <NuxtLink
          :to="basePath"
          class="text-sm text-muted-foreground hover:text-foreground"
        >
          ← Back to {{ resource.labelPlural.toLowerCase() }}
        </NuxtLink>
      </div>
    </div>

    <UiCard class="p-6">
      <div
        v-if="loading"
        class="py-12 text-center text-muted-foreground"
      >
        Loading…
      </div>
      <form
        v-else
        class="space-y-8"
        @submit.prevent="submit"
      >
        <FormSchemaRenderer :schema="schema" />
        <div class="flex items-center justify-end gap-2 border-t pt-6">
          <UiButton
            type="button"
            variant="outline"
            :disabled="submitting"
            @click="navigateTo(basePath)"
          >
            <ArrowLeftIcon /> Cancel
          </UiButton>
          <UiButton
            type="submit"
            :disabled="submitting"
          >
            {{ submitting ? 'Saving…' : mode === 'create' ? `Create ${resource.label}` : 'Save Changes' }}
          </UiButton>
        </div>
      </form>
    </UiCard>
  </div>
</template>
