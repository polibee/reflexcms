<script setup lang="ts">
import { cn } from '~/admin/utils/cn'

const props = defineProps<{
  modelValue?: string | number
  placeholder?: string
  rows?: number
  disabled?: boolean
  class?: string
}>()

const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
</script>

<template>
  <textarea
    :value="modelValue"
    :rows="rows || 4"
    :placeholder="placeholder"
    :disabled="disabled"
    :class="cn(
      'flex min-h-16 w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50',
      props.class
    )"
    @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  />
</template>
