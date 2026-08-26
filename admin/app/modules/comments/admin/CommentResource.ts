const statusBadges = {
  pending: { label: 'Pending', variant: 'warning' },
  approved: { label: 'Approved', variant: 'success' },
  rejected: { label: 'Rejected', variant: 'secondary' },
  spam: { label: 'Spam', variant: 'default' }
} as const

export default defineResource({
  name: 'comments',
  model: 'Comment',
  label: 'Comment',
  labelPlural: 'Comments',
  icon: 'message-circle',
  group: 'Content',
  sort: 50,
  searchable: ['content'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    numberColumn('article_id', 'Article'),
    textColumn('content', 'Content'),
    badgeColumn('status', 'Status', statusBadges),
    dateColumn('created_at', 'Created', { sortable: true }),
    actionsColumn([
      defineAction({
        name: 'approve',
        label: 'Approve',
        icon: 'badge-check',
        permission: 'comments.moderate',
        visible: record => record.status === 'pending' || record.status === 'rejected',
        handler: async ({ record }) => {
          await $fetch(`/api/admin/comments/${record!.id}/approve`, { method: 'POST' })
          notify('Comment approved')
          emitAdminEvent('comments:refresh')
        }
      }),
      defineAction({
        name: 'reject',
        label: 'Reject',
        icon: 'x',
        permission: 'comments.moderate',
        visible: record => record.status === 'pending' || record.status === 'approved',
        confirm: { title: 'Reject this comment?', confirmLabel: 'Reject' },
        handler: async ({ record }) => {
          await $fetch(`/api/admin/comments/${record!.id}/reject`, { method: 'POST' })
          notify('Comment rejected')
          emitAdminEvent('comments:refresh')
        }
      })
    ])
  ],

  bulkActions: [
    defineAction({
      name: 'bulk-approve',
      label: 'Approve Selected',
      icon: 'badge-check',
      permission: 'comments.moderate',
      confirm: { title: 'Approve all selected comments?', confirmLabel: 'Approve' },
      handler: async ({ ids }) => {
        for (const id of ids ?? []) {
          await $fetch(`/api/admin/comments/${id}/approve`, { method: 'POST' })
        }
        notify(`${ids!.length} comments approved`)
        emitAdminEvent('comments:refresh')
      }
    })
  ],

  form: () => [
    section('Comment (administrative edit)', [
      grid(2, [
        numberInput('article_id', 'Article ID', { required: true, min: 1 }),
        numberInput('parent_id', 'Parent Comment ID', { min: 0 })
      ]),
      textarea('content', 'Content', { rows: 5, colSpan: 2 })
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('article_id', 'Article ID'),
    textEntry('user_id', 'Author ID'),
    textEntry('content', 'Content'),
    badgeEntry('status', 'Status', statusBadges),
    dateEntry('created_at', 'Created')
  ]
})
