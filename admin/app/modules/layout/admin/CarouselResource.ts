export default defineResource({
  name: 'carousels',
  model: 'Carousel',
  label: '轮播图',
  labelPlural: '轮播图',
  icon: 'image',
  group: '布局管理',
  sort: 30,
  searchable: ['title'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('title', '标题'),
    textColumn('image_url', '图片'),
    textColumn('link_url', '链接'),
    numberColumn('sort', '排序', { sortable: true }),
    booleanColumn('is_active', '启用'),
    dateColumn('created_at', 'Created')
  ],

  form: () => [
    section('轮播图', [
      textInput('title', '标题（管理标识）', { required: true }),
      uploadInput('image_url', '轮播图片', { helpText: '支持 jpg/png/webp ≤5MB' }),
      textInput('link_url', '跳转链接', { placeholder: 'https://… 或 /articles/…' }),
      numberInput('sort', '排序', { min: 0, defaultValue: 100 }),
      switchInput('is_active', '启用', { defaultValue: true })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('title', '标题'),
    textEntry('image_url', '图片'),
    textEntry('link_url', '链接'),
    textEntry('sort', '排序')
  ]
})
