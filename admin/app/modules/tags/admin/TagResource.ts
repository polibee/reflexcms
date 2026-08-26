export default defineResource({
  name: 'tags',
  model: 'Tag',
  label: 'Tag',
  labelPlural: 'Tags',
  icon: 'tag',
  group: 'Content',
  sort: 40,
  searchable: ['name', 'slug'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('name', 'Name', { sortable: true }),
    textColumn('slug', 'Slug'),
    dateColumn('created_at', 'Created', { sortable: true })
  ],

  form: () => [
    section('Tag', [
      grid(2, [
        textInput('name', 'Name', { required: true }),
        textInput('slug', 'Slug', { placeholder: 'Leave empty to generate from the name' })
      ])
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('name', 'Name'),
    textEntry('slug', 'Slug'),
    dateEntry('created_at', 'Created')
  ]
})
