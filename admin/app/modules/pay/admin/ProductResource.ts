export default defineResource({
  name: 'products',
  model: 'Product',
  label: '商品',
  labelPlural: '商品',
  icon: 'shopping-bag',
  group: '商城',
  sort: 5,
  searchable: ['title'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('title', '标题'),
    numberColumn('price_cents', '价格（分）', { sortable: true }),
    textColumn('currency', '货币'),
    numberColumn('stock', '库存', { sortable: true }),
    numberColumn('grant_points', '赠送积分'),
    booleanColumn('is_active', '上架'),
    dateColumn('created_at', 'Created')
  ],

  form: () => [
    section('商品信息', [
      textInput('title', '商品标题', { required: true, colSpan: 2 }),
      textarea('description', '商品描述', { rows: 3, colSpan: 2 }),
      grid(2, [
        numberInput('price_cents', '价格（最小货币单位/分）', { required: true, min: 0, helpText: '例如 $9.99 填 999' }),
        textInput('currency', '货币代码', { placeholder: 'USD / CNY', defaultValue: 'USD' }),
        numberInput('stock', '库存', { required: true, min: 0 }),
        numberInput('grant_points', '购买后赠送积分', { min: 0, helpText: '0 表示不赠送' }),
        textInput('image', '封面图 URL', { colSpan: 2 }),
        switchInput('is_active', '上架')
      ]),
      numberInput('sort', '排序', { min: 0, helpText: '已自动接上当前最大序号，可自行修改' })
    ])
  ],

  defaultValues: async () => {
    try {
      const res = await $fetch<{ items: Array<Record<string, unknown>> }>('/api/admin/products', { query: { perPage: 200 } })
      const max = (res.items ?? []).reduce((m, w) => Math.max(m, Number(w.sort ?? 0)), 0)
      return { sort: max + 1 }
    } catch {
      return { sort: 100 }
    }
  },

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('title', '标题'),
    textEntry('description', '描述'),
    textEntry('price_cents', '价格（分）'),
    textEntry('currency', '货币'),
    textEntry('stock', '库存'),
    textEntry('grant_points', '赠送积分'),
    textEntry('is_active', '上架')
  ]
})
