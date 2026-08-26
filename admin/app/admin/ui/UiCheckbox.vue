<script setup lang="ts">
defineProps<{
  modelValue?: boolean
  disabled?: boolean
  label?: string
}>()

const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>

<template>
  <label class="inline-flex cursor-pointer items-center gap-2">
    <input
      type="checkbox"
      :checked="modelValue"
      :disabled="disabled"
      class="h-4 w-4 shrink-0 rounded border-input accent-[var(--primary)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
      @change="emit('update:modelValue', ($event.target as HTMLInputElement).checked)"
    >
    <span
      v-if="label"
      class="text-sm"
    >{{ label }}</span>
  </label>
</template>
