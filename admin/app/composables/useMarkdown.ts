/* Minimal, dependency-free Markdown renderer for public pages.
 * HTML is escaped FIRST so every tag below is ours — user content can never
 * inject markup (XSS-safe by construction). Supports the subset the reply
 * editor toolbar emits: code blocks, inline code, bold, italic, links
 * (http/https only), headings, blockquotes and unordered lists. */

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

export function renderMarkdown(src: string): string {
  let s = escapeHtml(src ?? '')

  // fenced code blocks (before other inline rules)
  s = s.replace(/```[\t ]*([a-zA-Z0-9+#-]*)\n([\s\S]*?)```/g,
    (_m, _lang: string, code: string) =>
      `<pre class="my-2 overflow-x-auto rounded-lg bg-gray-900 p-3 text-xs leading-relaxed text-gray-100"><code>${code.replace(/\n$/, '')}</code></pre>`)

  // headings (line-start)
  s = s.replace(/^### (.*)$/gm, '<h3 class="mt-3 mb-1 text-base font-semibold text-gray-900">$1</h3>')
  s = s.replace(/^## (.*)$/gm, '<h2 class="mt-3 mb-1 text-lg font-semibold text-gray-900">$1</h2>')
  s = s.replace(/^# (.*)$/gm, '<h2 class="mt-3 mb-1 text-lg font-semibold text-gray-900">$1</h2>')

  // blockquote lines (consecutive)
  s = s.replace(/(^|\n)((?:&gt; .*(?:\n|$))+)/g,
    (_m, pre: string, block: string) =>
      `${pre}<blockquote class="my-2 border-l-4 border-gray-300 bg-gray-50 py-1 pl-3 text-sm text-gray-600">${block.replace(/^&gt; /gm, '').replace(/\n$/, '')}</blockquote>`)

  // unordered list items (consecutive)
  s = s.replace(/(^|\n)((?:[-*] .*(?:\n|$))+)/g,
    (_m, pre: string, block: string) =>
      `${pre}<ul class="my-1 list-disc space-y-0.5 pl-5">${block.trimEnd().replace(/^[-*] (.*)$/gm, '<li>$1</li>')}</ul>`)

  // links — only http(s) targets
  s = s.replace(/\[([^\]]+)\]\((https?:\/\/[^\s)]+)\)/g,
    '<a href="$2" target="_blank" rel="noopener noreferrer" class="text-blue-600 hover:underline">$1</a>')

  // bold before italic
  s = s.replace(/\*\*([^*\n]+)\*\*/g, '<strong class="font-semibold">$1</strong>')
  s = s.replace(/\*([^*\n]+)\*/g, '<em>$1</em>')

  // inline code
  s = s.replace(/`([^`\n]+)`/g, '<code class="rounded bg-gray-100 px-1 py-0.5 font-mono text-xs text-pink-600">$1</code>')

  // paragraphs: collapse remaining blank runs, keep single newlines visible
  s = s.replace(/\n{2,}/g, '<br><br>').replace(/\n/g, '<br>')

  return s
}
