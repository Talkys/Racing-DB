import { ApiError, searchArticles } from '../api/client.ts'
import type { PreviewSegment } from '../api/types.ts'
import { fill, fillList, section, showError } from '../lib/dom.ts'
import { articleHref, entityHref, formatDate } from '../lib/format.ts'

// articles.html?q=<words>[&page=N]. Without q, lists every article newest first.
// Each result shows the text in the browser's language when that text qualifies.

const params = new URLSearchParams(location.search)
const q = params.get('q')?.trim() ?? ''
const page = Math.max(1, Number(params.get('page')) || 1)
const lang = navigator.language.toLowerCase()

const $ = <T extends HTMLElement>(selector: string) => document.querySelector<T>(selector)!

function withPage(n: number): string {
  const next = new URLSearchParams(params)
  next.set('page', String(n))
  return `?${next}`
}

/** Matched words become <mark>; everything stays text, so nothing from the article is parsed as HTML. */
function renderPreview(el: HTMLElement, segments: PreviewSegment[]): void {
  el.replaceChildren(
    ...segments.map((segment) => {
      if (!segment.match) return document.createTextNode(segment.text)
      const mark = document.createElement('mark')
      mark.textContent = segment.text
      return mark
    }),
  )
}

async function load(): Promise<void> {
  const { items, hasNext } = await searchArticles(q, lang, page)

  if (q) document.title = `"${q}" - Articles - RaceDB.org`
  fill(document.body, { summary: q ? `Articles matching "${q}"` : 'All articles, newest first' })

  const results = section('results')
  fillList(
    results,
    items.map((r) => {
      const aboutHref = entityHref(r.entity_type, r.entity_id)
      return {
        title: r.title ?? '(untitled)',
        href: articleHref(r.id, r.language ?? undefined),
        language: r.language ?? 'no text yet',
        kind: r.kind,
        about: r.entity_type,
        aboutHref: aboutHref ?? '',
        aboutPlain: aboutHref ? '' : r.entity_type,
        published: formatDate(r.published_at),
      }
    }),
  )
  results.querySelectorAll<HTMLElement>('[data-preview]').forEach((el, i) => renderPreview(el, items[i].preview))

  $('[data-empty]').hidden = items.length > 0
  $('[data-page]').textContent = String(page)
  const prev = $<HTMLAnchorElement>('[data-prev]')
  prev.hidden = page <= 1
  prev.href = withPage(page - 1)
  const next = $<HTMLAnchorElement>('[data-next]')
  next.hidden = !hasNext
  next.href = withPage(page + 1)
}

$<HTMLInputElement>('[data-article-search] [name="q"]').value = q

try {
  await load()
} catch (err) {
  showError(err instanceof ApiError && err.status === 400 ? new Error(`The server rejected this search: ${err.message}`) : err)
}
