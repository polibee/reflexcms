<script setup lang="ts">
/* Tiptap rich-text editor for admin forms (article module, reusable).
 * Stores HTML; toolbar covers the common authoring surface. Link/image
 * insertion uses an in-page modal — no window.prompt. The content area is
 * full-height and clickable everywhere (clicks below the text move the
 * caret to the end instead of doing nothing). */
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import Image from '@tiptap/extension-image'
import TextAlign from '@tiptap/extension-text-align'
import Highlight from '@tiptap/extension-highlight'
import { TextStyle, Color } from '@tiptap/extension-text-style'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import { ref, watch, nextTick } from 'vue'

const props = defineProps<{
  modelValue: string
  placeholder?: string
  disabled?: boolean
}>()

const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const editor = useEditor({
  content: props.modelValue || '',
  editable: !props.disabled,
  extensions: [
    StarterKit.configure({
      heading: { levels: [2, 3] }
    }),
    Link.configure({
      openOnClick: false,
      autolink: true,
      HTMLAttributes: { rel: 'noopener noreferrer', target: '_blank' },
      validate: href => /^https?:\/\//i.test(href)
    }),
    Image.configure({ inline: false, allowBase64: false }),
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Highlight.configure({ multicolor: true }),
    TextStyle,
    Color,
    TaskList,
    TaskItem.configure({ nested: true }),
    Placeholder.configure({ placeholder: props.placeholder ?? '开始编写内容…' })
  ],
  onUpdate: ({ editor: e }) => {
    emit('update:modelValue', e.isEmpty ? '' : e.getHTML())
  }
})

watch(() => props.modelValue, (val) => {
  if (editor.value && editor.value.getHTML() !== val) {
    editor.value.commands.setContent(val || '', false)
  }
})

watch(() => props.disabled, (d) => {
  editor.value?.setEditable(!d)
})

function isActive(name: string, attrs?: Record<string, unknown>) {
  return () => !!editor.value?.isActive(name, attrs)
}

/* Click-to-focus: clicks that land on the wrapper (the empty area below
 * the text) move the caret to the end instead of being swallowed. */
function focusEmptyArea(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.ProseMirror')) {
    editor.value?.commands.focus('end')
  }
}

/* ---- modal ---- */
const dialog = ref<null | 'link' | 'image'>(null)
const dialogUrl = ref('')
const dialogInput = ref<HTMLInputElement | null>(null)

async function openDialog(kind: 'link' | 'image') {
  dialogError.value = ''
  if (kind === 'link') {
    const existing = editor.value?.getAttributes('link').href as string | undefined
    dialogUrl.value = existing ?? ''
  } else {
    dialogUrl.value = ''
  }
  dialog.value = kind
  await nextTick()
  dialogInput.value?.focus()
  dialogInput.value?.select()
}

function closeDialog() {
  dialog.value = null
  editor.value?.commands.focus()
}

function submitDialog() {
  const url = dialogUrl.value.trim()
  if (!/^https?:\/\//i.test(url)) {
    dialogError.value = '地址必须以 http(s):// 开头'
    return
  }
  if (dialog.value === 'link') {
    editor.value?.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
  } else if (dialog.value === 'image') {
    editor.value?.chain().focus().setImage({ src: url, alt: '' }).run()
  }
  closeDialog()
}

function removeLink() {
  editor.value?.chain().focus().extendMarkRange('link').unsetLink().run()
  closeDialog()
}

const dialogError = ref('')

function onDialogKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') { e.preventDefault(); submitDialog() }
  if (e.key === 'Escape') { e.preventDefault(); closeDialog() }
}

function applyColor(v: string) {
  editor.value?.chain().focus().setColor(v).run()
}

/* Toolbar icons: inline SVG path data (lucide-style, 24×24 stroke). Text
 * glyphs like ↺/❝/☑ render as tofu boxes under some Windows fonts, so the
 * toolbar avoids Unicode symbols entirely. */
