# ReflexCMS · Self-Hosted Open-Source CMS + Community Forum

> An all-in-one content and community platform built with **Goravel v1.18 (Go)** and **Nuxt 4 / Vue 3**. Runs as `CMS`, `Forum`, or `Hybrid`, with a full community toolkit: **points & levels, daily check-ins, @mentions, direct messages, blocking, shop with payment gateways, and a static page manager**.

> 🤖 **This project was built entirely by Zhipu's GLM 5.3 FLASH** — backend architecture, frontend UI, and deployment scripts were all AI-coded; the human only supplied requirements and acceptance.

[中文](README.md) | English

## 📷 Screenshots

| Home (Hybrid mode) | Forum |
|---|---|
| ![Home](docs/screenshots/home-hybrid.png) | ![Forum](docs/screenshots/forums.png) |
| **Topic Detail** | **Admin Panel** |
| ![Topic Detail](docs/screenshots/topic-detail.png) | ![Admin Panel](docs/screenshots/admin-dashboard.png) |
| **User Center** | **Shop** |
| ![User Center](docs/screenshots/user-profile.png) | ![Shop](docs/screenshots/shop.png) |

---

## ✨ Features

### Content Management (CMS)
- Articles: draft / published / archived states, scheduled publishing, SEO meta
- Categories & tags, CJK-aware full-text search (unigram + PG fulltext)
- Cover images, buffered view counters, article favorites, moderated comments with @mentions

### Forum
- Boards, topics / replies / floor numbers, pin / feature / close / best reply
- **Posting**: board picker + Markdown editor + earn points per topic
- Reply pagination, @mention notifications, user blocking, **moderator system** (per-board assignment, mute/ban governance)
- **Invite-only registration**: one toggle in the admin panel; invite codes purchasable in the shop and auto-delivered after payment

### Users & Gamification
- Opaque-token session auth, RBAC roles with a visual permission editor
- **Points & levels**: currency name, per-action rewards, daily caps, level thresholds — all admin-configurable
- User-center card: level progress, daily tasks, currency balance, six counters
- Public profile (overview / topics / replies / favorites) + Markdown profile card + reply signatures
- **Direct messages**: Markdown editor, inbox/sent, unread counts

### Platform & Extensibility
- **Shop + payment gateway plugins**: Xcash / NOWPayments / CoinPayments / PayPal implemented; XunhuPay / CodePay (MD5 protocol) bundled behind a security-policy exception
- **Static page manager**: built-in WYSIWYG editor, public `/p/{slug}` pages (privacy policy, ToS, …)
- Layout: sidebar widgets with visual config, header/footer menus, carousels
- Runtime-switchable site modes: `CMS` / `Forum` / `Hybrid`
- Bulk email (manual or CSV import), audit logs, i18n (zh-CN / English)

---

## 🚀 One-Click Deploy

```bash
git clone https://github.com/polibee/reflexcms.git
cd reflexcms
./deploy/deploy.sh
```

The script handles everything: environment checks → Docker (PostgreSQL 17 + Redis 7) → migrations & seeds → super-admin bootstrap → backend build → frontend build. Then open `http://localhost:3000`.

> **Production admin account**: on a fresh install the deploy script runs `artisan admin:bootstrap`, which generates a strong random password and prints it **exactly once** — store it immediately. You can also run it manually (from `backend/`): `go run . artisan admin:bootstrap --email=you@example.com` (omit `--password` to auto-generate). The first account registered on an empty site also becomes the super-admin. Admins can change the password anytime under Settings → personal profile. The dev default `admin@reflexcms.dev` / `ReflexCMS@2026` must be reset before going live.

## 🛠 Manual Deployment

To deploy without the one-click script, follow these steps (full details and troubleshooting in [`docs/运维部署手册.md`](docs/运维部署手册.md)).

**1. Requirements**: Go ≥ 1.25, Node.js ≥ 20, PostgreSQL ≥ 15 (17 recommended), Redis ≥ 6.

**2. Create the database**:

```sql
CREATE USER reflexcms WITH PASSWORD 'your-strong-password';
CREATE DATABASE reflexcms OWNER reflexcms;
```

**3. Backend** (listens on `:9000` by default):

```bash
cd backend
cp .env.example .env                  # edit DB_* / APP_KEY / APP_ENV=production
go build -o reflexcms-api .
./reflexcms-api artisan migrate       # run migrations
./reflexcms-api artisan db:seed       # first time only: roles & base data
./reflexcms-api artisan admin:bootstrap   # create super-admin; password printed once
./reflexcms-api                       # start (use systemd in production)
```

**4. Frontend** (Nuxt server app hosting both the admin panel and the public site, default `:3000`):

