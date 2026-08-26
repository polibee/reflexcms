<script setup lang="ts">
import { FileTextIcon, MessagesSquareIcon, ReplyIcon, ShieldAlertIcon, UsersIcon } from 'lucide-vue-next'

interface Stats {
  usersTotal: number
  usersActive: number
  articlesTotal: number
  articlesPublished: number
  topicsTotal: number
  repliesTotal: number
  pendingComments: number
}

const stats = ref<Stats | null>(null)

onMounted(async () => {
  try {
    stats.value = await $fetch<Stats>('/api/admin/stats')
  } catch {
    stats.value = null
  }
})

const cards = computed(() => [
  { label: 'Users', value: stats.value?.usersTotal, sub: `${stats.value?.usersActive ?? 0} active`, icon: UsersIcon },
  { label: 'Articles', value: stats.value?.articlesTotal, sub: `${stats.value?.articlesPublished ?? 0} published`, icon: FileTextIcon },
  { label: 'Topics', value: stats.value?.topicsTotal, icon: MessagesSquareIcon },
  { label: 'Replies', value: stats.value?.repliesTotal, icon: ReplyIcon },
  { label: 'Pending Comments', value: stats.value?.pendingComments, icon: ShieldAlertIcon }
])
</script>

<template>
  <div class="grid grid-cols-2 gap-4 xl:grid-cols-5">
    <UiCard
      v-for="card in cards"
      :key="card.label"
      class="p-5"
    >
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium text-muted-foreground">{{ card.label }}</span>
        <component
          :is="card.icon"
          class="h-4 w-4 text-muted-foreground"
        />
      </div>
      <p class="mt-2 text-2xl font-semibold tabular-nums">
        {{ card.value ?? '—' }}
      </p>
      <p
        v-if="card.sub"
        class="mt-1 text-xs text-muted-foreground"
      >
        {{ card.sub }}
      </p>
    </UiCard>
  </div>
</template>
