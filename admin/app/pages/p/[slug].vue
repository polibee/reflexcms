<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Public static page (隐私政策 / TOS / 关于我们…), authored in the admin
 * page manager with the rich-text editor. Content is admin-authored HTML. */

interface StaticPage {
  title: string
  slug: string
  content: string
  updated_at: string
}

const route = useRoute()
const page = ref<StaticPage | null>(null)
const loading = ref(true)

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

onMounted(async () => {
  try {
    page.value = await rawFetch(`/api/v1/pages/${route.params.slug}`) as StaticPage
  } catch { /* silent */ }
  loading.value = false
})

useHead(() => ({ title: page.value?.title ?? '页面' }))
</script>

<template>
  <div class="mx-auto max-w-3xl">
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <article
      v-else-if="page"
      class="rounded-xl border bg-white p-8"
    >
      <h1 class="text-2xl font-bold text-gray-900">
        {{ page.title }}
      </h1>
      <p class="mt-1 text-xs text-gray-400">
        更新于 {{ page.updated_at?.slice(0, 10) }}
      </p>
      <!-- Admin-authored rich text: trusted source (admin panel only). -->
      <div
        class="prose prose-gray mt-6 max-w-none text-base leading-relaxed text-gray-700"
        v-html="page.content"
      />
    </article>

    <p
      v-else
      class="py-20 text-center text-gray-400"
    >
      页面不存在或未发布
    </p>
  </div>
</template>
