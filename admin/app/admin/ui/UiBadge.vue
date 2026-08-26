<script setup lang="ts">
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '~/admin/utils/cn'

const badgeVariants = cva(
  'inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium transition-colors',
  {
    variants: {
      variant: {
        default: 'border-transparent bg-primary text-primary-foreground',
        secondary: 'border-transparent bg-secondary text-secondary-foreground',
        destructive: 'border-transparent bg-destructive/15 text-destructive',
        success: 'border-transparent bg-[var(--success)]/15 text-[var(--success)]',
        warning: 'border-transparent bg-[var(--warning)]/20 text-[color-mix(in_oklab,var(--warning),black_35%)]'
      }
    },
    defaultVariants: { variant: 'default' }
  }
)

const props = defineProps<{
  variant?: VariantProps<typeof badgeVariants>['variant']
  class?: string
}>()
</script>

<template>
  <span :class="cn(badgeVariants({ variant }), props.class)">
    <slot />
  </span>
</template>
