-- Performance seed (plan §M5): 10,000 articles + 1,000 topics + 100,000 replies.
-- Idempotent-ish: prefixed with 'perf-' so it can be wiped via
--   DELETE FROM articles WHERE slug LIKE 'perf-%';  (etc.)
-- Run:  psql -U reflexcms -h 127.0.0.1 -p 5433 -d reflexcms -f scripts/perf-seed.sql

INSERT INTO articles (author_id, category_id, title, slug, content, summary, status, published_at, view_count, ai_keywords, created_at, updated_at)
SELECT
  1, NULL,
  'Perf Article #' || g,
  'perf-article-' || g,
  'Performance seed body for article ' || g || '. Lorem ipsum dolor sit amet.',
  'Perf summary ' || g,
  'published',
  NOW() - (g % 365) * INTERVAL '1 day',
  (g * 7) % 1000,
  '[]',
  NOW() - (g % 365) * INTERVAL '1 day',
  NOW()
FROM generate_series(1, 10000) AS g
ON CONFLICT (slug) DO NOTHING;

INSERT INTO forum_categories (name, slug, description, sort, topic_count, created_at, updated_at)
VALUES ('Perf Board', 'perf-board', 'Performance seeding target', 999,
  (SELECT COUNT(*) FROM topics WHERE title LIKE 'Perf Topic #%'), NOW(), NOW())
ON CONFLICT (slug) DO NOTHING;

INSERT INTO topics (user_id, forum_category_id, title, content, status, reply_count, view_count, last_reply_at, created_at, updated_at)
SELECT
  1,
  (SELECT id FROM forum_categories WHERE slug = 'perf-board'),
  'Perf Topic #' || g,
  'Performance seed topic body ' || g,
  'open',
  100,
  (g * 13) % 5000,
  NOW(),
  NOW() - (g % 90) * INTERVAL '1 day',
  NOW()
FROM generate_series(1, 1000) AS g
WHERE NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Perf Topic #' || g)
ON CONFLICT DO NOTHING;

INSERT INTO replies (topic_id, user_id, content, floor, like_count, created_at, updated_at)
SELECT
  t.id,
  1,
  'Performance seed reply ' || n || ' on topic ' || t.id,
  n,
  (n * 3) % 50,
  NOW(),
  NOW()
FROM (SELECT id FROM topics WHERE title LIKE 'Perf Topic #%' ORDER BY id LIMIT 1000) t
CROSS JOIN generate_series(1, 100) AS n
WHERE NOT EXISTS (
  SELECT 1 FROM replies r WHERE r.topic_id = t.id AND r.floor = n
);
