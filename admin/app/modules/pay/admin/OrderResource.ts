export default defineResource({
  name: 'orders',
  model: 'Order',
  label: '订单',
  labelPlural: '订单',
  icon: 'receipt',
  group: '商城',
  sort: 10,
  searchable: ['order_no', 'title'],

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
    textEntry('amount_cents', '金额（分）'),
    textEntry('currency', '货币'),
    textEntry('gateway', '渠道'),
    textEntry('gateway_ref', '渠道单号'),
    textEntry('status', '状态'),
    textEntry('paid_at', '支付时间'),
    textEntry('created_at', '创建时间')
  ]
})
