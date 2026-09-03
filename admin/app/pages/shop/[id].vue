<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Checkout: pick a gateway for one product, create the order, then jump to
 * the gateway's hosted payment page. */

interface ShopProduct {
  id: number
  title: string
  description: string
  price_cents: number
  currency: string
  stock: number
  grant_points: number
}

interface ShopPayload {
  items: ShopProduct[]
  aff_links: Record<string, string>
  gateways: string[]
}

const route = useRoute()
const router = useRouter()
const { user: me } = useAuthUser()

const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, unknown>) => Promise<unknown>

const product = ref<ShopProduct | null>(null)
const gateways = ref<string[]>([])
const affLinks = ref<Record<string, string>>({})
const loading = ref(true)

const gateway = ref('')
const submitting = ref(false)
const errorMsg = ref('')

const GATEWAY_LABELS: Record<string, string> = {
  xcash: 'Xcash（加密支付）',
  coinpayments: 'CoinPayments（加密支付）',
  xunhupay: '虎皮椒（微信/支付宝）',
  codepay: '码支付（个人收款）',
  paypal: 'PayPal'
}

const GATEWAY_REGISTER: Record<string, string> = {
  xcash: '注册 Xcash ↗',
  coinpayments: '注册 CoinPayments ↗',
  xunhupay: '注册虎皮椒 ↗',
  codepay: '注册码支付 ↗',
  paypal: '注册 PayPal ↗'
}

function fmtPrice(p: ShopProduct): string {
  const symbol = p.currency === 'CNY' ? '¥' : p.currency === 'USD' ? '$' : p.currency + ' '
  return symbol + (p.price_cents / 100).toFixed(2)
}

async function submitOrder() {
  if (!gateway.value) {
    errorMsg.value = '请选择支付方式'
    return
  }
  submitting.value = true
  errorMsg.value = ''
  try {
    const res = await rawFetch('/api/v1/orders', {
      method: 'POST',
      body: { product_id: product.value?.id, gateway: gateway.value }
    }) as { order_no: string, pay_url: string }
    if (res.pay_url) {
      window.location.href = res.pay_url
      return
    }
    await navigateTo(`/shop/orders/${res.order_no}`)
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string, message?: string } }
    errorMsg.value = err.data?.statusMessage ?? err.data?.message ?? '下单失败，请重试'
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const res = await rawFetch('/api/v1/products') as ShopPayload
    product.value = (res.items ?? []).find(p => String(p.id) === String(route.params.id)) ?? null
    gateways.value = res.gateways ?? []
    affLinks.value = res.aff_links ?? {}
    if (gateways.value.length && !gateway.value) gateway.value = gateways.value[0]
  } catch { /* silent */ }
  loading.value = false
})

useHead(() => ({ title: product.value ? `购买 ${product.value.title}` : '商品' }))
</script>

<template>
  <div class="mx-auto max-w-2xl">
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <template v-else-if="product">
      <div class="rounded-xl border bg-white p-6">
        <h1 class="text-xl font-bold text-gray-900">
          {{ product.title }}
        </h1>
        <p class="mt-2 text-sm text-gray-500">
          {{ product.description }}
        </p>
        <div class="mt-4 flex items-center justify-between border-t pt-4">
          <span class="text-2xl font-bold text-orange-600">{{ fmtPrice(product) }}</span>
          <span
            v-if="product.grant_points"
            class="rounded bg-amber-50 px-2 py-1 text-xs text-amber-600"
          >购买后赠送 {{ product.grant_points }} 积分</span>
        </div>
      </div>

      <template v-if="me">
        <div class="mt-5 rounded-xl border bg-white p-6">
          <h2 class="text-sm font-semibold text-gray-900">
            选择支付方式
          </h2>
          <div class="mt-3 space-y-2">
            <label
              v-for="g in gateways"
              :key="g"
              class="flex cursor-pointer items-center justify-between rounded-lg border p-3 text-sm hover:bg-gray-50"
              :class="gateway === g ? 'border-blue-500 ring-1 ring-blue-400' : 'border-gray-200'"
            >
              <span class="flex items-center gap-2">
                <input
                  v-model="gateway"
                  type="radio"
                  name="gateway"
                  :value="g"
                  class="accent-blue-600"
                >
                {{ GATEWAY_LABELS[g] ?? g }}
              </span>
              <a
                v-if="affLinks[g]"
                :href="affLinks[g]"
                target="_blank"
                rel="noopener noreferrer"
                class="text-xs text-blue-600 hover:underline"
                @click.stop
              >{{ GATEWAY_REGISTER[g] ?? '注册' }}</a>
            </label>
          </div>

          <p
            v-if="errorMsg"
            class="mt-3 text-sm text-red-600"
          >
            {{ errorMsg }}
          </p>

          <button
            type="button"
            :disabled="submitting || !gateway"
            class="mt-4 w-full rounded-lg bg-blue-600 py-2.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            @click="submitOrder"
          >
            {{ submitting ? '正在创建订单…' : `支付 ${fmtPrice(product)}` }}
          </button>
        </div>
      </template>

      <div
        v-else
        class="mt-5 rounded-xl border bg-white p-8 text-center"
      >
        <p class="text-gray-500">
          登录后即可购买
        </p>
        <NuxtLink
          :to="`/login?redirect=/shop/${product.id}`"
          class="mt-4 inline-block rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >登录</NuxtLink>
      </div>

      <button
        type="button"
        class="mt-4 text-sm text-gray-400 hover:text-gray-600"
        @click="router.back()"
      >
        ← 返回商城
      </button>
    </template>

    <p
      v-else
      class="py-20 text-center text-gray-400"
    >
      商品不存在或已下架
    </p>
  </div>
</template>
