<script setup lang="ts">
import { XIcon, CheckCircle2Icon, AlertCircleIcon, InfoIcon } from 'lucide-vue-next'

const toasts = useToasts()

const iconFor = {
  success: CheckCircle2Icon,
  error: AlertCircleIcon,
  info: InfoIcon
} as const

const colorFor = {
  success: 'text-[var(--success)]',
  error: 'text-destructive',
  info: 'text-primary'
} as const
</script>

<template>
  <Teleport to="body">
    <div class="pointer-events-none fixed bottom-4 right-4 z-[100] flex w-full max-w-sm flex-col gap-2">
      <TransitionGroup
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="opacity-0 translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-for="toast in toasts.items"
          :key="toast.id"
          class="pointer-events-auto flex items-start gap-3 rounded-lg border bg-card p-4 shadow-lg"
        >
          <component
            :is="iconFor[toast.variant]"
            class="mt-0.5 h-4 w-4 shrink-0"
            :class="colorFor[toast.variant]"
          />
          <div class="min-w-0 flex-1">
            <p class="text-sm font-medium">
              {{ toast.title }}
            </p>
            <p
              v-if="toast.description"
              class="mt-0.5 text-xs text-muted-foreground break-words"
            >
              {{ toast.description }}
            </p>
          </div>
          <button
            class="text-muted-foreground transition-colors hover:text-foreground"
            @click="toasts.dismiss(toast.id)"
          >
            <XIcon class="h-4 w-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
