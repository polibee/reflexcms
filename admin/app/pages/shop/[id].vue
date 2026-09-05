<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Product detail + checkout. Guests CAN buy: the unguessable order number
 * is kept in sessionStorage so the purchased invite code stays retrievable
 * without an account. */

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

const route = useRoute()
const router = useRouter()

const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, unknown>) => Promise<unknown>

const product = ref<ShopProduct | null>(null)
const gateways = ref<string[]>([])
const loading = ref(true)

const gateway = ref('')
const submitting = ref(false)
const errorMsg = ref('')

const GATEWAY_LABELS: Record<string, string> = {
  xcash: 'Xcash（加密支付）',
  coinpayments: 'CoinPayments（加密支付）',
  nowpayments: 'NOWPayments（加密支付）',
  xunhupay: '虎皮椒（微信/支付宝）',
  codepay: '码支付（个人收款）',
  paypal: 'PayPal'
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

    // Guests remember their orders via sessionStorage.
    try {
      const mine = JSON.parse(sessionStorage.getItem('shop_orders') ?? '[]') as string[]
      mine.push(res.order_no)
      sessionStorage.setItem('shop_orders', JSON.stringify(mine.slice(-20)))
    } catch { /* storage unavailable */ }

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
    // /api/v1/products only exposes active in-stock goods; detail falls
    // back gracefully for off-shelf items already on screen.
    const res = await rawFetch('/api/v1/products') as { items: ShopProduct[], gateways: string[] }
    product.value = (res.items ?? []).find(p => String(p.id) === String(route.params.id)) ?? null
    gateways.value = res.gateways ?? []
    if (gateways.value.length && !gateway.value) gateway.value = gateways.value[0]
  } catch { /* silent */ }
  loading.value = false
})

useHead(() => ({ title: product.value ? `${product.value.title} · 商品详情` : '商品详情' }))
</script>

<template>
  <div class="mx-auto max-w-3xl">
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <template v-else-if="product">
      <div class="overflow-hidden rounded-xl border bg-white">
        <div class="h-56 w-full overflow-hidden">
          <div
            v-if="product.type === 'invite'"
            class="flex h-full w-full items-center justify-center bg-gradient-to-br from-emerald-50 to-teal-100 text-5xl"
          >
            🎟
          </div>
          <div
            v-else
            class="flex h-full w-full items-center justify-center bg-gradient-to-br from-amber-50 to-orange-100 text-5xl"
          >
            🛍
          </div>
        </div>
        <div class="p-6">
          <div class="flex flex-wrap items-center gap-2">
            <span
              v-if="product.type === 'invite'"
              class="rounded bg-emerald-100 px-2 py-0.5 text-xs font-medium text-emerald-700"
            >邀请码商品 · 购买后自动发放注册邀请码</span>
            <span class="text-xs text-gray-400">库存 {{ product.stock }}</span>
            <span
              v-if="product.grant_points"
              class="rounded bg-amber-50 px-2 py-0.5 text-xs text-amber-600"
            >赠送 {{ product.grant_points }} 积分</span>
          </div>
          <h1 class="mt-2 text-2xl font-bold text-gray-900">
            {{ product.title }}
          </h1>
          <p class="mt-3 whitespace-pre-wrap text-sm leading-relaxed text-gray-600">
            {{ product.description }}
          </p>
          <p class="mt-4 text-2xl font-bold text-orange-600">
            {{ fmtPrice(product) }}
          </p>
        </div>
      </div>

      <div class="mt-5 rounded-xl border bg-white p-6">
        <h2 class="text-sm font-semibold text-gray-900">
          选择支付方式
        </h2>
        <div class="mt-3 space-y-2">
          <label
            v-for="g in gateways"
            :key="g"
            class="flex cursor-pointer items-center gap-2 rounded-lg border p-3 text-sm hover:bg-gray-50"
            :class="gateway === g ? 'border-blue-500 ring-1 ring-blue-400' : 'border-gray-200'"
          >
            <input
              v-model="gateway"
              type="radio"
              name="gateway"
              :value="g"
              class="accent-blue-600"
            >
            {{ GATEWAY_LABELS[g] ?? g }}
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
          {{ submitting ? '正在创建订单…' : `立即购买 ${fmtPrice(product)}` }}
        </button>
        <p class="mt-2 text-xs text-gray-400">
          无需注册账号即可购买；订单凭据保存在本浏览器，支付完成后回到本页即可查看商品。
        </p>
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
