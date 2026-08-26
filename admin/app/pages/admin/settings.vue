<script setup lang="ts">
definePageMeta({ layout: 'default' })

/* eslint-disable @typescript-eslint/no-explicit-any */
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

interface FieldDef {
  key: string
  label: string
  type: 'text' | 'switch' | 'number' | 'textarea' | 'select'
  options?: Array<{ label: string, value: string }>
  helpText?: string
}

const GROUP_DEFS: Array<{ group: string, title: string, fields: FieldDef[] }> = [
  {
    group: 'site', title: '站点基本',
    fields: [
      { key: 'site.name', label: '站点名称', type: 'text' },
      { key: 'site.mode', label: '站点模式', type: 'select', helpText: 'CMS 文章站 / 论坛站 / 两者同时展示（运行时可切换）', options: [
        { label: 'CMS + Forum 同时展示', value: 'hybrid' },
        { label: 'CMS 文章站', value: 'blog' },
        { label: 'Forum 论坛站', value: 'forum' }
      ] },
      { key: 'site.home', label: '首页内容', type: 'select', helpText: '首页主体展示 CMS 文章还是论坛帖子', options: [
        { label: '两者都显示', value: 'both' },
        { label: '仅 CMS 文章', value: 'cms' },
        { label: '仅论坛帖子', value: 'forum' }
      ] },
      { key: 'site.description', label: '站点描述', type: 'textarea', helpText: '用于 SEO meta description' },
      { key: 'site.keywords', label: 'SEO 关键词', type: 'text', helpText: '逗号分隔' },
      { key: 'site.icp', label: '备案号', type: 'text' }
    ]
  },
  {
    group: 'registration', title: '注册限制',
    fields: [
      { key: 'reg.enabled', label: '开放注册', type: 'switch' },
      { key: 'reg.mainstream_only', label: '仅主流邮箱', type: 'switch', helpText: 'Gmail、Outlook、QQ 邮箱等' },
      { key: 'reg.email_whitelist', label: '邮箱域名白名单', type: 'textarea', helpText: '每行一个域名；为空则不限制' },
      { key: 'reg.email_blacklist', label: '邮箱域名黑名单', type: 'textarea', helpText: '每行一个域名（永久生效）' },
      { key: 'reg.turnstile_enabled', label: 'Turnstile 人机验证', type: 'switch', helpText: '仅生产环境生效' },
      { key: 'reg.turnstile_site_key', label: 'Turnstile Site Key', type: 'text' },
      { key: 'reg.turnstile_secret_key', label: 'Turnstile Secret Key', type: 'text' }
    ]
  },
  {
    group: 'email', title: '邮件配置',
    fields: [
      { key: 'email.driver', label: '发送驱动', type: 'select', options: [
        { label: 'Log（开发调试）', value: 'log' },
        { label: 'Resend API', value: 'resend' },
        { label: 'SMTP', value: 'smtp' }
      ] },
      { key: 'email.from_address', label: '发件人地址', type: 'text' },
      { key: 'email.from_name', label: '发件人名称', type: 'text' },
      { key: 'email.resend_api_key', label: 'Resend API Key', type: 'text' },
      { key: 'email.smtp_host', label: 'SMTP 主机', type: 'text' },
      { key: 'email.smtp_port', label: 'SMTP 端口', type: 'number' },
      { key: 'email.smtp_user', label: 'SMTP 用户名', type: 'text' },
      { key: 'email.smtp_pass', label: 'SMTP 密码', type: 'text' }
    ]
  },
  {
    group: 'seo', title: 'SEO 与广告验证',
    fields: [
      { key: 'seo.google_site_verification', label: 'Google Search Console', type: 'text' },
      { key: 'seo.baidu_verification', label: '百度站长验证', type: 'text' },
      { key: 'seo.ads_txt', label: 'ads.txt 内容', type: 'textarea' },
      { key: 'seo.head_scripts', label: '<head> 自定义脚本', type: 'textarea', helpText: '广告联盟加载脚本等' }
    ]
  },
  {
    group: 'analytics', title: '统计分析',
    fields: [
      { key: 'analytics.ga_property_id', label: 'GA4 Property ID', type: 'text' },
      { key: 'analytics.ga_credentials_json', label: '服务账号 JSON', type: 'textarea' }
    ]
  },
  {
    group: 'ai', title: 'AI 设置',
    fields: [
      { key: 'ai.provider', label: '模型厂商', type: 'select', options: [{ label: 'OpenAI', value: 'openai' }] },
      { key: 'ai.model', label: '模型标识', type: 'text' },
      { key: 'ai.api_key', label: 'API Key', type: 'text' },
      { key: 'ai.base_url', label: '自定义端点', type: 'text' }
    ]
  },
  {
    group: 'points', title: '积分与等级',
    fields: [
      { key: 'points.currency_name', label: '货币名称', type: 'text', helpText: '社区货币的展示名，如"鸡腿"' },
      { key: 'points.level_thresholds', label: '等级阈值', type: 'text', helpText: '累计获得货币的升级阈值，逗号分隔，如 0,100,300,600,1000' },
      { key: 'points.earn_topic', label: '发帖奖励', type: 'number' },
      { key: 'points.topic_daily_cap', label: '发帖每日上限', type: 'number', helpText: '每日通过发帖最多可获得的货币量' },
      { key: 'points.earn_reply', label: '回复奖励', type: 'number' },
      { key: 'points.reply_daily_cap', label: '回复每日上限', type: 'number' },
      { key: 'points.earn_signin', label: '签到奖励', type: 'number', helpText: '每日签到一次的奖励' }
    ]
  },
  {
    group: 'privacy', title: '隐私与 Cookie',
    fields: [
      { key: 'privacy.cookie_banner', label: 'Cookie 同意横幅', type: 'switch' },
      { key: 'privacy.privacy_policy_url', label: '隐私政策链接', type: 'text' },
      { key: 'privacy.policy_updated_at', label: '政策更新日期', type: 'text' }
    ]
  }
]

