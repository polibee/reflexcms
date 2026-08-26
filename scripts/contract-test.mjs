#!/usr/bin/env node
/**
 * Contract regression for the ReflexCMS admin contract (plan §5).
 * Verifies that the Goravel backend implements exactly what nuxtadmin's
 * frontend expects: auth trio, generic resource gateway envelope, error
 * semantics and the permission matrix.
 *
 * Usage: node scripts/contract-test.mjs [backendBaseURL]
 * The backend must be running; default http://127.0.0.1:9000.
 */

const BASE = process.argv[2] || 'http://127.0.0.1:9000'

const ADMIN = { email: 'admin@reflexcms.dev', password: 'ReflexCMS@2026' }

let passed = 0
let failed = 0

function check(name, cond, extra = '') {
  if (cond) {
    passed++
    console.log(`  ok   ${name}`)
  } else {
    failed++
    console.error(`  FAIL ${name}${extra ? ` — ${extra}` : ''}`)
  }
}

async function api(method, path, { token, body } = {}) {
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: {
      'content-type': 'application/json',
      ...(token ? { authorization: `Bearer ${token}` } : {})
    },
    ...(body !== undefined ? { body: JSON.stringify(body) } : {})
  })
  let json = null
  try { json = await res.json() } catch {}
  return { status: res.status, json }
}

const stamp = Date.now()

