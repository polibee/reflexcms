# ReflexCMS · 可自部署的开源 CMS + 社区论坛

> 基于 **Goravel v1.18（Go）** 与 **Nuxt 4 / Vue 3** 构建的一体化内容与社区平台。开箱即用，支持 `CMS 文章`、`论坛社区` 与 `混合模式` 三种运行形态，内置**积分等级、签到、@提及、私信、拉黑、静态页面管理**等完整社区能力。

ReflexCMS 是一套「文章发布 + 论坛讨论」双引擎的 Web 平台：既能当作简洁的博客/CMS 运营内容，也能架设起带用户体系、发帖回复、积分激励的社区论坛，两种模式可随时在后台切换。它采用清晰的**模块化产品组合架构**，前端基于自维护的声明式后台框架 nuxtadmin，公共前台与管理后台同仓一体。

---

## ✨ 功能总览

### 内容管理（CMS）
- 文章发布：草稿 / 已发布 / 已归档三态流转，定时发布，SEO 标题摘要
- 分类 / 标签体系，中文分词全文检索（unigram + PG fulltext）
- 封面图，浏览量缓存，文章收藏
- 文章评论（审核制）

### 论坛社区（Forum）
- 版块管理，主题帖 / 回复 / 楼层号，置顶 / 精华 / 关闭，最佳回复
- **发帖**：板块选择 + Markdown 编辑器 + 每日发帖积分奖励
- 帖子列表多维度排序（最新 / 最新回复 / 热门），侧边栏版块导航
- 回复分页、@提及通知、拉黑（屏蔽后其回复不再显示）

### 用户与激励
- 会话鉴权（不透明 token，SHA-256 落库）、注册与邮箱域名限制
- RBAC 角色权限 + 可视化权限编辑器
- **积分等级系统**：管理员可配置货币名称、发帖/回复/签到奖励、每日上限、升级阈值
- **用户中心卡片**：等级进度条、今日任务（发帖/评论/签到）、货币余额、六项计数
- 用户公开主页（概览 / 主题帖 / 回复 / 收藏）+ Markdown 个人卡片
- **私信**：Markdown 编辑、收件箱/已发送、未读计数
- **@提及**：回复/评论中 `@用户名` 自动通知

### 平台与扩展
- **静态页面管理**：创建隐私政策、服务条款等页面，**内置富文本（WYSIWYG）编辑器**，前台 `/p/{slug}` 展示
- 布局管理：侧栏卡片（可视化配置、发帖卡片、用户中心卡片）、页眉/页脚导航、轮播图
- 系统设置：站点模式 / 首页内容可视化切换、SEO、注册限制、邮件、隐私、积分等级分组配置
- 审计操作日志、站内通知中心
- 站点三种模式：`CMS` / `Forum` / `Hybrid`（CMS + Forum 混合），首页内容与导航联动

---

## 🚀 快速开始

前置要求：Go ≥ 1.25 · Node ≥ 20 · PostgreSQL 17 · Redis 7

```bash
# 1. 启动 PostgreSQL + Redis（或用 Docker）
make infra            # 或: docker compose -f deploy/docker-compose.yml up -d

# 2. 后端（终端 A）
cp backend/.env.example backend/.env   # 首次；随后 go run . artisan key:generate / jwt:secret
cd backend && go run .                 # http://127.0.0.1:9000 （GET /healthz 探活）

# 3. 管理后台 + 公共前台（终端 B）
cd admin && npm install && npm run dev # http://localhost:3000

# 4. 数据库迁移 / 种子（M1 起使用）
make migrate && make seed
```

默认管理员：`admin@reflexcms.dev` / `ReflexCMS@2026`（可在 `.env` 覆盖）。

---

## 🧱 技术架构

### 后端 —— Goravel v1.18（Go）
- `gin` 驱动，`gorm` + `pgx` / PostgreSQL，`go-redis` 缓存
- **模块化产品组合**：`products` 层是唯一知道所有模块的地方，可组合出 `Full / Blog / Forum` 三种可独立部署的产品
- 模块接口 `kernel.Module{ Name, Migrations, Routes, Boot }`，可选 Seeder / Event / Scheduler Provider
- 通用资源网关 `adminhub`：一套声明式 Spec 驱动全后台 CRUD（搜索/排序/软删/钩子）
- RBAC 通配权限、不透明会话令牌、限流与锁定
- CJK 中文搜索（unigram 分词 + PG fulltext）、SSRF 防护、参数绑定防注入

