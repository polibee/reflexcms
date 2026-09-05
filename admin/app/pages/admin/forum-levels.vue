<script setup lang="ts">
definePageMeta({ layout: 'default' })

/* Community growth settings under the Forum nav group: currency, earn
 * rates, daily caps, level ladder and level perks. Edits go through the
 * same settings-batch endpoint as the main settings page. */

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

interface Field {
  key: string
  label: string
  helpText?: string
  number?: boolean
}

const GROUPS: Array<{ title: string, fields: Field[] }> = [
  {
    title: '货币与奖励',
    fields: [
      { key: 'points.currency_name', label: '货币名称' },
      { key: 'points.earn_topic', label: '发帖奖励', number: true },
      { key: 'points.topic_daily_cap', label: '发帖每日上限', number: true, helpText: '每日通过发帖最多可获得的货币量' },
      { key: 'points.earn_reply', label: '回复奖励', number: true },
      { key: 'points.reply_daily_cap', label: '回复每日上限', number: true },
      { key: 'points.earn_signin', label: '签到奖励', number: true, helpText: '每日签到一次的奖励' }
    ]
  },
  {
    title: '等级成长',
    fields: [
      { key: 'points.level_thresholds', label: '升级阈值', helpText: '累计货币达到阈值即升级，逗号分隔升序，如 0,100,300,600,1000' },
      { key: 'points.level_names', label: '等级名称', helpText: '与阈值一一对应的称号，逗号分隔，如 新手,活跃,达人,元老,传奇' },
      { key: 'points.level_perks', label: '各等级特权说明', helpText: '展示在用户中心，用 | 分隔每个等级的特权，如 发帖|回复|自定义头衔|专属徽章' }
    ]
  }
]

const settings = ref<Record<string, any>>({})
const loading = ref(false)
const saving = ref(false)
const msg = ref('')

function getVal(key: string): any {
  return settings.value[key] ?? ''
}

function setVal(key: string, val: any) {
  settings.value[key] = val
}

async function fetchAll() {
  loading.value = true
  try {
    const all = await rawFetch('/api/admin/settings-all') as Record<string, Record<string, any>>
    const flat: Record<string, any> = {}
    for (const group of Object.values(all)) {
      Object.assign(flat, group)
    }
    settings.value = flat
  } catch { /* silent */ } finally { loading.value = false }
}

async function save() {
  saving.value = true
  msg.value = ''
  try {
    const values: Record<string, any> = {}
    for (const g of GROUPS) {
      for (const f of g.fields) values[f.key] = getVal(f.key)
    }
    await rawFetch('/api/admin/settings-batch', { method: 'PUT', body: { values } })
    msg.value = '已保存，前台即时生效'
  } catch (e: unknown) {
    const err = e as { data?: { statusMessage?: string, message?: string } }
    msg.value = err.data?.statusMessage ?? '保存失败'
  } finally {
    saving.value = false
  }
}

onMounted(fetchAll)
useHead({ title: '等级与签到设置' })
</script>

<template>
  <div class="mx-auto max-w-3xl space-y-6 p-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">
        等级与签到设置
      </h1>
      <button
        type="button"
        :disabled="saving || loading"
        class="rounded-lg bg-blue-600 px-5 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        @click="save"
      >
        {{ saving ? '保存中…' : '保存设置' }}
      </button>
    </div>

    <p
      v-if="msg"
      class="text-sm"
      :class="msg.includes('失败') ? 'text-red-600' : 'text-green-600'"
    >
      {{ msg }}
    </p>

    <div
      v-for="group in GROUPS"
      :key="group.title"
      class="space-y-4 rounded-xl border bg-white p-6"
    >
      <h2 class="text-base font-semibold text-gray-900">
        {{ group.title }}
      </h2>
      <div class="grid grid-cols-1 gap-x-6 gap-y-4 md:grid-cols-2">
        <div
          v-for="f in group.fields"
          :key="f.key"
          :class="f.number ? '' : 'md:col-span-2'"
        >
          <label class="mb-1 block text-sm text-gray-600">{{ f.label }}</label>
          <input
            :type="f.number ? 'number' : 'text'"
            class="h-9 w-full rounded-md border border-input bg-background px-3 text-sm"
            :value="getVal(f.key)"
            @input="setVal(f.key, ($event.target as HTMLInputElement).value)"
          >
          <p
            v-if="f.helpText"
            class="mt-0.5 text-xs text-gray-400"
          >
            {{ f.helpText }}
          </p>
        </div>
      </div>
    </div>

    <div class="rounded-xl border bg-white p-6 text-sm text-gray-500">
      <h2 class="mb-2 font-semibold text-gray-900">
        等级特权说明
      </h2>
      <p>等级由用户累计获得的货币自动计算（不可手动指定）。升级阈值与等级名称一一对应；特权说明会展示在前台用户中心，作为激励引导。需要按等级限制的具体功能（如自定义头衔、专属版块）可在后续版本中接入 <code class="rounded bg-gray-100 px-1">points.level_names</code> 判断。</p>
    </div>
  </div>
</template>