async function main() {
  console.log(`contract-test against ${BASE}\n`)

  /* ---- auth ---- */
  // Use a throwaway email so the shared admin account never accumulates
  // lockout failures across runs.
  const badLogin = await api('POST', '/api/auth/login', {
    body: { email: `wrong+${stamp}@reflexcms.dev`, password: 'nope' }
  })
  check('bad login rejected with 401', badLogin.status === 401)

  const login = await api('POST', '/api/auth/login', { body: ADMIN })
  check('login returns token + AuthUser shape', login.status === 200
    && typeof login.json?.token === 'string'
    && login.json?.user && typeof login.json.user.id === 'number'
    && typeof login.json.user.name === 'string'
    && Array.isArray(login.json.user.permissions))
  const token = login.json.token

  const me = await api('GET', '/api/auth/me', { token })
  check('me echoes identity', me.status === 200 && me.json?.id === login.json.user.id)
  const perms = me.json.permissions
  check('super-admin holds wildcard permission', Array.isArray(perms) && perms.includes('*'))

  /* ---- gateway list envelope ---- */
  const users = await api('GET', '/api/admin/users?page=1&perPage=5&q=', { token })
  check('list is 200', users.status === 200)
  const env = users.json ?? {}
  check('Paginated envelope exact keys',
    ['items', 'total', 'page', 'perPage', 'totalPages'].every(k => k in env),
    JSON.stringify(Object.keys(env)))

  /* ---- create / show / update / bulk-delete roundtrip ---- */
  const email = `contract+${stamp}@reflexcms.dev`
  const created = await api('POST', '/api/admin/users', {
    token,
    body: { username: `contract${stamp}`, email, password: 'Contract@123', status: 'active' }
  })
  check('create returns row without password leak', created.status === 200
    && created.json?.id > 0 && !('password' in (created.json ?? {})),
    JSON.stringify(created.json))
  const id = created.json?.id

  const shown = await api('GET', `/api/admin/users/${id}`, { token })
  check('show finds created row', shown.status === 200 && shown.json?.email === email)

  const updated = await api('PUT', `/api/admin/users/${id}`, {
    token, body: { bio: 'touched by contract test' }
  })
  check('update applies partial payload', updated.status === 200 && updated.json?.bio === 'touched by contract test')

  const noPerms = await api('POST', '/api/admin/roles', {
    body: { name: 'x', display_name: 'X', permissions: '' }
  })
  check('anonymous write gets 401', noPerms.status === 401)

  const removed = await api('POST', '/api/admin/users/bulk-delete', { token, body: { ids: [id] } })
  check('bulk-delete reports removed count', removed.status === 200 && removed.json?.removed === 1)

  const gone = await api('GET', `/api/admin/users/${id}`, { token })
  check('deleted row is 404', gone.status === 404)

  /* ---- stats ---- */
  const stats = await api('GET', '/api/admin/stats', { token })
  check('stats reachable with expected keys', stats.status === 200
    && ['usersTotal', 'articlesTotal', 'topicsTotal', 'pendingComments'].every(k => k in (stats.json ?? {})))

  /* ---- articles lifecycle (M2) ---- */
  const ssrf = await api('POST', '/api/admin/articles', {
    token,
    body: { title: `SSRF Probe ${stamp}`, cover: 'http://169.254.169.254/latest/meta-data' }
  })
  check('private-URL cover rejected (SSRF guard)', ssrf.status === 422,
    `got ${ssrf.status}`)

  const cat = await api('POST', '/api/admin/categories', {
    token, body: { name: `Tech ${stamp}` }
  })
  check('category created with generated slug', cat.status === 200
    && typeof cat.json?.slug === 'string' && cat.json.slug.length > 0)

  const art = await api('POST', '/api/admin/articles', {
    token,
    body: {
      title: `Contract Article ${stamp}`,
      content: 'body mentioning contractsearchterm',
      summary: 'summary for fulltext',
      tags: 'alpha\nbeta',
      category_id: cat.json?.id
    }
  })
  check('article draft created with slug + tag echo', art.status === 200
    && art.json?.status === 'draft'
    && typeof art.json?.slug === 'string' && art.json.slug.length > 0
    && art.json?.tags === 'alpha\nbeta')
  const aid = art.json?.id

  const pub = await api('POST', `/api/admin/articles/${aid}/publish`, { token })
  check('publish action sets status and published_at', pub.status === 200
    && pub.json?.status === 'published'
    && Boolean(pub.json?.published_at))

  const searchHit = await api('GET', `/api/admin/articles?q=${stamp}`, { token })
  check('full-text search hits the published article', searchHit.status === 200
    && Array.isArray(searchHit.json?.items)
    && searchHit.json.items.some(i => i.id === aid))

  const cjkHit = await api('GET', '/api/admin/articles?q=契约文章', { token })
  check('CJK query falls back to ILIKE without error', cjkHit.status === 200)

  const arch = await api('POST', `/api/admin/articles/${aid}/archive`, { token })
  check('archive transitions published→archived', arch.status === 200
    && arch.json?.status === 'archived')

  const republish = await api('POST', `/api/admin/articles/${aid}/publish`, { token })
  check('archived→publish rejected as 422', republish.status === 422)

  const cleanupArt = await api('POST', '/api/admin/articles/bulk-delete', { token, body: { ids: [aid] } })
  const cleanupCat = await api('POST', '/api/admin/categories/bulk-delete', { token, body: { ids: [cat.json?.id] } })
  check('lifecycle fixtures cleaned up', cleanupArt.status === 200 && cleanupCat.status === 200)

  /* ---- forum lifecycle (M3) ---- */
  const board = await api('POST', '/api/admin/forums', { token, body: { name: `Board ${stamp}` } })
  check('forum board created with slug', board.status === 200 && board.json?.slug)

  const topic = await api('POST', '/api/admin/topics', {
    token,
    body: { title: `Topic ${stamp}`, content: 'opening post', forum_category_id: board.json?.id }
  })
  check('topic created as open', topic.status === 200 && topic.json?.status === 'open')
  const tid = topic.json?.id

  // Concurrency probe: parallel replies must all land, with unique floors and
  // an exact reply_count (DoD of M3).
  const REPLY_N = 10
  const replies = await Promise.all(Array.from({ length: REPLY_N }, (_, i) =>
    api('POST', '/api/admin/forum/reply', { token, body: { topic_id: tid, content: `reply #${i}` } })
  ))
  check('all concurrent replies accepted', replies.every(r => r.status === 200),
    JSON.stringify(replies.filter(r => r.status !== 200).map(r => r.json)))
  const floors = new Set(replies.map(r => r.json?.floor))
  check('floors unique under concurrency', floors.size === REPLY_N,
    JSON.stringify([...floors]))

  const shownTopic = await api('GET', `/api/admin/topics/${tid}`, { token })
  check('reply_count matches reply rows exactly', shownTopic.json?.reply_count === REPLY_N,
    `reply_count=${shownTopic.json?.reply_count}`)

  const likeOn = await api('POST', '/api/admin/forum/like', {
    token, body: { likeable_type: 'topics', likeable_id: tid }
  })
  check('first toggle likes the topic', likeOn.status === 200 && likeOn.json?.liked === true
    && likeOn.json?.like_count === 1)
  const likeOff = await api('POST', '/api/admin/forum/like', {
    token, body: { likeable_type: 'topics', likeable_id: tid }
  })
  check('second toggle unlikes the topic', likeOff.status === 200 && likeOff.json?.liked === false
    && likeOff.json?.like_count === 0)

  const pin = await api('POST', `/api/admin/topics/${tid}/pin`, { token })
  check('pin action sets is_pinned', pin.status === 200 && pin.json?.is_pinned === true)

  const hot = await api('GET', '/api/admin/topics?sortBy=hot&sortDir=desc', { token })
  check('hot sort alias accepted', hot.status === 200 && Array.isArray(hot.json?.items))

  const close = await api('POST', `/api/admin/topics/${tid}/close`, { token })
  check('close action sets status closed', close.status === 200 && close.json?.status === 'closed')

  const replyClosed = await api('POST', '/api/admin/forum/reply', {
    token, body: { topic_id: tid, content: 'should be rejected' }
  })
  check('closed topic rejects replies', replyClosed.status === 422)

  const cleanupReplies = await api('POST', `/api/admin/topics/bulk-delete`, { token, body: { ids: [tid] } })
  const cleanupBoard = await api('POST', `/api/admin/forums/bulk-delete`, { token, body: { ids: [board.json?.id] } })
  check('forum fixtures cleaned up', cleanupReplies.status === 200 && cleanupBoard.status === 200)

  /* ---- M4: settings + site.mode + comments moderation + notifications ---- */
  const cfg1 = await api('GET', '/api/v1/site/config')
  check('site config is anonymous-readable with mode', cfg1.status === 200
    && typeof cfg1.json?.mode === 'string' && Array.isArray(cfg1.json?.sections))

  // flip mode via the settings resource, expect config to follow
  const modeRow = await api('GET', '/api/admin/settings?q=site.mode', { token })
  const modeId = modeRow.json?.items?.[0]?.id
  check('site.mode row exists in settings resource', !!modeId)
  if (modeId) {
    const setForum = await api('PUT', `/api/admin/settings/${modeId}`, {
      token, body: { value: JSON.stringify('forum') }
    })
    const cfg2 = await api('GET', '/api/v1/site/config')
    check('mode switch reflects in site config (cache invalidated)',
      setForum.status === 200 && cfg2.json?.mode === 'forum'
      && !cfg2.json.sections.includes('articles'))
    await api('PUT', `/api/admin/settings/${modeId}`, {
      token, body: { value: JSON.stringify('hybrid') }
    })
  }

  /* comments: pending default → approve/reject + article count + notify */
  const cArticle = await api('POST', '/api/admin/articles', {
    token, body: { title: `Commented ${stamp}`, content: 'body' }
  })
  const caid = cArticle.json?.id

  const cm1 = await api('POST', '/api/admin/comments', {
    token, body: { article_id: caid, user_id: 1, content: 'nice article' }
  })
  check('comment defaults to pending', cm1.status === 200 && cm1.json?.status === 'pending')
  const cm2 = await api('POST', '/api/admin/comments', {
    token, body: { article_id: caid, user_id: 1, content: 'second opinion' }
  })

  const appr = await api('POST', `/api/admin/comments/${cm1.json?.id}/approve`, { token })
  check('approve transitions to approved', appr.status === 200 && appr.json?.status === 'approved')

  const artAfter = await api('GET', `/api/admin/articles/${caid}`, { token })
  check('article comment_count counts approved only', artAfter.json?.comment_count === 1,
    `count=${artAfter.json?.comment_count}`)

  const rej = await api('POST', `/api/admin/comments/${cm2.json?.id}/reject`, { token })
  check('reject transitions to rejected', rej.status === 200 && rej.json?.status === 'rejected')

  const notif = await api('GET', '/api/admin/notifications?perPage=50', { token })
  const types = (notif.json?.items ?? []).map(n => n.type)
  check('moderation produced a notification', types.includes('comment.rejected'),
    JSON.stringify(types))

  const readAll = await api('POST', '/api/admin/notifications/read-all', { token })
  check('read-all returns updated count', readAll.status === 200 && readAll.json?.updated >= 1)

  /* cleanup M4 fixtures */
  await api('POST', `/api/admin/comments/bulk-delete`, { token, body: { ids: [cm1.json?.id, cm2.json?.id] } })
  await api('POST', `/api/admin/articles/bulk-delete`, { token, body: { ids: [caid] } })

  /* ---- UX round: permission catalog / scheduled / best reply / recycle bin ---- */
  const catalogRes = await api('GET', '/api/admin/permissions', { token })
  check('permission catalog served grouped', catalogRes.status === 200
    && Array.isArray(catalogRes.json)
    && catalogRes.json.some(g => g.permissions?.some(p => p.name === 'articles.publish')))

  const badRole = await api('POST', '/api/admin/roles', {
    token, body: { name: `bad${stamp}`, display_name: 'Bad', permissions: 'articles.vew' }
  })
  check('unknown permission rejected with 422', badRole.status === 422,
    `got ${badRole.status}`)

  // scheduled release: future "at" parks the draft as scheduled.
  // Built from LOCAL time components to match the server's clock.
  const sArt = await api('POST', '/api/admin/articles', {
    token, body: { title: `Scheduled ${stamp}`, content: 'later' }
  })
  const said = sArt.json?.id
  const due = new Date(Date.now() + 3600 * 1000)
  const pad = n => String(n).padStart(2, '0')
  const futureAt = `${due.getFullYear()}-${pad(due.getMonth() + 1)}-${pad(due.getDate())} ` +
    `${pad(due.getHours())}:${pad(due.getMinutes())}`
  const sched = await api('POST', `/api/admin/articles/${said}/schedule`, {
    token, body: { at: futureAt }
  })
  check('schedule action parks article as scheduled', sched.status === 200
    && sched.json?.status === 'scheduled'
    && Boolean(sched.json?.published_at))

  // Manual publish overrides the schedule (WordPress semantics): scheduled
  // → published is a legal transition when an editor forces it early.
  const earlyPub = await api('POST', `/api/admin/articles/${said}/publish`, { token })
  check('manual publish overrides schedule', earlyPub.status === 200
    && earlyPub.json?.status === 'published')

  // recycle bin: soft-deleted rows appear under trashed=only and can be restored
  const delRes = await api('POST', `/api/admin/articles/bulk-delete`, { token, body: { ids: [said] } })
  check('bulk-delete soft-deletes the article', delRes.status === 200 && delRes.json?.removed === 1)
  const binList = await api('GET', '/api/admin/articles?trashed=only&perPage=100', { token })
  const binItems = binList.json?.items ?? []
  check('trashed=only lists the deleted article', binList.status === 200
    && binItems.some(i => i.id === said))

  const restored2 = await api('POST', `/api/admin/articles/${said}/restore`, { token })
  check('restore revives the soft-deleted article', restored2.status === 200
    && restored2.json?.id === said)

  // forum best reply
  const fBoard = await api('POST', '/api/admin/forums', { token, body: { name: `BestBoard ${stamp}` } })
  const fTopic = await api('POST', '/api/admin/topics', {
    token, body: { title: `BestTopic ${stamp}`, content: 'q?', forum_category_id: fBoard.json?.id }
  })
  const fReply = await api('POST', '/api/admin/forum/reply', {
    token, body: { topic_id: fTopic.json?.id, content: 'the answer' }
  })
  const setBest = await api('POST', `/api/admin/topics/${fTopic.json?.id}/best-reply`, {
    token, body: { reply_id: fReply.json?.id }
  })
  check('best reply marker set', setBest.status === 200
    && setBest.json?.best_reply_id === fReply.json?.id)

  const wrongBest = await api('POST', `/api/admin/topics/999999/best-reply`, {
    token, body: { reply_id: fReply.json?.id }
  })
  check('cross-topic best reply rejected', wrongBest.status !== 200)

  const clearBest = await api('POST', `/api/admin/topics/${fTopic.json?.id}/best-reply`, {
    token, body: { reply_id: 0 }
  })
  check('clearing best reply works', clearBest.status === 200 && clearBest.json?.best_reply_id === 0)

  await api('POST', `/api/admin/topics/bulk-delete`, { token, body: { ids: [fTopic.json?.id] } })
  await api('POST', `/api/admin/forums/bulk-delete`, { token, body: { ids: [fBoard.json?.id] } })

  /* ---- M5: audit trail + account lockout ---- */
  const audit = await api('GET', '/api/admin/operation_logs?perPage=50', { token })
  const auditActions = (audit.json?.items ?? []).map(l => l.action)
  check('audit trail recorded gateway mutations', audit.status === 200
    && auditActions.includes('create') && auditActions.includes('bulk-delete'),
    JSON.stringify(auditActions.slice(0, 8)))
  const auditSensitive = (audit.json?.items ?? [])
    .some(l => typeof l.payload === 'string' && l.payload.includes('Contract@123'))
  check('audit payloads do not leak passwords', !auditSensitive)

  // account lockout: hammer bad logins on a fresh email until locked (429).
  // The per-IP limiter may trip first — both produce 429, which is the
  // assertion target.
  let saw429 = false
  for (let i = 0; i < 12; i++) {
    const r = await api('POST', '/api/auth/login', {
      body: { email: `lockprobe+${stamp}@reflexcms.dev`, password: 'wrong' }
    })
    if (r.status === 429) { saw429 = true; break }
  }
  check('brute force eventually throttled to 429', saw429)

  console.log(`\n${passed} passed, ${failed} failed`)
  process.exit(failed > 0 ? 1 : 0)
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
