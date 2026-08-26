export default defineResource({
  name: 'menus',
  model: 'Menu',
  label: '导航菜单',
  labelPlural: '导航菜单',
  icon: 'menu',
  group: '布局管理',
  sort: 20,
  searchable: ['label', 'url'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('location', '位置'),
    textColumn('label', '菜单名'),
    textColumn('url', '链接'),
    numberColumn('sort', '排序', { sortable: true }),
    booleanColumn('is_visible', '可见'),
    dateColumn('created_at', 'Created')
  ],

  form: () => [
    section('菜单项', [
      grid(2, [
        textInput('label', '菜单名称', { required: true }),
        textInput('url', '链接地址', { required: true, placeholder: '/articles' })
      ]),
      grid(2, [
        selectInput('location', '显示位置', [
          { label: '顶部导航', value: 'top_nav' },
          { label: '底部页脚', value: 'footer_nav' }
        ], { required: true }),
        numberInput('sort', '排序', { min: 0, defaultValue: 100 })
      ]),
      switchInput('is_visible', '是否显示', { defaultValue: true })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('location', '位置'),
    textEntry('label', '菜单名'),
    textEntry('url', '链接'),
    textEntry('sort', '排序')
  ]
})
