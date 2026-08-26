export default defineResource({
  name: 'forums',
  model: 'ForumCategory',
  label: 'Board',
  labelPlural: 'Boards',
  icon: 'messages-square',
  group: 'Forum',
  sort: 10,
  searchable: ['name', 'slug'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('name', 'Name', { sortable: true }),
    textColumn('slug', 'Slug'),
    numberColumn('topic_count', 'Topics'),
    numberColumn('sort', 'Sort'),
    dateColumn('created_at', 'Created', { sortable: true })
  ],

  form: () => [
    section('Board', [
      grid(2, [
        textInput('name', 'Name', { required: true }),
        numberInput('sort', 'Sort Order', { min: 0 })
      ]),
      grid(2, [
        textInput('slug', 'Slug', { placeholder: 'Leave empty to generate from the name' }),
        textInput('icon', 'Icon Key')
      ]),
      textarea('description', 'Description', { rows: 3, colSpan: 2 })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('name', 'Name'),
    textEntry('slug', 'Slug'),
    textEntry('icon', 'Icon'),
    textEntry('sort', 'Sort Order'),
    textEntry('description', 'Description')
  ]
})
