<script setup lang="ts">
import { useField, useFormValues } from 'vee-validate'
import type { FieldOption, FieldNode } from '~/admin/core/types'
import type { Paginated } from '#shared/types/api'
import type { RowRecord } from '~/admin/tables/useResourceTable'

const props = defineProps<{ node: FieldNode }>()

/* conditional visibility: evaluated against the live form values so fields
 * can show/hide as their dependencies change (e.g. per-type widget config) */
const formValues = useFormValues()
const hidden = computed(() =>
  props.node.visibleIf ? !props.node.visibleIf((formValues.value ?? {}) as Record<string, unknown>) : false)

const name = toRef(() => props.node.name)
const { value, errorMessage } = useField<string | number | boolean | null>(name)

const { value: arrayValue, errorMessage: arrayErrorMessage } = useField<string[]>(name)
const arrayItems = computed<string[]>(() =>
  Array.isArray(arrayValue.value) ? arrayValue.value.map(String) : []
)
function toggleOption(opt: FieldOption, checked: boolean) {
  const set = new Set(arrayItems.value)
  const key = String(opt.value)
  if (checked) set.add(key)
  else set.delete(key)
  arrayValue.value = [...set]
}

const stringValue = computed<string>({
  get: () => (value.value ?? '') as string,
  set: (v) => { value.value = v }
})

const relationOptions = ref<FieldOption[]>([])
onMounted(async () => {
  const rel = props.node.relation
  if (props.node.kind !== 'relation' || !rel) return
  try {
    const res = await $fetch<Paginated<RowRecord>>(`/api/admin/${rel.resource}`, {
      query: { perPage: 200 }
    })
    relationOptions.value = res.items.map(r => ({
      label: String(r[rel.labelKey] ?? r.id),
      value: r.id as string | number
    }))
  } catch {
    relationOptions.value = []
  }
})

const options = computed<FieldOption[]>(() =>
  props.node.kind === 'relation' ? relationOptions.value : (props.node.options ?? [])
)

const inputType = computed(() => {
  switch (props.node.kind) {
    case 'email': return 'email'
    case 'password': return 'password'
    case 'number': return 'number'
    case 'date': return 'date'
    default: return 'text'
  }
})

/* upload */
const uploadBusy = ref(false)
async function onUpload(e: Event) {
  const inputEl = e.target as HTMLInputElement
  const file = inputEl.files?.[0]
  if (!file) return
  uploadBusy.value = true
  try {
    const fd = new FormData()
    fd.append('file', file)
    const res = await $fetch<{ url: string }>('/api/admin/uploads', { method: 'POST', body: fd })
    value.value = res.url
  } finally {
    uploadBusy.value = false
    inputEl.value = ''
  }
}

/* markdown toolbar for textarea/richtext */
const mdRef = ref<HTMLTextAreaElement>()
function wrapSel(before: string, after = before) {
  const el = mdRef.value
  if (!el) return
  const s = el.selectionStart
  const e = el.selectionEnd
  const sel = el.value.slice(s, e)
  el.value = el.value.slice(0, s) + before + sel + after + el.value.slice(e)
  el.selectionStart = s + before.length
  el.selectionEnd = s + before.length + sel.length
  el.focus()
  value.value = el.value
}

const isRichtext = computed(() => props.node.kind === 'richtext')

/* ---------- wysiwyg: contenteditable rich text (dependency-free) ---------- */
const wysiwygRef = ref<HTMLElement>()
const wysiwygHtml = computed<string>({
  get: () => (value.value ?? '') as string,
  set: (v) => { value.value = v }
})

