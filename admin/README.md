# Nuxt Admin

基于 **Nuxt 4 + Vue 3.5 + TypeScript** 的可扩展后台管理框架，设计思想致敬 Laravel Filament。
不是模板，而是框架：业务以「资源声明」方式接入，CRUD 页面由引擎自动生成。

## 技术栈

| 层 | 选型 |
|---|---|
| 运行时 | Nuxt 4 · Nitro · Vue 3.5 |
| UI | Tailwind CSS v4 · shadcn 风格组件 · Reka UI |
| 表格 | TanStack Table Vue（服务端分页/排序/搜索）|
| 表单 | VeeValidate + Zod（Schema 自动编译校验规则）|
| 状态 | Pinia（仅 UI 偏好与会话，不存服务端数据）|

## 快速开始

```bash
npm install
npm run dev        # http://localhost:3000
```

演示账号（密码均为 `password`）：

- `admin@demo.dev` — 全部权限
- `editor@demo.dev` — 仅内容（posts 增删改、orders/users 只读）
- `viewer@demo.dev` — 全局只读

## 目录结构

```
app/
├── admin/                # 框架层
│   ├── core/             # 类型系统 · AdminRegistry 注册中心 · 事件总线
│   ├── panel/            # definePanel：多面板入口（/admin、未来 /partner…）
│   ├── navigation/       # 由注册表+权限自动生成导航
│   ├── permissions/      # User→Role→Permission 通配匹配 can()
│   ├── schemas/builders/ # DSL：textInput/selectInput/grid/section/textColumn…
│   ├── forms/            # Schema→Zod 编译器 + useFormSchema 引擎
│   ├── tables/           # useResourceTable 服务端表格状态机
│   ├── infolists/        # 详情页 Entry 构建器
│   ├── actions/          # defineAction：确认弹窗/表单弹窗/权限/通知管线
│   ├── widgets/          # 仪表盘 Widget 注册
│   ├── notifications/    # Toast 通知
│   ├── framework/        # 通用渲染组件（DataTable/Resource*Page/ActionHost）
│   └── ui/               # 基础 UI 原子组件
├── modules/              # 业务模块（users/posts/orders/dashboard）
├── pages/admin/[...path] # 资源路由分发
└── plugins/admin.ts      # 组合根：注册 Panel 与全部模块

server/
├── api/auth/             # 登录/会话
├── api/admin/[resource]/ # 通用 REST：list/create/read/update/delete/bulk-delete
└── utils/                # db(仓储) · auth(会话) · resourceConfigs(服务端校验)

shared/types/api.ts       # 前后端共享契约
```

## 声明一个资源（核心用法）

```ts
// app/modules/products/admin/ProductResource.ts
export default defineResource({
  name: 'products',
  label: 'Product',
  labelPlural: 'Products',
  icon: 'package',
  group: 'Catalog',
  searchable: ['name'],

  table: () => [
    textColumn('name', 'Name', { sortable: true }),
    moneyColumn('price', 'Price'),
    badgeColumn('status', 'Status', {
      active: { label: 'Active', variant: 'success' },
      archived: { label: 'Archived', variant: 'secondary' }
    })
  ],

  form: () => [
    section('General', [
      grid(2, [
        textInput('name', 'Name', { required: true }),
        numberInput('price', 'Price', { required: true, min: 0, step: 0.01 }),
      ]),
    ]),
  ],

  infolist: () => [
    textEntry('name', 'Name'),
    moneyEntry('price', 'Price'),
  ],
})
```

在模块中注册即可自动获得：

```
/admin/products          列表（搜索/排序/分页/批量删除/CSV导出）
/admin/products/create   创建（Zod 校验）
/admin/products/:id      详情（Infolist）
/admin/products/:id/edit 编辑
```

## 脚本

```bash
npm run dev         # 开发
npm run build       # 生产构建
npm run preview     # 预览产物
npm run lint        # ESLint（0 错误基线）
npm run typecheck   # vue-tsc（0 错误基线）
npm test            # Vitest 单测（24 个用例，覆盖三个核心纯函数）
```

### 单测覆盖的核心纯函数

| 模块 | 验证点 |
|---|---|
| `admin/forms/schemaToZod.ts` | Schema→Zod 编译：必填/邮箱/数字强转/select·relation 空值拒绝/rules 追加/嵌套布局收集 |
| `admin/permissions/can()` | 通配匹配：`*`、`posts.*` 前缀、精确匹配、无用户拒绝、未声明权限放行 |
| `server/utils/db.ts` `applyQuery` | 大小写不敏感搜索、字符串/数字排序、分页边界、perPage 钳制、非法输入兜底、不可变排序 |

CI 流水线（`.github/workflows/ci.yml`）依次执行 **Lint → Typecheck → Test → Build** 四道闸门。

## 文档导航

| 文档 | 读者 | 内容 |
|---|---|---|
| [`docs/介绍文档.md`](docs/介绍文档.md) | 所有人 | 项目定位、特性总览、架构、快速开始、路线图 |
| [`docs/开发指南.md`](docs/开发指南.md) | 开发者 | 目录职责、核心概念、自动导入清单、新建模块实战、API 参考、规范 |
| [`docs/工程审计报告.md`](docs/工程审计报告.md) | 维护者 | 审计发现与处置、验证矩阵、生产前置清单 |
| [`docs/开发文档.md`](docs/开发文档.md) | 设计溯源 | 原始设计蓝本（Nuxt 版 Filament 设计文档）|

## 当前定位

**框架骨架已完成并验证（构建/静态检查/运行时冒烟全绿），数据层为内存 Demo 实现。**
投入生产前需完成审计报告中列出的 P1/P2 项（真实数据库、真实鉴权、测试覆盖等）。
