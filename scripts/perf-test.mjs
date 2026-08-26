#!/usr/bin/env node
/**
 * Performance baseline probe (plan §M5): measures P95 latency of the core
 * admin list endpoints against a seeded database (scripts/perf-seed.sql:
 * 10k articles / 100k replies). Results are archived in docs/性能基线.md.
 *
 * Usage: node scripts/perf-test.mjs [backendBaseURL]
 */

const BASE = process.argv[2] || 'http://127.0.0.1:9000'
const N = Number(process.env.PERF_N || 40)

const ADMIN = { email: 'admin@reflexcms.dev', password: 'ReflexCMS@2026' }

async function login() {
  const res = await fetch(`${BASE}/api/auth/login`, {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(ADMIN)
  })
  if (!res.ok) throw new Error(`login failed: ${res.status}`)
  return (await res.json()).token
}

function p95(durations) {
  const sorted = [...durations].sort((a, b) => a - b)
  const idx = Math.max(0, Math.ceil(sorted.length * 0.95) - 1)
  return sorted[idx]
}

async function measure(label, path, token, n = N) {
  // warmup
  for (let i = 0; i < 3; i++) await fetch(`${BASE}${path}`, { headers: { authorization: `Bearer ${token}` } })

  const durations = []
  let ok = 0
  for (let i = 0; i < n; i++) {
    const start = performance.now()
    const res = await fetch(`${BASE}${path}`, { headers: { authorization: `Bearer ${token}` } })
    durations.push(performance.now() - start)
    if (res.ok) ok++
  }
  durations.sort((a, b) => a - b)
  const result = {
    label,
    p50: Math.round(durations[Math.floor(durations.length / 2)] * 10) / 10,
    p95: Math.round(p95(durations) * 10) / 10,
    max: Math.round(durations[durations.length - 1] * 10) / 10,
    okRatio: `${ok}/${n}`
  }
  console.log(
    `${result.label.padEnd(34)} p50=${String(result.p50).padStart(7)}ms  ` +
    `p95=${String(result.p95).padStart(7)}ms  max=${String(result.max).padStart(7)}ms  (${result.okRatio} ok)`
  )
  return result
}

const token = await login()
console.log(`perf-test against ${BASE} (${N} samples each)\n`)

const results = []
results.push(await measure('articles list (10k rows)', '/api/admin/articles?perPage=20&sortBy=published_at&sortDir=desc', token))
results.push(await measure('articles FTS search', '/api/admin/articles?q=lorem&perPage=20', token))
results.push(await measure('articles CJK ILIKE fallback', `/api/admin/articles?q=${encodeURIComponent('入门指南')}&perPage=20`, token))
results.push(await measure('topics list hot sort', '/api/admin/topics?perPage=20&sortBy=hot&sortDir=desc', token))
results.push(await measure('replies list (100k rows)', '/api/admin/replies?perPage=20', token))
results.push(await measure('stats aggregate', '/api/admin/stats', token))

const worst = Math.max(...results.map(r => r.p95))
console.log(`\nworst p95 = ${worst}ms (target: < 200ms)`)
process.exit(worst < 200 ? 0 : 1)
