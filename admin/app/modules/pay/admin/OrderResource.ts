export default defineResource({
  name: 'orders',
  model: 'Order',
  label: '订单',
  labelPlural: '订单',
  icon: 'receipt',
  group: '商城',
  sort: 10,
  searchable: ['order_no', 'title'],

  // Manual delivery confirmation for offline payments and re-delivery of
  // missing invite codes (idempotent server-side).
  rowActions: [
    defineAction({
      name: 'mark-paid',
      label: '标记已支付',
      icon: 'check-circle',
      permission: 'orders.edit',
      visible: record => record.status === 'pending',
      confirm: {
        title: '确认该订单已线下支付？将触发库存扣减与商品（邀请码）发放。',
        confirmLabel: '标记已支付'
      },
      handler: async ({ record }) => {
        await $fetch('/api/admin/order-mark-paid', {
          method: 'POST',
          body: { order_id: record!.id }
        })
        notify('订单已标记为已支付')
        emitAdminEvent('orders:refresh')
      }
    })
  ],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('order_no', '订单号', { sortable: true }),
    textColumn('title', '商品'),
    numberColumn('amount_cents', '金额（分）', { sortable: true }),
    textColumn('currency', '货币'),
    textColumn('gateway', '渠道', { sortable: true }),
    textColumn('status', '状态', { sortable: true }),
    dateColumn('paid_at', '支付时间'),
    dateColumn('created_at', 'Created', { sortable: true })
  ],

  // Orders are immutable records — no create/edit form, no row mutations.
  form: () => [],

  infolist: () => [
    textEntry('order_no', '订单号'),
    textEntry('title', '商品'),
    moneyCentsEntry('amount_cents', '金额'),
    textEntry('currency', '货币'),
    textEntry('gateway', '渠道'),
    textEntry('gateway_ref', '渠道单号'),
    badgeEntry('status', '状态', {
      pending: { label: '待支付', variant: 'warning' },
      paid: { label: '已支付', variant: 'success' },
      failed: { label: '支付失败', variant: 'destructive' },
      refunded: { label: '已退款', variant: 'secondary' },
      closed: { label: '已关闭', variant: 'secondary' }
    }),
    dateEntry('paid_at', '支付时间'),
    datetimeEntry('created_at', '创建时间')
  ]
})
