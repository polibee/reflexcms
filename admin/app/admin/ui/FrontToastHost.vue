<script setup lang="ts">
/* Global host for public-site toasts and confirm dialogs — render once in
 * the public layout. Styled consistently with the modern Tiptap modal. */
import { useFrontUi } from '~/composables/useFrontUi'

const { toasts, confirmState, settleConfirm } = useFrontUi()

const iconFor: Record<string, string> = {
  success: 'M9 12l2 2 4-4',
  error: 'M12 8v4m0 4h.01',
  info: 'M12 8v4m0 4h.01'
}
const toneFor: Record<string, string> = {
  success: 'text-green-600 bg-green-50',
  error: 'text-red-600 bg-red-50',
  info: 'text-blue-600 bg-blue-50'
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') { e.preventDefault(); settleConfirm(true) }
  if (e.key === 'Escape') { e.preventDefault(); settleConfirm(false) }
}
</script>

<template>
  <!-- toasts -->
  <div class="pointer-events-none fixed right-4 top-4 z-[70] flex w-80 flex-col gap-2">
    <TransitionGroup
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="translate-x-6 opacity-0"
      enter-to-class="translate-x-0 opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-for="t in toasts"
        :key="t.id"
        class="pointer-events-auto flex items-start gap-2.5 rounded-xl border border-gray-100 bg-white p-3.5 shadow-lg"
      >
        <span
          class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full"
          :class="toneFor[t.type] ?? toneFor.info"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-3 w-3"
            aria-hidden="true"
          >
            <path :d="iconFor[t.type] ?? iconFor.info" />
          </svg>
        </span>
        <p class="text-sm leading-5 text-gray-700">
          {{ t.message }}
        </p>
      </div>
    </TransitionGroup>
  </div>

  <!-- confirm dialog -->
  <Teleport to="body">
    <div
      v-if="confirmState"
      class="fixed inset-0 z-[60] flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm"
      @click.self="settleConfirm(false)"
    >
      <div
        class="w-full max-w-sm rounded-xl border border-gray-100 bg-white p-5 shadow-2xl"
        role="dialog"
        aria-modal="true"
        tabindex="-1"
        @keydown="onKeydown"
      >
        <div class="flex items-start gap-3">
          <span
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full"
            :class="confirmState.danger ? 'bg-red-50 text-red-600' : 'bg-blue-50 text-blue-600'"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-4.5 w-4.5"
              aria-hidden="true"
            >
              <path d="M12 9v4m0 4h.01M10.3 3.9 1.8 18a2 2 0 0 0 1.7 3h17a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z" />
            </svg>
          </span>
          <div class="min-w-0">
            <h3 class="text-sm font-semibold text-gray-900">
              {{ confirmState.title }}
            </h3>
            <p
              v-if="confirmState.message"
              class="mt-1 text-sm leading-5 text-gray-500"
            >
              {{ confirmState.message }}
            </p>
          </div>
        </div>
        <div class="mt-5 flex justify-end gap-2">
          <button
            type="button"
            class="rounded-lg border border-gray-200 px-3.5 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50"
            @click="settleConfirm(false)"
          >
            {{ confirmState.cancelLabel ?? '取消' }}
          </button>
          <button
            type="button"
            class="rounded-lg px-3.5 py-1.5 text-sm text-white transition-colors"
            :class="confirmState.danger ? 'bg-red-600 hover:bg-red-700' : 'bg-blue-600 hover:bg-blue-700'"
            @click="settleConfirm(true)"
          >
            {{ confirmState.confirmLabel ?? '确定' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
