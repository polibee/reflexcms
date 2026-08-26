<script setup lang="ts">
import type { SchemaNode } from '~/admin/core/types'
import { useFormSchema } from '~/admin/forms/useFormSchema'

const props = defineProps<{
  schema: SchemaNode[]
  initial?: Record<string, unknown>
  busy?: boolean
}>()

const emit = defineEmits<{ submit: [values: Record<string, unknown>], cancel: [] }>()

const { submit, submitting } = useFormSchema({
  schema: () => props.schema,
  initialValues: props.initial,
  onSubmit: values => emit('submit', values)
})
</script>

<template>
  <form
    class="space-y-6"
    @submit.prevent="submit"
  >
    <FormSchemaRenderer :schema="schema" />
    <div class="flex justify-end gap-2 pt-2">
      <UiButton
        type="button"
        variant="outline"
        :disabled="submitting || busy"
        @click="emit('cancel')"
      >
        Cancel
      </UiButton>
      <UiButton
        type="submit"
        :disabled="submitting || busy"
      >
        {{ submitting || busy ? 'Saving…' : 'Confirm' }}
      </UiButton>
    </div>
  </form>
</template>
