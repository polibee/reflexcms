export default defineResource({
  name: 'topics',
  model: 'Topic',
  label: 'Topic',
  labelPlural: 'Topics',
  icon: 'message-square',
  group: 'Forum',
  sort: 20,
  searchable: ['title'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('title', 'Title', { sortable: true }),
    booleanColumn('is_pinned', 'Pinned'),
    booleanColumn('is_featured', 'Featured'),
    badgeColumn('status', 'Status', {
      open: { label: 'Open', variant: 'success' },
      closed: { label: 'Closed', variant: 'secondary' }
    }),
    numberColumn('forum_category_id', 'Board'),
    numberColumn('reply_count', 'Replies', { sortable: true }),
    numberColumn('best_reply_id', 'Best Reply'),
    dateColumn('last_reply_at', 'Last Reply', { sortable: true })
  ],

  rowActions: [
    defineAction({
      name: 'set-best-reply',
      label: '设最佳回复',
      icon: 'award',
      permission: 'topics.edit',
      visible: record => !record.best_reply_id,
      form: () => [
        section('最佳回复', [
          numberInput('reply_id', '回复 ID', { required: true, min: 1, helpText: '必须属于本帖，否则将被拒绝' })
        ])
      ],
      handler: async ({ record, values }) => {
        await $fetch(`/api/admin/topics/${record!.id}/best-reply`, {
          method: 'POST',
          body: { reply_id: values?.reply_id }
        })
        notify(`已将回复 #${values?.reply_id} 设为最佳`)
        emitAdminEvent('topics:refresh')
      }
    }),
    defineAction({
      name: 'clear-best-reply',
      label: '取消最佳回复',
      icon: 'x-circle',
      permission: 'topics.edit',
      visible: record => !!record.best_reply_id,
      handler: async ({ record }) => {
        await $fetch(`/api/admin/topics/${record!.id}/best-reply`, {
          method: 'POST',
          body: { reply_id: 0 }
        })
        notify('已取消最佳回复')
        emitAdminEvent('topics:refresh')
      }
    })
  ],

  form: () => [
    section('Topic', [
      textInput('title', 'Title', { required: true, colSpan: 2 }),
      relationInput('forum_category_id', 'Board', { resource: 'forums', labelKey: 'name' }, { required: true }),
      textarea('content', 'Content (Markdown)', { rows: 8, colSpan: 2 })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('title', 'Title'),
    textEntry('status', 'Status'),
    booleanEntry('is_pinned', 'Pinned'),
    booleanEntry('is_featured', 'Featured'),
    textEntry('forum_category_id', 'Board ID'),
    textEntry('best_reply_id', 'Best Reply ID'),
    textEntry('reply_count', 'Replies'),
    textEntry('like_count', 'Likes'),
    textEntry('view_count', 'Views'),
    datetimeEntry('last_reply_at', 'Last Reply'),
    dateEntry('created_at', 'Created')
  ]
})
