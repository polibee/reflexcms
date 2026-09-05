<script setup lang="ts">
/* Tiptap rich-text editor for admin forms (article module).
 * Stores HTML; toolbar covers the common authoring surface. */
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import { watch } from 'vue'

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
      heading: { levels: [2, 3] },
      codeBlock: false
    }),
    Link.configure({
      openOnClick: false,
      autolink: true,
      HTMLAttributes: { rel: 'noopener noreferrer', target: '_blank' }
    }).configure({
      validate: href => /^https?:\/\//i.test(href)
    }),
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

function exec(fn: () => void) {
  fn()
}

function setLink() {
  const url = window.prompt('链接地址（https://…）')
  if (!url) return
  if (!/^https?:\/\//i.test(url)) {
    notifyError('链接必须以 http(s):// 开头')
    return
  }
  editor.value?.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
}

const tools = [
  { label: 'B', title: '粗体', run: () => editor.value?.chain().focus().toggleBold().run(), active: () => editor.value?.isActive('bold') },
  { label: 'I', title: '斜体', run: () => editor.value?.chain().focus().toggleItalic().run(), active: () => editor.value?.isActive('italic') },
  { label: 'S', title: '删除线', run: () => editor.value?.chain().focus().toggleStrike().run(), active: () => editor.value?.isActive('strike') },
  { label: 'H2', title: '标题', run: () => editor.value?.chain().focus().toggleHeading({ level: 2 }).run(), active: () => editor.value?.isActive('heading', { level: 2 }) },
  { label: 'H3', title: '小标题', run: () => editor.value?.chain().focus().toggleHeading({ level: 3 }).run(), active: () => editor.value?.isActive('heading', { level: 3 }) },
  { label: '❝', title: '引用', run: () => editor.value?.chain().focus().toggleBlockquote().run(), active: () => editor.value?.isActive('blockquote') },
  { label: '•', title: '无序列表', run: () => editor.value?.chain().focus().toggleBulletList().run(), active: () => editor.value?.isActive('bulletList') },
  { label: '1.', title: '有序列表', run: () => editor.value?.chain().focus().toggleOrderedList().run(), active: () => editor.value?.isActive('orderedList') },
  { label: '</>', title: '代码块', run: () => editor.value?.chain().focus().toggleCodeBlock().run(), active: () => editor.value?.isActive('codeBlock') },
  { label: '—', title: '分隔线', run: () => editor.value?.chain().focus().setHorizontalRule().run(), active: () => false },
  { label: '🔗', title: '链接', run: setLink, active: () => editor.value?.isActive('link') },
  { label: '⤶', title: '清除格式', run: () => editor.value?.chain().focus().unsetAllMarks().clearNodes().run(), active: () => false }
]
</script>

<template>
  <div class="rounded-md border border-input bg-background">
    <div class="flex flex-wrap gap-1 border-b p-1.5">
      <button
        v-for="tool in tools"
        :key="tool.title"
        type="button"
        :title="tool.title"
        class="rounded px-2 py-1 text-xs text-muted-foreground hover:bg-accent hover:text-foreground"
        :class="{ 'bg-accent text-foreground': tool.active?.() }"
        @click.prevent="exec(tool.run)"
      >
        {{ tool.label }}
      </button>
    </div>
    <editor-content
      :editor="editor"
      class="prose prose-sm prose-gray max-w-none px-3 py-2 min-h-56 [&_p]:my-1 [&_h2]:mt-3 [&_h2]:mb-1 [&_h2]:text-lg [&_h2]:font-semibold [&_h3]:mt-2 [&_h3]:mb-1 [&_h3]:font-semibold [&_ul]:list-disc [&_ul]:pl-5 [&_ol]:list-decimal [&_ol]:pl-5 [&_blockquote]:border-l-4 [&_blockquote]:border-muted [&_blockquote]:pl-3 [&_blockquote]:italic [&_pre]:rounded [&_pre]:bg-muted [&_pre]:p-2 [&_code]:text-xs [&_a]:text-primary [&_a]:underline"
    />
  </div>
</template>
