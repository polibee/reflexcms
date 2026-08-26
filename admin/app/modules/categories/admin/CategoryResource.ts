export default defineResource({
  name: 'categories',
  model: 'ArticleCategory',
  label: 'Category',
  labelPlural: 'Categories',
  icon: 'folder',
  group: 'Content',
  sort: 30,
  searchable: ['name', 'slug'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('name', 'Name', { sortable: true }),
    numberColumn('parent_id', 'Parent'),
    textColumn('slug', 'Slug'),
    numberColumn('sort', 'Sort'),
    dateColumn('created_at', 'Created', { sortable: true })
  ],

  form: () => [
    section('Category', [
      grid(2, [
        textInput('name', 'Name', { required: true }),
        relationInput('parent_id', 'Parent Category', { resource: 'categories', labelKey: 'name' })
      ]),
      grid(2, [
        numberInput('sort', 'Sort Order', { min: 0 }),
        textInput('description', 'Short Description')
      ])
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('name', 'Name'),
    textEntry('parent_id', 'Parent ID'),
    textEntry('slug', 'Slug'),
    textEntry('sort', 'Sort Order'),
    textEntry('description', 'Description')
  ]
})