```
backend/
├── bootstrap/         应用装配（迁移聚合、Provider）
├── products/          产品组合根（Full / Blog / Forum，唯一知晓所有模块处）
├── kernel/            Module 接口定义
├── modules/
│   ├── access/        RBAC、用户、收藏、私信、拉黑、积分面板接口
│   ├── auth/          会话鉴权（不透明 token）
│   ├── cms/           文章/分类/标签
│   ├── comment/       文章评论（审核制 + 提及）
│   ├── forum/         版块/主题/回复/点赞/最佳回复
│   ├── layout/        侧栏卡片/菜单/轮播
│   ├── notification/  站内通知 + @提及
│   ├── pages/         静态页面（富文本）
│   ├── settings/      分组 KV 设置 + 积分/等级服务
│   └── adminhub/      通用资源网关
└── database/migrations/  按模块分组迁移
```

### 前端 —— Nuxt 4 + Vue 3（nuxtadmin 自维护 fork）
- 声明式后台资源定义（表/表单/操作）、TanStack Table、VeeValidate + Zod 校验
- **BFF 代理**：Nitro server 把 `/api/admin/**` 与 `/api/v1/**` 转发到 Goravel 后端，`admin_session` cookie 在服务端转为 Bearer 头
- 框架扩展：`visibleIf` 条件字段、资源级 `transformIn/Out`、`defaultValues` 钩子、WYSIWYG 富文本字段
- 公共前台（文章 / 论坛 / 用户 / 通知 / 私信 / 静态页面）与后台同仓
- 自研零依赖安全 Markdown 渲染器（先转义后转换，免疫 XSS）

---

## 📈 开发进度

| 阶段 | 内容 | 状态 |
|---|---|---|
| M0 | 平台底座、模块化骨架、BFF、基础设施 | ✅ 完成 |
| M1 | 认证（会话）+ RBAC + 通用资源网关 | ✅ 完成 |
| M2 | CMS：文章/分类/标签、CJK 搜索、定时发布 | ✅ 完成 |
| M3 | Forum：版块/主题/回复/点赞、楼层与计数事务、最佳回复 | ✅ 完成 |
| M4 | 社区能力：评论审核、通知、积分、邀请码 | ✅ 完成 |
| M5 | 加固：限流、锁定、SSRF、审计日志、CI、契约测试 | ✅ 完成 |
| 公共前台 | 首页模式化、文章卡片、论坛双列表、发帖、回复 | ✅ 完成 |
| 用户体系 | 用户主页、个人卡片、积分等级、签到、每日任务 | ✅ 完成 |
| 社区进阶 | @提及、私信、拉黑、收藏（文章+帖子）、通知中心 | ✅ 完成 |
| 页面管理 | 静态页面 + 内置 WYSIWYG 富文本编辑器 | ✅ 完成 |
| 布局可视化 | 侧栏卡片可视化配置、页眉/页脚导航、站点模式切换 | ✅ 完成 |

### 规划中
- 搜索增强（Meilisearch 接入）
- OAuth 登录（GitHub / 微信）
- 完整 CKEditor 深度融合
- 前端主题层分离（Phase B）

---

## ✅ 质量闸门

| 目录 | 命令 | 标准 |
|---|---|---|
| backend | `go vet ./...` · `go test ./...` · `go build ./...` | 0 error |
| admin | `npm run lint` · `npm run typecheck` · `npm test` · `npm run build` | 维持 0 基线 |

CI 流水线见 `.github/workflows/`（backend + admin 双条）。

---

## 📚 文档

- [`docs/工程化开发计划.md`](docs/工程化开发计划.md) — 总计划：架构决策、接口契约、里程碑 M0–M5
- [`docs/模块化架构与产品组合.md`](docs/模块化架构与产品组合.md) — 模块边界与产品组合
- [`docs/扩展性架构设计.md`](docs/扩展性架构设计.md) — 插件/主题/扩展点取舍
- [`docs/积分与邀请系统.md`](docs/积分与邀请系统.md) — 激励与等级设计
- [`docs/设置系统与插件架构.md`](docs/设置系统与插件架构.md) — 分组设置可视化
- [`docs/运维部署手册.md`](docs/运维部署手册.md) — 部署运维
- [`docs/性能基线.md`](docs/性能基线.md) — 性能基准

---

## 📄 License

私有 / 待定。仅供学习与二次开发。