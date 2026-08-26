/* Sidebar widgets. The config column is JSON in storage, but admins never
 * see raw JSON: each widget type renders its own visual fields (cfg_*),
 * and transformIn/transformOut pack them into the config payload. Preset
 * types get sensible defaults so picking a card type is enough. */

type Values = Record<string, unknown>

const PRESET_DEFAULTS: Record<string, Record<string, unknown>> = {
  categories: { show_count: true, count: 10 },
  tags: { max_tags: 20 },
  latest_articles: { count: 5 },
  hottest_articles: { count: 5 },
  latest_topics: { count: 5 },
  hottest_topics: { count: 5 },
  boards: {},
  user_stats: {},
  author_center: {},
  composer: {}
}

function num(v: unknown, fallback: number): number {
  const n = Number(v)
  return Number.isFinite(n) && n > 0 ? n : fallback
}

export default defineResource({
  name: 'widgets',
  model: 'Widget',
  label: '侧栏卡片',
  labelPlural: '侧栏卡片',
  icon: 'layout-dashboard',
  group: '布局管理',
  sort: 10,
  searchable: ['title', 'widget_type'],

  table: () => [
    textColumn('id', 'ID', { sortable: true }),
    textColumn('area', '显示区域'),
    textColumn('widget_type', '类型'),
    textColumn('title', '标题'),
    numberColumn('sort', '排序', { sortable: true }),
    booleanColumn('is_active', '启用'),
    dateColumn('created_at', 'Created')
  ],

  form: () => [
    section('Widget', [
      grid(2, [
        selectInput('area', '显示区域', [
          { label: '博客侧边栏（blog_sidebar）', value: 'blog_sidebar' },
          { label: '论坛侧边栏（forum_sidebar）', value: 'forum_sidebar' }
        ], { required: true }),
        selectInput('widget_type', '卡片类型', [
          { label: '用户中心', value: 'author_center' },
          { label: '发帖卡片', value: 'composer' },
          { label: '分类卡片', value: 'categories' },
          { label: '标签卡片', value: 'tags' },
          { label: '最新文章', value: 'latest_articles' },
          { label: '最热文章', value: 'hottest_articles' },
          { label: '最新帖子', value: 'latest_topics' },
          { label: '最热帖子', value: 'hottest_topics' },
          { label: '板块卡片', value: 'boards' },
          { label: '用户统计', value: 'user_stats' },
          { label: '自定义 HTML', value: 'custom_html' },
          { label: '自定义图片链接', value: 'custom_image' },
          { label: '自定义文本链接', value: 'custom_text_link' }
        ], {
          required: true,
          helpText: '预设卡片无需填写配置，选择类型即可使用；数量选项可微调'
        })
      ]),
      textInput('title', '卡片标题', { colSpan: 2 }),

      /* ---- per-type visual config (no raw JSON anywhere) ---- */
      grid(2, [
        switchInput('cfg_show_count', '显示数量统计', { visibleIf: v => v.widget_type === 'categories' }),
        numberInput('cfg_count_cat', '显示条数', {
          min: 1,
          max: 50,
          placeholder: '默认 10',
          helpText: '留空使用默认值',
          visibleIf: v => v.widget_type === 'categories'
        }),
        numberInput('cfg_max_tags', '标签数量上限', {
          min: 1,
          max: 100,
          placeholder: '默认 20',
          helpText: '留空使用默认值',
          visibleIf: v => v.widget_type === 'tags'
        }),
        numberInput('cfg_count_articles', '文章数量', {
          min: 1,
          max: 20,
          placeholder: '默认 5',
          helpText: '留空使用默认值',
          visibleIf: v => v.widget_type === 'latest_articles' || v.widget_type === 'hottest_articles'
        }),
        numberInput('cfg_count_topics', '帖子数量', {
          min: 1,
          max: 20,
          placeholder: '默认 5',
          helpText: '留空使用默认值',
          visibleIf: v => v.widget_type === 'latest_topics' || v.widget_type === 'hottest_topics'
        })
      ]),
      textarea('cfg_html', 'HTML / JS 代码', {
        rows: 6,
        colSpan: 2,
        placeholder: '<div style="...">直接粘贴 HTML 或 <script> 代码</div>',
        helpText: '无需写成 JSON，保存时自动包装',
        visibleIf: v => v.widget_type === 'custom_html'
      }),
      grid(2, [
        textInput('cfg_image_url', '图片地址', {
          placeholder: 'https://…',
          visibleIf: v => v.widget_type === 'custom_image'
        }),
        textInput('cfg_image_link', '点击跳转链接', {
          placeholder: '/articles 或 https://…',
          visibleIf: v => v.widget_type === 'custom_image'
        })
      ]),
      grid(2, [
        textInput('cfg_link_text', '链接文字', {
          visibleIf: v => v.widget_type === 'custom_text_link'
        }),
        textInput('cfg_link_url', '链接地址', {
          placeholder: '/forums 或 https://…',
          visibleIf: v => v.widget_type === 'custom_text_link'
        })
      ]),

      numberInput('sort', '排序', { min: 0, helpText: '已自动接上当前最大序号，可自行修改' }),
      switchInput('is_active', '启用')
    ])
  ],

  infolist: () => [
    textEntry('id', 'ID'),
    textEntry('area', '区域'),
    textEntry('widget_type', '类型'),
    textEntry('title', '标题'),
    textEntry('config', '配置'),
    textEntry('sort', '排序')
  ],

  /* CREATE form defaults: sort continues from the current highest value */
  defaultValues: async () => {
    try {
      const res = await $fetch<{ items: Array<Record<string, unknown>> }>('/api/admin/widgets', { query: { perPage: 200 } })
      const max = (res.items ?? []).reduce((m, w) => Math.max(m, Number(w.sort ?? 0)), 0)
      return { sort: max + 1 }
    } catch {
      return { sort: 100 }
    }
  },

  /* API record → form values */
  transformIn: (record) => {
    const r = (record.widget ?? record) as Values
    const cfg: Values = {}
    if (typeof r.config === 'string' && r.config.trim()) {
      try { Object.assign(cfg, JSON.parse(r.config) as Values) } catch { /* raw html handled below */ }
    }
    return {
      ...r,
      cfg_show_count: cfg.show_count ?? true,
      cfg_count_cat: cfg.count ?? '',
      cfg_max_tags: cfg.max_tags ?? '',
      cfg_count_articles: cfg.count ?? '',
      cfg_count_topics: cfg.count ?? '',
      // custom_html unwraps to the raw code for editing
      cfg_html: typeof cfg.html === 'string' ? cfg.html : (typeof r.config === 'string' && !r.config.trim().startsWith('{') ? r.config : ''),
      cfg_image_url: cfg.image_url ?? '',
      cfg_image_link: cfg.link_url ?? '',
      cfg_link_text: cfg.text ?? '',
      cfg_link_url: cfg.url ?? ''
    }
  },

  /* form values → API payload */
  transformOut: (values) => {
    const t = String(values.widget_type ?? '')
    let config = ''

    if (t === 'custom_html') {
      const html = String(values.cfg_html ?? '').trim()
      config = html.startsWith('{') ? html : JSON.stringify({ html })
    } else if (t === 'custom_image') {
      config = JSON.stringify({ image_url: String(values.cfg_image_url ?? ''), link_url: String(values.cfg_image_link ?? '') })
    } else if (t === 'custom_text_link') {
      config = JSON.stringify({ text: String(values.cfg_link_text ?? ''), url: String(values.cfg_link_url ?? '') })
    } else if (t === 'categories') {
      config = JSON.stringify({ show_count: values.cfg_show_count !== false, count: num(values.cfg_count_cat, PRESET_DEFAULTS.categories.count as number) })
    } else if (t === 'tags') {
      config = JSON.stringify({ max_tags: num(values.cfg_max_tags, PRESET_DEFAULTS.tags.max_tags as number) })
    } else if (t === 'latest_articles' || t === 'hottest_articles') {
      config = JSON.stringify({ count: num(values.cfg_count_articles, 5) })
    } else if (t === 'latest_topics' || t === 'hottest_topics') {
      config = JSON.stringify({ count: num(values.cfg_count_topics, 5) })
    }

    const payload: Values = { ...values }
    for (const k of Object.keys(payload)) {
      if (k.startsWith('cfg_')) delete payload[k]
    }
    payload.config = config
    return payload
  }
})
