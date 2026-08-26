const statusBadges = {
  draft: { label: 'Draft', variant: 'warning' },
  scheduled: { label: 'Scheduled', variant: 'default' },
  published: { label: 'Published', variant: 'success' },
  archived: { label: 'Archived', variant: 'secondary' }
} as const

export default defineResource({
  name: 'articles',
  model: 'Article',
  label: 'Article',
  labelPlural: 'Articles',
  icon: 'file-text',
  group: 'Content',
  sort: 20,
  searchable: ['title', 'summary'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('title', 'Title', { sortable: true }),
    booleanColumn('is_pinned', 'Pinned'),
    badgeColumn('status', 'Status', statusBadges),
    numberColumn('category_id', 'Category'),
    numberColumn('view_count', 'Views'),
    dateColumn('published_at', 'Published', { sortable: true }),
    actionsColumn([
      defineAction({
        name: 'publish',
        label: '发布',
        icon: 'badge-check',
        permission: 'articles.publish',
        visible: record => record.status === 'draft' || record.status === 'scheduled',
        confirm: { title: '立即发布这篇文章？', confirmLabel: '发布' },
        handler: async ({ record }) => {
          await $fetch(`/api/admin/articles/${record!.id}/publish`, { method: 'POST' })
          notify('文章已发布')
          emitAdminEvent('articles:refresh')
        }
      }),
      defineAction({
        name: 'schedule',
        label: '定时发布',
        icon: 'calendar-clock',
        permission: 'articles.publish',
        visible: record => record.status === 'draft',
        form: () => [
          section('定时发布', [
            textInput('at', '发布时间', {
              required: true,
              placeholder: '2026-09-01 09:00',
              helpText: '格式：YYYY-MM-DD HH:mm；到点后由调度器自动发布'
            })
          ])
        ],
        handler: async ({ record, values }) => {
          await $fetch(`/api/admin/articles/${record!.id}/schedule`, {
            method: 'POST',
            body: { at: values?.at }
          })
          notify(`将于 ${values?.at} 自动发布`)
          emitAdminEvent('articles:refresh')
        }
      }),
      defineAction({
        name: 'pin',
        label: '置顶',
        icon: 'pin',
        permission: 'articles.edit',
        visible: record => !record.is_pinned,
        handler: async ({ record }) => {
          await $fetch(`/api/admin/articles/${record!.id}/pin`, { method: 'POST' })
          notify('已置顶')
          emitAdminEvent('articles:refresh')
        }
      }),
      defineAction({
        name: 'unpin',
        label: '取消置顶',
        icon: 'pin-off',
        permission: 'articles.edit',
        visible: record => !!record.is_pinned,
        handler: async ({ record }) => {
          await $fetch(`/api/admin/articles/${record!.id}/unpin`, { method: 'POST' })
          notify('已取消置顶')
          emitAdminEvent('articles:refresh')
        }
      }),
      defineAction({
        name: 'restore-bin',
        label: '从回收站恢复',
        icon: 'undo-2',
        permission: 'articles.delete',
        visible: record => record.status === 'archived' && !!record.deleted_at,
        handler: async ({ record }) => {
          await $fetch(`/api/admin/articles/${record!.id}/restore`, { method: 'POST' })
          notify('已恢复为草稿')
          emitAdminEvent('articles:refresh')
        }
      }),
      defineAction({
        name: 'unpublish',
        label: '转为草稿',
        icon: 'undo-2',
        permission: 'articles.edit',
        visible: record => record.status === 'published',
        confirm: { title: '将这篇文章转回草稿？', confirmLabel: '转草稿' },
        handler: async ({ record }) => {
          await $fetch(`/api/admin/articles/${record!.id}/unpublish`, { method: 'POST' })
          notify('已转回草稿')
          emitAdminEvent('articles:refresh')
        }
      })
    ])
  ],

  bulkActions: [
    defineAction({
      name: 'bulk-publish',
      label: 'Publish Selected',
      icon: 'badge-check',
      permission: 'articles.publish',
      confirm: { title: 'Publish all selected articles?', confirmLabel: 'Publish' },
      handler: async ({ ids }) => {
        for (const id of ids ?? []) {
          await $fetch(`/api/admin/articles/${id}/publish`, { method: 'POST' })
        }
        notify(`${ids!.length} articles published`)
        emitAdminEvent('articles:refresh')
      }
    })
  ],

  form: () => [
    section('Content', [
      textInput('title', 'Title', {
        required: true,
        placeholder: 'Getting started with Goravel',
        colSpan: 2
      }),
      textInput('slug', 'Slug', {
        placeholder: 'Leave empty to generate from the title',
        colSpan: 2
      }),
      richtext('content', 'Content (Markdown)', { colSpan: 2 }),
      textarea('summary', 'Summary', { rows: 3, colSpan: 2 }),
      textarea('tags', 'Tags', {
        rows: 2,
        colSpan: 2,
        placeholder: 'One tag per line, e.g.\ngoravel\nnuxt'
      })
    ]),
    section('Metadata', [
      grid(2, [
        relationInput('category_id', 'Category', { resource: 'categories', labelKey: 'name' }),
        uploadInput('cover', '封面图片', { helpText: '支持 jpg/png/webp ≤5MB' })
      ])
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('title', 'Title'),
    badgeEntry('status', 'Status', statusBadges),
    booleanEntry('is_pinned', 'Pinned'),
    textEntry('slug', 'Slug'),
    textEntry('summary', 'Summary'),
    textEntry('category_id', 'Category ID'),
    datetimeEntry('published_at', 'Published'),
    dateEntry('created_at', 'Created')
  ]
})
