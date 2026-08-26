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
    <button
      type="button"
      role="switch"
      :aria-checked="modelValue"
      :disabled="disabled"
      class="relative inline-flex h-5 w-9 shrink-0 items-center rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
      :class="modelValue ? 'bg-primary' : 'bg-input'"
      @click="emit('update:modelValue', !modelValue)"
    >
      <span
        class="pointer-events-none block h-4 w-4 rounded-full bg-background shadow-lg transition-transform"
        :class="modelValue ? 'translate-x-4' : 'translate-x-0'"
      />
    </button>
    <span
      v-if="label"
      class="text-sm"
    >{{ label }}</span>
  </label>
</template>