const ICONS: Record<string, string> = {
  undo: '<path d="M9 14 4 9l5-5"/><path d="M4 9h10.5a5.5 5.5 0 0 1 5.5 5.5a5.5 5.5 0 0 1-5.5 5.5H11"/>',
  redo: '<path d="m15 14 5-5-5-5"/><path d="M20 9H9.5A5.5 5.5 0 0 0 4 14.5A5.5 5.5 0 0 0 9.5 20H13"/>',
  bold: '<path d="M6 12h9a4 4 0 0 1 0 8H7a1 1 0 0 1-1-1V5a1 1 0 0 1 1-1h7a4 4 0 0 1 0 8"/>',
  italic: '<line x1="19" x2="10" y1="4" y2="4"/><line x1="14" x2="5" y1="20" y2="20"/><line x1="15" x2="9" y1="4" y2="20"/>',
  strike: '<path d="M16 4H9a3 3 0 0 0-2.83 4"/><path d="M14 12a4 4 0 0 1 0 8H6"/><line x1="4" x2="20" y1="12" y2="12"/>',
  underline: '<path d="M6 4v6a6 6 0 0 0 12 0V4"/><line x1="4" x2="20" y1="20" y2="20"/>',
  code: '<path d="m16 18 6-6-6-6"/><path d="m8 6-6 6 6 6"/>',
  highlighter: '<path d="m9 11-6 6v3h9l3-3"/><path d="m22 12-4.6 4.6a2 2 0 0 1-2.8 0l-5.2-5.2a2 2 0 0 1 0-2.8L14 4l8 8Z"/>',
  quote: '<path d="M3 21c3 0 7-1 7-8V5c0-1.25-.756-2.017-2-2H4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2 1 0 1 0 1 1v1c0 1-1 2-2 2s-1 .008-1 1.031V20c0 1 0 1 1 1z"/><path d="M15 21c3 0 7-1 7-8V5c0-1.25-.757-2.017-2-2h-4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2h.75c0 2.25.25 4-2.75 4v3c0 1 0 1 1 1z"/>',
  list: '<line x1="8" x2="21" y1="6" y2="6"/><line x1="8" x2="21" y1="12" y2="12"/><line x1="8" x2="21" y1="18" y2="18"/><line x1="3" x2="3.01" y1="6" y2="6"/><line x1="3" x2="3.01" y1="12" y2="12"/><line x1="3" x2="3.01" y1="18" y2="18"/>',
  ordered: '<line x1="10" x2="21" y1="6" y2="6"/><line x1="10" x2="21" y1="12" y2="12"/><line x1="10" x2="21" y1="18" y2="18"/><path d="M4 6h1v4"/><path d="M4 10h2"/><path d="M6 18H4c0-1 2-2 2-3s-1-1.5-2-1"/>',
  tasks: '<path d="m3 17 2 2 4-4"/><path d="m3 7 2 2 4-4"/><path d="M13 6h8"/><path d="M13 12h8"/><path d="M13 18h8"/>',
  alignLeft: '<line x1="21" x2="3" y1="6" y2="6"/><line x1="15" x2="3" y1="12" y2="12"/><line x1="17" x2="3" y1="18" y2="18"/>',
  alignCenter: '<line x1="21" x2="3" y1="6" y2="6"/><line x1="17" x2="7" y1="12" y2="12"/><line x1="19" x2="5" y1="18" y2="18"/>',
  alignRight: '<line x1="21" x2="3" y1="6" y2="6"/><line x1="21" x2="9" y1="12" y2="12"/><line x1="21" x2="7" y1="18" y2="18"/>',
  codeBlock: '<rect width="18" height="18" x="3" y="3" rx="2"/><path d="m10 10-2 2 2 2"/><path d="m14 10 2 2-2 2"/>',
  hr: '<line x1="5" x2="19" y1="12" y2="12"/>',
  link: '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>',
  image: '<rect width="18" height="18" x="3" y="3" rx="2" ry="2"/><circle cx="9" cy="9" r="2"/><path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"/>',
  clear: '<path d="m7 21-4.3-4.3c-1-1-1-2.5 0-3.4l9.6-9.6c1-1 2.5-1 3.4 0l5.6 5.6c1 1 1 2.5 0 3.4L13 21"/><path d="M22 21H7"/><path d="m5 11 9 9"/>'
}

