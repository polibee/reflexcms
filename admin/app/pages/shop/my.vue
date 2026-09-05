<script setup lang="ts">
definePageMeta({ layout: 'public' })

/* My orders: for signed-in users lists could come from the API, but the
 * primary flow is guests — whose orders are tracked by unguessable order
 * numbers kept in sessionStorage. */

interface OrderView {
  order_no: string
  title: string
  amount_cents: number
  currency: string
  status: string
  granted_code?: string
}

const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, unknown>) => Promise<unknown>

const orders = ref<OrderView[]>([])
const loading = ref(true)

onMounted(async () => {
  let nos: string[] = []
  try {
    nos = JSON.parse(sessionStorage.getItem('shop_orders') ?? '[]') as string[]
  } catch { /* ignore */ }
  const results = await Promise.allSettled(
    nos.map(no => rawFetch(`/api/v1/orders/${no}`) as Promise<OrderView>)
  )
  orders.value = results
    .filter((r): r is PromiseFulfilledResult<OrderView> => r.status === 'fulfilled')
    .map(r => r.value)
  loading.value = false
})

function fmtPrice(o: OrderView): string {
  const symbol = o.currency === 'CNY' ? '¥' : o.currency === 'USD' ? '$' : o.currency + ' '
  return symbol + (o.amount_cents / 100).toFixed(2)
}

const statusText: Record<string, string> = {
  pending: '待支付',
  paid: '已支付',
  failed: '已失败'
}
</script>

<template>
  <div class="mx-auto max-w-2xl">
    <h1 class="mb-5 text-2xl font-bold text-gray-900">
      我的订单
    </h1>

    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <div
      v-else-if="orders.length"
      class="divide-y divide-gray-100 rounded-xl border bg-white"
    >
      <NuxtLink
        v-for="o in orders"
        :key="o.order_no"
        :to="`/shop/orders/${o.order_no}`"
        class="flex items-center justify-between px-5 py-3.5 hover:bg-gray-50"
      >
        <div class="min-w-0">
          <p class="truncate text-sm font-medium text-gray-900">{{ o.title }}</p>
          <p class="text-xs text-gray-400">{{ o.order_no }}</p>
        </div>
        <div class="ml-3 shrink-0 text-right">
          <p class="text-sm font-semibold text-orange-600">{{ fmtPrice(o) }}</p>
          <p class="text-xs text-gray-400">{{ statusText[o.status] ?? o.status }}</p>
        </div>
      </NuxtLink>
    </div>

    <div
      v-else
      class="rounded-xl border bg-white p-16 text-center"
    >
      <p class="text-gray-400">
        本浏览器暂无订单记录
      </p>
      <NuxtLink
        to="/shop"
        class="mt-4 inline-block text-sm text-blue-600 hover:underline"
      >去商城看看 →</NuxtLink>
    </div>
  </div>
</template>
