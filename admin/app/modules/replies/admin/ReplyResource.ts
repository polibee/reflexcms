export default defineResource({
  name: 'replies',
  model: 'Reply',
  label: 'Reply',
  labelPlural: 'Replies',
  icon: 'reply',
  group: 'Forum',
  sort: 30,
  searchable: ['content'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    numberColumn('topic_id', 'Topic'),
    numberColumn('floor', 'Floor'),
    numberColumn('like_count', 'Likes'),
    dateColumn('created_at', 'Created', { sortable: true })
  ],

  form: () => [
    section('Reply (administrative edit)', [
      grid(2, [
        numberInput('topic_id', 'Topic ID', { required: true, min: 1 }),
        numberInput('parent_id', 'Parent Reply ID', { min: 0 })
      ]),
      textarea('content', 'Content', { rows: 5, colSpan: 2 })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('topic_id', 'Topic ID'),
    textEntry('floor', 'Floor'),
    textEntry('user_id', 'Author ID'),
    textEntry('content', 'Content'),
    textEntry('like_count', 'Likes'),
    dateEntry('created_at', 'Created')
  ]
})