const wysiwygTools = [
  { cmd: 'bold', label: 'B', title: '粗体', cls: 'font-bold' },
  { cmd: 'italic', label: 'I', title: '斜体', cls: 'italic' },
  { cmd: 'underline', label: 'U', title: '下划线', cls: 'underline' },
  { cmd: 'formatBlock', arg: '<h2>', label: 'H2', title: '标题', cls: 'font-semibold' },
  { cmd: 'formatBlock', arg: '<h3>', label: 'H3', title: '小标题', cls: 'font-semibold' },
  { cmd: 'formatBlock', arg: '<blockquote>', label: '❝', title: '引用', cls: '' },
  { cmd: 'insertUnorderedList', label: '•', title: '无序列表', cls: '' },
  { cmd: 'insertOrderedList', label: '1.', title: '有序列表', cls: '' },
  { cmd: 'removeFormat', label: '✕', title: '清除格式', cls: '' }
]

function execWysiwyg(tool: { cmd: string, arg?: string }) {
  const el = wysiwygRef.value
  if (!el) return
  el.focus()
  document.execCommand(tool.cmd, false, tool.arg ?? undefined)
  wysiwygHtml.value = el.innerHTML
}

function insertWysiwygLink() {
  const el = wysiwygRef.value
  if (!el) return
  const url = window.prompt('链接地址（https://…）')
  if (!url) return
  if (!/^https?:\/\//i.test(url)) {
    notifyError('链接必须以 http(s):// 开头')
    return
  }
  el.focus()
  document.execCommand('createLink', false, url)
  wysiwygHtml.value = el.innerHTML
}

function onWysiwygInput(e: Event) {
  wysiwygHtml.value = (e.target as HTMLElement).innerHTML
}
</script>