const toolGroups: Array<Array<{ title: string; icon: string; run: () => void; active?: () => boolean }>> = [
  [
    { title: '撤销', icon: 'undo', run: () => editor.value?.chain().focus().undo().run() },
    { title: '重做', icon: 'redo', run: () => editor.value?.chain().focus().redo().run() }
  ],
  [
    { title: '粗体', icon: 'bold', run: () => editor.value?.chain().focus().toggleBold().run(), active: isActive('bold') },
    { title: '斜体', icon: 'italic', run: () => editor.value?.chain().focus().toggleItalic().run(), active: isActive('italic') },
    { title: '删除线', icon: 'strike', run: () => editor.value?.chain().focus().toggleStrike().run(), active: isActive('strike') },
    { title: '下划线', icon: 'underline', run: () => editor.value?.chain().focus().toggleUnderline().run(), active: isActive('underline') },
    { title: '行内代码', icon: 'code', run: () => editor.value?.chain().focus().toggleCode().run(), active: isActive('code') },
    { title: '高亮', icon: 'highlighter', run: () => editor.value?.chain().focus().toggleHighlight().run(), active: isActive('highlight') }
  ],
  [
    { title: '标题', icon: 'H2', run: () => editor.value?.chain().focus().toggleHeading({ level: 2 }).run(), active: isActive('heading', { level: 2 }) },
    { title: '小标题', icon: 'H3', run: () => editor.value?.chain().focus().toggleHeading({ level: 3 }).run(), active: isActive('heading', { level: 3 }) },
    { title: '引用', icon: 'quote', run: () => editor.value?.chain().focus().toggleBlockquote().run(), active: isActive('blockquote') }
  ],
  [
    { title: '无序列表', icon: 'list', run: () => editor.value?.chain().focus().toggleBulletList().run(), active: isActive('bulletList') },
    { title: '有序列表', icon: 'ordered', run: () => editor.value?.chain().focus().toggleOrderedList().run(), active: isActive('orderedList') },
    { title: '任务列表', icon: 'tasks', run: () => editor.value?.chain().focus().toggleTaskList().run(), active: isActive('taskList') }
  ],
  [
    { title: '左对齐', icon: 'alignLeft', run: () => editor.value?.chain().focus().setTextAlign('left').run(), active: isActive('textAlign', { textAlign: 'left' }) },
    { title: '居中', icon: 'alignCenter', run: () => editor.value?.chain().focus().setTextAlign('center').run(), active: isActive('textAlign', { textAlign: 'center' }) },
    { title: '右对齐', icon: 'alignRight', run: () => editor.value?.chain().focus().setTextAlign('right').run(), active: isActive('textAlign', { textAlign: 'right' }) }
  ],
  [
    { title: '代码块', icon: 'codeBlock', run: () => editor.value?.chain().focus().toggleCodeBlock().run(), active: isActive('codeBlock') },
    { title: '分隔线', icon: 'hr', run: () => editor.value?.chain().focus().setHorizontalRule().run() },
    { title: '链接', icon: 'link', run: () => openDialog('link'), active: isActive('link') },
    { title: '图片', icon: 'image', run: () => openDialog('image') },
    { title: '清除格式', icon: 'clear', run: () => editor.value?.chain().focus().unsetAllMarks().clearNodes().run() }
  ]
]
</script>

