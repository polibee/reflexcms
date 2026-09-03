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
  image: string
  grant_points: number
}

interface ShopPayload {
  items: ShopProduct[]
  currency: string
  aff_links: Record<string, string>
  gateways: string[]
}

const { config } = useSiteConfig()
const { get } = usePublicApi()

const items = ref<ShopProduct[]>([])
const gateways = ref<string[]>([])
const affLinks = ref<Record<string, string>>({})
const shopEnabled = computed(() => (config.value?.sections ?? []).includes('shop'))
const loading = ref(true)

const GATEWAY_LABELS: Record<string, string> = {
  xcash: 'Xcash（加密支付）',
  coinpayments: 'CoinPayments（加密支付）',
  xunhupay: '虎皮椒（微信/支付宝）',
  codepay: '码支付（个人收款）',
  paypal: 'PayPal'
}

const GATEWAY_REGISTER: Record<string, string> = {
  xcash: '注册 Xcash',
  coinpayments: '注册 CoinPayments',
  xunhupay: '注册虎皮椒',
  codepay: '注册码支付',
  paypal: '注册 PayPal'
}

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
      const res = await get<ShopPayload>('products')
      items.value = res.items ?? []
      gateways.value = res.gateways ?? []
      affLinks.value = res.aff_links ?? {}
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
      <h1 class="mb-2 text-2xl font-bold text-gray-900">
        商城
      </h1>
      <p class="mb-6 text-sm text-gray-500">
        购买积分套餐与社区服务，支付完成后积分自动到账
      </p>

      <div
        v-if="items.length"
        class="grid gap-5 md:grid-cols-2 lg:grid-cols-3"
      >
        <NuxtLink
          v-for="p in items"
          :key="p.id"
          :to="`/shop/${p.id}`"
          class="group overflow-hidden rounded-xl border bg-white transition-shadow hover:shadow-md"
        >
          <div class="h-36 w-full overflow-hidden">
            <img
              v-if="p.image"
              :src="p.image"
              :alt="p.title"
              class="h-full w-full object-cover"
            >
            <div
              v-else
              class="flex h-full w-full items-center justify-center bg-gradient-to-br from-amber-50 to-orange-100 text-3xl"
            >
              🛍
            </div>
          </div>
          <div class="p-4">
            <h2 class="font-semibold text-gray-900 group-hover:text-blue-600">
              {{ p.title }}
            </h2>
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
          商店暂无在售商品
        </p>
      </div>

      <div
        v-if="gateways.length"
        class="mt-8 rounded-xl border bg-white p-5 text-sm"
      >
        <h2 class="font-semibold text-gray-900">
          支持的支付渠道
        </h2>
        <div class="mt-2 flex flex-wrap gap-2">
          <span
            v-for="g in gateways"
            :key="g"
            class="rounded-md bg-gray-100 px-2.5 py-1 text-xs text-gray-600"
          >{{ GATEWAY_LABELS[g] ?? g }}</span>
        </div>
        <p class="mt-3 text-xs text-gray-400">
          还没有支付渠道账号？
          <a
            v-for="(link, g) in affLinks"
            v-show="link"
            :key="g"
            :href="link"
            target="_blank"
            rel="noopener noreferrer"
            class="mr-3 text-blue-600 hover:underline"
          >{{ GATEWAY_REGISTER[g] ?? g }} ↗</a>
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
