<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* Order status page: guests land here after payment (return_url) and —
 * for invite-code goods — see the granted code with one-click copy. */

interface OrderView {
  order_no: string
  title: string
  amount_cents: number
  currency: string
  gateway: string
  status: string
  paid_at: string | null
  granted_code?: string
  created_at: string
}

const route = useRoute()
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, unknown>) => Promise<unknown>

const order = ref<OrderView | null>(null)
const notFound = ref(false)
const copied = ref(false)
let poll: ReturnType<typeof setInterval> | null = null

const GATEWAY_LABELS: Record<string, string> = {
  xcash: 'Xcash（加密支付）',
  coinpayments: 'CoinPayments（加密支付）',
  nowpayments: 'NOWPayments（加密支付）',
  xunhupay: '虎皮椒',
  codepay: '码支付',
  paypal: 'PayPal'
}

function fmtPrice(o: OrderView): string {
  const symbol = o.currency === 'CNY' ? '¥' : o.currency === 'USD' ? '$' : o.currency + ' '
  return symbol + (o.amount_cents / 100).toFixed(2)
}

async function fetchOrder() {
  try {
    order.value = await rawFetch(`/api/v1/orders/${route.params.no}`) as OrderView
    // stop polling once paid/failed
    if (order.value && order.value.status !== 'pending' && poll) {
      clearInterval(poll)
      poll = null
    }
  } catch {
    notFound.value = true
    if (poll) {
      clearInterval(poll)
      poll = null
    }
  }
}

async function copyCode() {
  if (!order.value?.granted_code) return
  try {
    await navigator.clipboard.writeText(order.value.granted_code)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch { /* clipboard unavailable */ }
}

onMounted(() => {
  void fetchOrder()
  poll = setInterval(() => {
    if (order.value?.status === 'pending') void fetchOrder()
  }, 5000)
})

onBeforeUnmount(() => {
  if (poll) clearInterval(poll)
})

useHead(() => ({ title: `订单 ${order.value?.order_no ?? ''}` }))
</script>

<template>
  <div class="mx-auto max-w-xl">
    <div
      v-if="notFound"
      class="rounded-xl border bg-white p-16 text-center"
    >
      <p class="text-gray-400">
        订单不存在
      </p>
      <NuxtLink
        to="/shop"
        class="mt-4 inline-block text-sm text-blue-600 hover:underline"
      >← 返回商城</NuxtLink>
    </div>

    <div
      v-else-if="order"
      class="overflow-hidden rounded-xl border bg-white"
    >
      <div class="p-6">
        <div class="flex items-center justify-between">
          <h1 class="text-lg font-bold text-gray-900">
            订单详情
          </h1>
          <span
            class="rounded-full px-3 py-1 text-xs font-medium"
            :class="{
              'bg-yellow-100 text-yellow-700': order.status === 'pending',
              'bg-green-100 text-green-700': order.status === 'paid',
              'bg-red-100 text-red-600': order.status === 'failed'
            }"
          >
            {{ order.status === 'pending' ? '待支付' : order.status === 'paid' ? '已支付' : '已失败' }}
          </span>
        </div>

        <div class="mt-4 space-y-2 text-sm text-gray-600">
          <p>{{ order.title }}</p>
          <p class="text-lg font-bold text-orange-600">
            {{ fmtPrice(order) }}
          </p>
          <p class="text-xs text-gray-400">
            订单号 {{ order.order_no }} · {{ GATEWAY_LABELS[order.gateway] ?? order.gateway }}
          </p>
        </div>

        <!-- invite-code delivery -->
        <div
          v-if="order.granted_code"
          class="mt-5 rounded-lg border border-emerald-200 bg-emerald-50 p-4"
        >
          <p class="text-sm font-medium text-emerald-700">
            🎟 你的邀请码已发放
          </p>
          <div class="mt-2 flex items-center gap-2">
            <code class="flex-1 rounded bg-white px-3 py-2 font-mono text-sm tracking-wider text-gray-900">{{ order.granted_code }}</code>
            <button
              type="button"
              class="rounded-md bg-emerald-600 px-3 py-2 text-xs font-medium text-white hover:bg-emerald-700"
              @click="copyCode"
            >
              {{ copied ? '✓ 已复制' : '复制' }}
            </button>
          </div>
          <p class="mt-2 text-xs text-emerald-600">
            请妥善保存，注册时填写此邀请码
          </p>
        </div>

        <p
          v-if="order.status === 'pending'"
          class="mt-4 text-xs text-gray-400"
        >
          支付完成后本页会自动更新（每 5 秒刷新一次）。
        </p>
      </div>

      <div class="border-t p-4 text-center">
        <NuxtLink
          to="/shop"
          class="text-sm text-blue-600 hover:underline"
        >← 返回商城</NuxtLink>
      </div>
    </div>

    <div
      v-else
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>
  </div>
</template>
