<script setup lang="ts">
definePageMeta({ layout: 'public' })

interface Topic { id: number, title: string, reply_count: number, like_count: number, last_reply_at: string }
const route = useRoute()
const boardId = route.params.id as string
const topics = ref<Topic[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await $fetch<{ items: Topic[] }>(`/api/v1/topics?forum_category_id=${boardId}&perPage=50`)
    topics.value = res.items ?? []
  } catch { /* silent */ }
  loading.value = false
})
</script>

<template>
  <div>
    <h1 class="mb-6 text-2xl font-bold">
      板块帖子
    </h1>
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>
    <div
      v-else
      class="space-y-2"
    >
      <NuxtLink
        v-for="t in topics"
        :key="t.id"
        :to="`/topics/${t.id}`"
        class="block rounded-lg border bg-white p-4 hover:shadow-sm"
      >
        <div class="flex items-center justify-between">
          <span class="font-medium text-gray-900">{{ t.title }}</span>
          <span class="text-xs text-gray-400">{{ t.reply_count }} 回复 · {{ t.like_count }} 赞</span>
        </div>
      </NuxtLink>
      <p
        v-if="!topics.length"
        class="py-10 text-center text-gray-400"
      >
        暂无帖子
      </p>
    </div>
  </div>
</template>