<template>
  <div
    v-if="!hidden"
    class="space-y-1.5"
  >
    <UiLabel
      v-if="node.kind !== 'switch' && node.kind !== 'checkbox'"
      :for="node.name"
    >
      {{ node.label }}<span
        v-if="node.required"
        class="text-destructive"
      > *</span>
    </UiLabel>

    <UiInput
      v-if="['text', 'email', 'password', 'number', 'date'].includes(node.kind)"
      :id="node.name"
      v-model="stringValue"
      :type="inputType"
      :placeholder="node.placeholder"
      :disabled="node.disabled"
    />

    <!-- textarea / richtext with markdown toolbar -->
    <template v-else-if="node.kind === 'textarea' || isRichtext">
      <div
        v-if="isRichtext"
        class="mb-1 flex gap-1"
      >
        <button
          v-for="btn in [
            { label: 'B', title: '粗体', before: '**', after: '**', cls: 'font-bold' },
            { label: 'I', title: '斜体', before: '*', after: '*', cls: 'italic' },
            { label: 'H2', title: '标题', before: '## ', after: '', cls: 'font-semibold' },
            { label: 'H3', title: '子标题', before: '### ', after: '', cls: 'font-semibold' },
            { label: '🔗', title: '链接', before: '[', after: '](url)', cls: '' },
            { label: '</>', title: '代码', before: '`', after: '`', cls: 'font-mono text-xs' },
            { label: '📝', title: '引用', before: '> ', after: '', cls: '' },
            { label: '•', title: '列表', before: '- ', after: '', cls: '' }
          ]"
          :key="btn.title"
          type="button"
          :title="btn.title"
          class="rounded border border-input px-2 py-0.5 text-xs hover:bg-accent"
          :class="btn.cls"
          @click.prevent="wrapSel(btn.before, btn.after)"
        >
          {{ btn.label }}
        </button>
      </div>
      <textarea
        :id="node.name"
        :ref="isRichtext ? (el: any) => { mdRef = el.$el ?? el } : undefined"
        v-model="stringValue"
        :rows="isRichtext ? (node.rows ?? 12) : (node.rows ?? 4)"
        :placeholder="node.placeholder"
        :disabled="node.disabled"
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
      />
    </template>

    <!-- wysiwyg rich text editor (contenteditable, no dependencies) -->
    <template v-else-if="node.kind === 'wysiwyg'">
      <div class="mb-1 flex flex-wrap gap-1">
        <button
          v-for="tool in wysiwygTools"
          :key="tool.title"
          type="button"
          :title="tool.title"
          class="rounded border border-input px-2 py-0.5 text-xs hover:bg-accent"
          :class="tool.cls"
          @click.prevent="execWysiwyg(tool)"
        >
          {{ tool.label }}
        </button>
        <button
          type="button"
          title="链接"
          class="rounded border border-input px-2 py-0.5 text-xs hover:bg-accent"
          @click.prevent="insertWysiwygLink"
        >
          🔗
        </button>
      </div>
      <div
        :id="node.name"
        ref="wysiwygRef"
        contenteditable="true"
        class="min-h-40 w-full rounded-md border border-input bg-background px-3 py-2 text-sm leading-relaxed focus:outline-none focus:ring-1 focus:ring-ring [&_blockquote]:border-l-4 [&_blockquote]:pl-3 [&_blockquote]:text-muted-foreground [&_h2]:text-lg [&_h2]:font-semibold [&_h3]:font-semibold [&_ol]:list-decimal [&_ol]:pl-6 [&_ul]:list-disc [&_ul]:pl-6"
        @input="onWysiwygInput"
        @blur="wysiwygHtml = ($event.target as HTMLElement).innerHTML"
      />
    </template>

    <template v-else-if="node.kind === 'select' || node.kind === 'relation'">
      <UiSelect
        v-if="options.length > 0 || node.kind === 'select'"
        :id="node.name"
        v-model="stringValue"
        :options="options"
        :placeholder="node.required ? `Select ${node.label.toLowerCase()}` : undefined"
        :disabled="node.disabled"
      />
      <div
        v-else
        class="flex h-9 items-center rounded-md border border-input px-3 text-sm text-muted-foreground"
      >
        Loading…
      </div>
    </template>

    <div
      v-else-if="node.kind === 'checkboxgroup'"
      class="grid grid-cols-2 gap-x-4 gap-y-2"
    >
      <UiCheckbox
        v-for="opt in options"
        :key="String(opt.value)"
        :model-value="arrayItems.includes(String(opt.value))"
        :disabled="node.disabled"
        :label="opt.label"
        @update:model-value="toggleOption(opt, $event)"
      />
    </div>

    <div
      v-else-if="node.kind === 'upload'"
      class="space-y-2"
    >
      <img
        v-if="stringValue"
        :src="stringValue"
        alt=""
        class="h-28 w-auto rounded-md border object-cover"
      >
      <label
        class="inline-flex cursor-pointer items-center rounded-md border border-input bg-background px-3 py-1.5 text-sm hover:bg-accent"
        :class="{ 'pointer-events-none opacity-60': uploadBusy }"
      >
        {{ uploadBusy ? 'Uploading…' : '选择文件…' }}
        <input
          type="file"
          class="hidden"
          :accept="node.accept ?? 'image/*'"
          @change="onUpload"
        >
      </label>
      <span
        v-if="stringValue"
        class="ml-2 text-xs text-gray-400"
      >{{ stringValue }}</span>
    </div>

    <UiSwitch
      v-else-if="node.kind === 'switch'"
      :model-value="Boolean(value)"
      :disabled="node.disabled"
      @update:model-value="value = $event"
    >
      <span v-if="!node.helpText">{{ node.label }}</span>
    </UiSwitch>

    <UiCheckbox
      v-else-if="node.kind === 'checkbox'"
      :model-value="Boolean(value)"
      :disabled="node.disabled"
      :label="node.label"
      @update:model-value="value = $event"
    />

    <p
      v-if="errorMessage || arrayErrorMessage"
      class="text-xs text-destructive"
    >
      {{ errorMessage || arrayErrorMessage }}
    </p>
    <p
      v-else-if="node.helpText"
      class="text-xs text-muted-foreground"
    >
      {{ node.helpText }}
    </p>
  </div>
</template>
