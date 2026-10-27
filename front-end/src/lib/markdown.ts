import DOMPurify from 'dompurify'
import { marked } from 'marked'

// Article bodies are stored as Markdown and never converted by the back end.
// This is the only place that assigns HTML: always sanitized first.

DOMPurify.addHook('afterSanitizeAttributes', (node) => {
  if (node instanceof HTMLAnchorElement) {
    node.target = '_blank'
    node.rel = 'noopener noreferrer'
  }
  if (node instanceof HTMLImageElement) {
    const src = node.getAttribute('src') ?? ''
    if (!src.startsWith('https://') && !src.startsWith('/')) node.removeAttribute('src')
  }
})

export function renderMarkdown(target: HTMLElement, md: string): void {
  const raw = marked.parse(md, { async: false, gfm: true, breaks: false })
  target.innerHTML = DOMPurify.sanitize(raw, { USE_PROFILES: { html: true } })
}
