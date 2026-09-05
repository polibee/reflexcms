<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Shop: admin-configurable product listing. The whole module can be
 * toggled off from the admin panel (404 when disabled). */

interface ShopProduct {
  id: number
  title: string
  description: string
  price_cents: number
  currency: string
  stock: number
  grant_points: number
  type?: string
}

const { config } = useSiteConfig()
const { get } = usePublicApi()

const items = ref<ShopProduct[]>([])
const shopEnabled = computed(() => (config.value?.sections ?? []).includes('shop'))
const loading = ref(true)

const filter = ref<'all' | 'invite'>('all')
const filtered = computed(() =>
  filter.value === 'invite'
    ? items.value.filter(p => p.type === 'invite')
    : items.value)

function fmtPrice(p: ShopProduct): string {
  const symbol = p.currency === 'CNY' ? '¥' : p.currency === 'USD' ? '$' : p.currency + ' '
  return symbol + (p.price_cents / 100).toFixed(2)
}

onMounted(async () => {
  if (!config.value) {
    try {
      config.value = await get<SiteConfig>('site/config')
    } catch { /* silent */ }
  }
  if (shopEnabled.value) {
    try {
      const res = await get<{ items: ShopProduct[] }>('products')
      items.value = res.items ?? []
    } catch { /* silent */ }
  }
  loading.value = false
})

useHead({ title: '商城' })
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <template v-else-if="shopEnabled">
      <div class="mb-6 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">
            商城
          </h1>
          <p class="mt-1 text-sm text-gray-500">
            购买邀请码与社区服务，支付完成后自动发放
          </p>
        </div>
        <div class="flex items-center gap-3">
          <div class="flex rounded-lg border border-gray-200 bg-white p-1 text-sm">
            <button
              type="button"
              class="rounded-md px-3 py-1"
              :class="filter === 'all' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-50'"
              @click="filter = 'all'"
            >
              全部
            </button>
            <button
              type="button"
              class="rounded-md px-3 py-1"
              :class="filter === 'invite' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-50'"
              @click="filter = 'invite'"
            >
              邀请码
            </button>
          </div>
          <NuxtLink
            to="/shop/my"
            class="text-sm text-gray-500 hover:text-gray-900"
          >我的订单</NuxtLink>
        </div>
      </div>

      <div
        v-if="filtered.length"
        class="grid gap-5 md:grid-cols-2 lg:grid-cols-3"
      >
        <NuxtLink
          v-for="p in filtered"
          :key="p.id"
          :to="`/shop/${p.id}`"
          class="group overflow-hidden rounded-xl border bg-white transition-shadow hover:shadow-md"
        >
          <div class="h-36 w-full overflow-hidden">
            <div
              v-if="p.type === 'invite'"
              class="flex h-full w-full items-center justify-center bg-gradient-to-br from-emerald-50 to-teal-100 text-4xl"
            >
              🎟
            </div>
            <div
              v-else
              class="flex h-full w-full items-center justify-center bg-gradient-to-br from-amber-50 to-orange-100 text-4xl"
            >
              🛍
            </div>
          </div>
          <div class="p-4">
            <div class="flex items-center gap-2">
              <span
                v-if="p.type === 'invite'"
                class="rounded bg-emerald-100 px-1.5 py-0.5 text-xs text-emerald-700"
              >邀请码</span>
              <h2 class="truncate font-semibold text-gray-900 group-hover:text-blue-600">
                {{ p.title }}
              </h2>
            </div>
            <p class="mt-1 line-clamp-2 text-sm text-gray-500">
              {{ p.description }}
            </p>
            <div class="mt-3 flex items-center justify-between">
              <span class="text-lg font-bold text-orange-600">{{ fmtPrice(p) }}</span>
              <span
                v-if="p.grant_points"
                class="rounded bg-amber-50 px-2 py-0.5 text-xs text-amber-600"
              >送 {{ p.grant_points }} 积分</span>
            </div>
            <p class="mt-2 text-xs text-gray-400">
              库存 {{ p.stock }}
            </p>
          </div>
        </NuxtLink>
      </div>

      <div
        v-else
        class="rounded-xl border bg-white p-16 text-center"
      >
        <p class="text-gray-400">
          暂无相关商品
        </p>
      </div>
    </template>

    <p
      v-else
      class="py-20 text-center text-gray-400"
    >
      商城未开放
    </p>
  </div>
</template>
