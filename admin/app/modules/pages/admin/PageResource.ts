export default defineResource({
  name: 'pages',
  model: 'Page',
  label: '静态页面',
  labelPlural: '静态页面',
  icon: 'file-text',
  group: '页面管理',
  sort: 5,
  searchable: ['title', 'slug'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('title', '标题'),
    textColumn('slug', '路径（/p/slug）'),
    textColumn('status', '状态'),
    dateColumn('created_at', 'Created')
  ],

  form: () => [
    section('页面内容', [
      grid(2, [
        textInput('title', '页面标题', { required: true, placeholder: '隐私政策 / 服务条款 / 关于我们' }),
        textInput('slug', '访问路径', {
          required: true,
          placeholder: 'privacy（前台地址 /p/privacy，留空按标题自动生成）'
        })
      ]),
      wysiwyg('content', '页面内容（富文本）', { colSpan: 2 }),
      selectInput('status', '状态', [
        { label: '已发布', value: 'published' },
        { label: '草稿', value: 'draft' }
      ], { defaultValue: 'draft', helpText: '仅已发布页面对前台可见' })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('title', '标题'),
    textEntry('slug', '路径'),
    textEntry('status', '状态'),
    textEntry('updated_at', '更新时间')
  ]
})
