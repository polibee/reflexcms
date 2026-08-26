<script setup lang="ts">
import { cn } from '~/admin/utils/cn'

const props = defineProps<{ class?: string }>()
</script>

<template>
  <label :class="cn('text-sm font-medium leading-none', props.class)">
    <slot />
  </label>
</template>
