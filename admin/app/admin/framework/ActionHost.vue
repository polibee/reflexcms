<script setup lang="ts">
import {
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle
} from 'reka-ui'
import { getIcon } from '~/admin/utils/iconMap'

const runner = useActionRunner()

const ConfirmIcon = computed(() => getIcon(runner.confirmTarget?.action.icon))
const isDestructive = computed(() => runner.confirmTarget?.action.variant === 'destructive')
</script>

<template>
  <!-- confirmation dialog -->
  <DialogRoot
    :open="!!runner.confirmTarget"
    @update:open="v => !v && runner.close()"
  >
    <DialogPortal v-if="runner.confirmTarget">
      <DialogOverlay class="fixed inset-0 z-50 bg-black/60" />
      <DialogContent
        class="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl border bg-card p-6 shadow-lg focus:outline-none"
      >
        <div class="flex items-start gap-4">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full"
            :class="isDestructive ? 'bg-destructive/10 text-destructive' : 'bg-primary/10 text-primary'"
          >
            <component
              :is="ConfirmIcon"
              class="h-5 w-5"
            />
          </div>
          <div class="space-y-1.5">
            <DialogTitle class="text-base font-semibold">
              {{ runner.confirmTarget.action.confirm?.title }}
            </DialogTitle>
            <DialogDescription class="text-sm text-muted-foreground">
              {{ runner.confirmTarget.action.confirm?.description || `This will run "${runner.confirmTarget.action.label}".` }}
            </DialogDescription>
          </div>
        </div>
        <div class="mt-6 flex justify-end gap-2">
          <UiButton
            variant="outline"
            :disabled="runner.busy"
            @click="runner.close()"
          >
            Cancel
          </UiButton>
          <UiButton
            :variant="isDestructive ? 'destructive' : 'default'"
            :disabled="runner.busy"
            @click="runner.execute(runner.confirmTarget!.action, runner.confirmTarget!.ctx)"
          >
            {{ runner.confirmTarget.action.confirm?.confirmLabel || runner.confirmTarget.action.label }}
          </UiButton>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>

  <!-- action form modal -->
  <DialogRoot
    :open="!!runner.formTarget"
    @update:open="v => !v && runner.close()"
  >
    <DialogPortal v-if="runner.formTarget">
      <DialogOverlay class="fixed inset-0 z-50 bg-black/60" />
      <DialogContent
        class="fixed left-1/2 top-1/2 z-50 max-h-[85vh] w-full max-w-lg -translate-x-1/2 -translate-y-1/2 overflow-y-auto rounded-xl border bg-card p-6 shadow-lg focus:outline-none"
      >
        <DialogTitle class="text-base font-semibold">
          {{ runner.formTarget.action.label }}
        </DialogTitle>
        <div class="mt-4">
          <ActionFormModal
            :key="runner.formTarget.action.name"
            :schema="runner.formTarget.action.form!()"
            :initial="runner.formTarget.ctx.record"
            :busy="runner.busy"
            @submit="values => runner.execute(runner.formTarget!.action, runner.formTarget!.ctx, values)"
            @cancel="runner.close()"
          />
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