const settings = ref<Record<string, Record<string, any>>>({})
const loading = ref(false)
const saving = ref(false)

async function fetchAll() {
  loading.value = true
  try {
    settings.value = await rawFetch('/api/admin/settings-all') as Record<string, Record<string, any>>
  } catch { /* silent */ } finally { loading.value = false }
}

function getVal(group: string, key: string): any {
  return settings.value[group]?.[key] ?? ''
}

function setVal(group: string, key: string, val: any) {
  if (!settings.value[group]) settings.value[group] = {}
  settings.value[group][key] = val
}

async function saveAll() {
  saving.value = true
  try {
    const values: Record<string, any> = {}
    for (const g of GROUP_DEFS) {
      for (const f of g.fields) {
        values[f.key] = getVal(g.group, f.key)
      }
    }
    await rawFetch('/api/admin/settings-batch', {
      method: 'PUT',
      body: { values }
    })
    notify('设置已保存')
  } catch (e: any) {
    notifyError(e?.message ?? '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(fetchAll)

useHead({ title: '站点设置' })
</script>

<template>
  <div class="mx-auto max-w-4xl space-y-8 p-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">
        站点设置
      </h1>
      <UiButton
        :disabled="saving || loading"
        @click="saveAll"
      >
        {{ saving ? 'Saving…' : '保存全部设置' }}
      </UiButton>
    </div>

    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>

    <template v-else>
      <div
        v-for="group in GROUP_DEFS"
        :key="group.title"
        class="space-y-4 rounded-xl border bg-white p-6"
      >
        <h2 class="text-base font-semibold text-gray-900">
          {{ group.title }}
        </h2>
        <div class="grid grid-cols-1 gap-x-6 gap-y-4 md:grid-cols-2">
          <template
            v-for="field in group.fields"
            :key="field.key"
          >
            <div
              v-if="field.type === 'switch'"
              class="flex items-center justify-between rounded-lg border p-3"
            >
              <div>
                <span class="text-sm">{{ field.label }}</span>
                <p
                  v-if="field.helpText"
                  class="mt-0.5 text-xs text-gray-400"
                >
                  {{ field.helpText }}
                </p>
              </div>
              <UiSwitch
                :model-value="Boolean(getVal(group.group, field.key))"
                @update:model-value="setVal(group.group, field.key, $event)"
              />
            </div>

            <div
              v-else-if="field.type === 'select'"
              class="flex items-center gap-3"
            >
              <span class="w-32 shrink-0 text-sm text-gray-600">{{ field.label }}</span>
              <div class="flex-1">
                <select
                  :value="getVal(group.group, field.key)"
                  class="h-9 w-full rounded-md border border-input bg-background px-3 text-sm"
                  @change="setVal(group.group, field.key, ($event.target as HTMLSelectElement).value)"
                >
                  <option
                    v-for="opt in field.options"
                    :key="String(opt.value)"
                    :value="opt.value"
                  >
                    {{ opt.label }}
                  </option>
                </select>
                <p
                  v-if="field.helpText"
                  class="mt-0.5 text-xs text-gray-400"
                >
                  {{ field.helpText }}
                </p>
              </div>
            </div>

            <div
              v-else-if="field.type === 'textarea'"
              class="md:col-span-2"
            >
              <label class="mb-1 block text-sm text-gray-600">{{ field.label }}</label>
              <textarea
                :rows="3"
                :placeholder="field.helpText"
                class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                :value="getVal(group.group, field.key)"
                @input="setVal(group.group, field.key, ($event.target as HTMLTextAreaElement).value)"
              />
            </div>

            <div v-else>
              <label class="mb-1 block text-sm text-gray-600">{{ field.label }}</label>
              <input
                :type="field.type === 'number' ? 'number' : 'text'"
                class="h-9 w-full rounded-md border border-input bg-background px-3 text-sm"
                :value="getVal(group.group, field.key)"
                @input="setVal(group.group, field.key, ($event.target as HTMLInputElement).value)"
              >
              <p
                v-if="field.helpText"
                class="mt-0.5 text-xs text-gray-400"
              >
                {{ field.helpText }}
              </p>
            </div>
          </template>
        </div>
      </div>
    </template>
  </div>
</template>
