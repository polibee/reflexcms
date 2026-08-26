<script setup lang="ts">
/* Windowed pagination bar for public list pages: prev/next plus a page
 * window around the current page. Emits the target page; the parent owns
 * the query state. */
const props = withDefaults(defineProps<{
  page: number
  totalPages: number
  window?: number
}>(), { window: 2 })

const emit = defineEmits<{ change: [page: number] }>()

const pages = computed<number[]>(() => {
  const total = props.totalPages
  const cur = props.page
  const start = Math.max(1, Math.min(cur - props.window, total - props.window * 2))
  const end = Math.min(total, start + props.window * 2)
  const list: number[] = []
  for (let i = start; i <= end; i++) list.push(i)
  return list
})

function go(p: number) {
  if (p < 1 || p > props.totalPages || p === props.page) return
  emit('change', p)
}
</script>

<template>
  <nav
    v-if="totalPages > 1"
    class="mt-6 flex items-center justify-center gap-1"
  >
    <button
      type="button"
      class="rounded-md border px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40"
      :disabled="page <= 1"
      @click="go(page - 1)"
    >
      上一页
    </button>

    <button
      v-if="pages[0] > 1"
      type="button"
      class="rounded-md px-3 py-1.5 text-sm text-gray-500 hover:bg-gray-50"
      @click="go(1)"
    >
      1
    </button>
    <span
      v-if="pages[0] > 2"
      class="px-1 text-sm text-gray-400"
    >…</span>

    <button
      v-for="p in pages"
      :key="p"
      type="button"
      class="min-w-9 rounded-md px-3 py-1.5 text-sm"
      :class="p === page
        ? 'bg-blue-600 font-medium text-white'
        : 'text-gray-600 hover:bg-gray-50'"
      @click="go(p)"
    >
      {{ p }}
    </button>

    <span
      v-if="pages[pages.length - 1] < totalPages - 1"
      class="px-1 text-sm text-gray-400"
    >…</span>
    <button
      v-if="pages[pages.length - 1] < totalPages"
      type="button"
      class="rounded-md px-3 py-1.5 text-sm text-gray-500 hover:bg-gray-50"
      @click="go(totalPages)"
    >
      {{ totalPages }}
    </button>

    <button
      type="button"
      class="rounded-md border px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40"
      :disabled="page >= totalPages"
      @click="go(page + 1)"
    >
      下一页
    </button>
  </nav>
</template>