<template>
  <div class="overflow-hidden rounded-md border border-input bg-background">
    <div class="flex flex-wrap items-center gap-0.5 border-b bg-muted/30 p-1">
      <template v-for="(group, gi) in toolGroups" :key="gi">
        <span v-if="gi > 0" class="mx-1 h-5 w-px bg-border" aria-hidden="true" />
        <button
          v-for="tool in group"
          :key="tool.title"
          type="button"
          :title="tool.title"
          class="inline-flex h-7 w-7 items-center justify-center rounded text-muted-foreground transition-colors hover:bg-accent hover:text-foreground disabled:opacity-40"
          :class="{ 'bg-accent text-foreground shadow-sm': tool.active?.() }"
          :disabled="disabled"
          @click.prevent="tool.run()"
        >
          <svg
            v-if="ICONS[tool.icon]"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-4 w-4"
            aria-hidden="true"
            v-html="ICONS[tool.icon]"
          />
          <span v-else class="text-xs font-semibold">{{ tool.icon }}</span>
        </button>
        <input
          v-if="gi === 1"
          type="color"
          title="文字颜色"
          class="h-7 w-7 cursor-pointer rounded border-0 bg-transparent p-0.5"
          :disabled="disabled"
          @input="applyColor(($event.target as HTMLInputElement).value)"
        >
      </template>
    </div>

    <div class="cursor-text" @click="focusEmptyArea">
      <editor-content
        :editor="editor"
        class="prose prose-sm prose-gray max-w-none px-4 py-3 [&_a]:text-primary [&_a]:underline [&_blockquote]:border-l-4 [&_blockquote]:border-muted [&_blockquote]:pl-3 [&_blockquote]:italic [&_code]:rounded [&_code]:bg-muted [&_code]:px-1 [&_code]:text-xs [&_h2]:mb-1 [&_h2]:mt-3 [&_h2]:text-lg [&_h2]:font-semibold [&_h3]:mb-1 [&_h3]:mt-2 [&_h3]:font-semibold [&_img]:mx-auto [&_img]:max-w-full [&_img]:rounded [&_ol]:list-decimal [&_ol]:pl-5 [&_p]:my-1 [&_pre]:rounded [&_pre]:bg-muted [&_pre]:p-2 [&_ul]:list-disc [&_ul]:pl-5 [&_ul[data-type=taskList]]:list-none [&_ul[data-type=taskList]_li]:flex [&_ul[data-type=taskList]_li]:gap-2 [&_ul[data-type=taskList]_input]:mt-1 [&_mark]:rounded [&_mark]:bg-yellow-200 [&_mark]:px-0.5"
      />
    </div>

    <!-- Link / image modal -->
    <Teleport to="body">
      <div
        v-if="dialog"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm"
        @click.self="closeDialog"
      >
        <div
          class="w-full max-w-md rounded-xl border border-border bg-background p-5 shadow-2xl"
          role="dialog"
          aria-modal="true"
          @keydown="onDialogKeydown"
        >
          <h3 class="mb-1 text-sm font-semibold">
            {{ dialog === 'link' ? '插入链接' : '插入图片' }}
          </h3>
          <p class="mb-3 text-xs text-muted-foreground">
            {{ dialog === 'link' ? '输入要跳转的完整网址' : '输入图片的完整网址（外链）' }}
          </p>
          <input
            ref="dialogInput"
            v-model="dialogUrl"
            type="url"
            :placeholder="dialog === 'link' ? 'https://example.com' : 'https://example.com/image.png'"
            class="mb-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none ring-primary/40 focus:ring-2"
          >
          <p v-if="dialogError" class="mb-2 text-xs text-destructive">
            {{ dialogError }}
          </p>
          <div class="mt-3 flex items-center justify-between">
            <button
              v-if="dialog === 'link' && editor?.isActive('link')"
              type="button"
              class="rounded-md px-3 py-1.5 text-sm text-destructive hover:bg-destructive/10"
              @click="removeLink"
            >
              移除链接
            </button>
            <span v-else />
            <div class="flex gap-2">
              <button
                type="button"
                class="rounded-md border border-border px-3 py-1.5 text-sm hover:bg-accent"
                @click="closeDialog"
              >
                取消
              </button>
              <button
                type="button"
                class="rounded-md bg-primary px-3 py-1.5 text-sm text-primary-foreground hover:bg-primary/90"
                @click="submitDialog"
              >
                确定
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
:deep(.ProseMirror) {
  min-height: 320px;
  outline: none;
}
:deep(.ProseMirror p.is-editor-empty:first-child)::before {
  content: attr(data-placeholder);
  color: hsl(var(--muted-foreground) / 0.6);
  float: left;
  height: 0;
  pointer-events: none;
}
</style>