```bash
cd admin
echo "NUXT_PUBLIC_API_BASE=http://127.0.0.1:9000" > .env
npm ci && npm run build
BACKEND_URL=http://127.0.0.1:9000 node .output/server/index.mjs
```

**5. Process management & reverse proxy**: host the backend under systemd (`Restart=always`) or NSSM (Windows); terminate TLS at Caddy/Nginx and set `SESSION_SECURE=true` plus a production build of admin to enable Secure cookies.

For manual / production deployment, see the **Manual Deployment** section above and [`docs/运维部署手册.md`](docs/运维部署手册.md).

---

## 🧱 Architecture

**Backend**: Goravel v1.18 (Go ≥ 1.25) · gin · gorm/pgx · PostgreSQL 17 · Redis 7
**Frontend**: Nuxt 4 · Vue 3.5 · TanStack Table · VeeValidate + Zod · BFF proxy

- **Modular product composition**: the `products` layer is the only place that knows all modules; `Full / Blog / Forum` are independently deployable products
- Generic resource gateway `adminhub`: declarative specs drive all admin CRUD
- Payment gateways are plugins: a `gateways.Gateway` interface with self-registration — adding a provider is a single file
- Zero-dependency XSS-safe Markdown renderer (escape-first, then transform)
- CJK full-text search, SSRF protection, full parameter binding

---

## 📈 Roadmap Status

| Phase | Scope | Status |
|---|---|---|
| M0 | Foundation, modular skeleton, BFF | ✅ |
| M1 | Auth + RBAC + generic resource gateway | ✅ |
| M2 | CMS: articles, categories, tags, CJK search, scheduled publishing | ✅ |
| M3 | Forum: boards, topics, replies, likes, transactional floors | ✅ |
| M4 | Comment moderation, notifications, points, invites | ✅ |
| M5 | Rate limiting, lockout, SSRF protection, audit log, CI | ✅ |
| Public frontend | Mode-aware home, article cards, forum feeds, posting | ✅ |
| User system | Profiles, cards, points & levels, check-ins, DMs | ✅ |
| Governance | Moderators, mute/ban, blocking, mentions, notifications | ✅ |
| Commerce | Products/orders/stock, 5 gateway plugins, invite-code goods | ✅ |
| Page manager | Static pages + WYSIWYG editor | ✅ |
| Platform | i18n (zh-CN / en), deploy script, bulk email | ✅ |

**Planned**: Meilisearch integration, OAuth login, theme layer separation.

---

<a id="bigmodel"></a>
## 🤖 How This Project Was Built & BigModel Promo

**This project was built entirely by GLM 5.3 FLASH** — backend architecture, frontend UI, payment gateways, deployment scripts, and docs were all coded by the GLM 5.3 FLASH model on [Zhipu's BigModel.cn open platform](https://www.bigmodel.cn/); the human only supplied requirements and acceptance.

**I'm building AI applications on BigModel.cn (Zhipu AI's open platform).** Their new flagship model **GLM-5.3** is now live, reaching open-source SOTA in **reasoning, coding, and agentic** capabilities. Sign up through my invite link to get a **20-million Tokens welcome pack**:

**[👉 Invite link: https://www.bigmodel.cn/invite?icode=LHOTDtGgt9PRlQPrnhcrXnHEaazDlIZGj9HxftzTbt4%3D](https://www.bigmodel.cn/invite?icode=LHOTDtGgt9PRlQPrnhcrXnHEaazDlIZGj9HxftzTbt4%3D)**

- 🧠 **GLM-5.3 flagship**: open-source SOTA in reasoning / coding / agents
- 🎁 **Invite bonus**: 20M Tokens for sign-ups via the link above
- 🛠 **Proof of full-lifecycle AI development**: the 60+ commits in this repo are GLM 5.3 FLASH's deliverables

<a id="vast-ai"></a>
## ☁️ Need GPU Compute?

**[Vast — cost-effective GPU cloud rental](https://cloud.vast.ai/?ref_id=91181)**: if you need GPUs for AI summarization, content moderation, or model hosting, try [Vast](https://cloud.vast.ai/?ref_id=91181) — per-second billing, a huge community GPU marketplace, at a fraction of mainstream cloud pricing.

- **The world's largest GPU marketplace**: 20,000+ rentable GPUs across 40+ data centers, from RTX 4090 to A100 / H100
- **Per-second billing, stop anytime**: prices are set by real-time supply and demand — usually a fraction of what hyperscalers charge; a few dollars gets you started
- **Deploy in minutes**: pick a PyTorch / TensorFlow image, launch an instance, connect via SSH or Jupyter
- **Flexible form factors**: single-GPU instances, large elastic clusters, and Serverless inference endpoints

## 📄 License

MIT
