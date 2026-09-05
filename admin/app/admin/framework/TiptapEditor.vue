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

const toolGroups: Array<Array<{ label: string; title: string; run: () => void; active?: () => boolean }>> = [
  [
    { label: '↺', title: '撤销', run: () => editor.value?.chain().focus().undo().run() },
    { label: '↻', title: '重做', run: () => editor.value?.chain().focus().redo().run() }
  ],
  [
    { label: 'B', title: '粗体', run: () => editor.value?.chain().focus().toggleBold().run(), active: isActive('bold') },
    { label: 'I', title: '斜体', run: () => editor.value?.chain().focus().toggleItalic().run(), active: isActive('italic') },
    { label: 'S', title: '删除线', run: () => editor.value?.chain().focus().toggleStrike().run(), active: isActive('strike') },
    { label: 'U', title: '下划线', run: () => editor.value?.chain().focus().toggleUnderline().run(), active: isActive('underline') },
    { label: '</>', title: '行内代码', run: () => editor.value?.chain().focus().toggleCode().run(), active: isActive('code') },
    { label: '🖍', title: '高亮', run: () => editor.value?.chain().focus().toggleHighlight().run(), active: isActive('highlight') }
  ],
  [
    { label: 'H2', title: '标题', run: () => editor.value?.chain().focus().toggleHeading({ level: 2 }).run(), active: isActive('heading', { level: 2 }) },
    { label: 'H3', title: '小标题', run: () => editor.value?.chain().focus().toggleHeading({ level: 3 }).run(), active: isActive('heading', { level: 3 }) },
    { label: '❝', title: '引用', run: () => editor.value?.chain().focus().toggleBlockquote().run(), active: isActive('blockquote') }
  ],
  [
    { label: '•', title: '无序列表', run: () => editor.value?.chain().focus().toggleBulletList().run(), active: isActive('bulletList') },
    { label: '1.', title: '有序列表', run: () => editor.value?.chain().focus().toggleOrderedList().run(), active: isActive('orderedList') },
    { label: '☑', title: '任务列表', run: () => editor.value?.chain().focus().toggleTaskList().run(), active: isActive('taskList') }
  ],
  [
    { label: '⯇', title: '左对齐', run: () => editor.value?.chain().focus().setTextAlign('left').run(), active: isActive('textAlign', { textAlign: 'left' }) },
    { label: '≡', title: '居中', run: () => editor.value?.chain().focus().setTextAlign('center').run(), active: isActive('textAlign', { textAlign: 'center' }) },
    { label: '⯈', title: '右对齐', run: () => editor.value?.chain().focus().setTextAlign('right').run(), active: isActive('textAlign', { textAlign: 'right' }) }
  ],
  [
    { label: '</>', title: '代码块', run: () => editor.value?.chain().focus().toggleCodeBlock().run(), active: isActive('codeBlock') },
    { label: '—', title: '分隔线', run: () => editor.value?.chain().focus().setHorizontalRule().run() },
    { label: '🔗', title: '链接', run: () => openDialog('link'), active: isActive('link') },
    { label: '🖼', title: '图片', run: () => openDialog('image') },
    { label: '⌫', title: '清除格式', run: () => editor.value?.chain().focus().unsetAllMarks().clearNodes().run() }
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
          class="inline-flex h-7 min-w-7 items-center justify-center rounded px-1.5 text-xs font-medium text-muted-foreground transition-colors hover:bg-accent hover:text-foreground disabled:opacity-40"
          :class="{ 'bg-accent text-foreground shadow-sm': tool.active?.() }"
          :disabled="disabled"
          @click.prevent="tool.run()"
        >
          {{ tool.label }}
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
