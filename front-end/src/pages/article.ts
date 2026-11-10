import { get } from '../api/client.ts'
import { textsOf } from '../lib/articles.ts'
import { fill, fillList, section, showError } from '../lib/dom.ts'
import { aboutFields, resolveRefs } from '../lib/entities.ts'
import { articleHref, formatDate, formatDateTime } from '../lib/format.ts'
import { pickLanguage } from '../lib/language.ts'
import { renderMarkdown } from '../lib/markdown.ts'

// article.html?id=<article id>[&lang=<tag>]
// Listed on articles.html. Each language is an independent text, shown exactly as stored.

async function load(id: string, requestedLang: string | null): Promise<void> {
  const article = await get('articles', id)
  if (!article) {
    fill(document.body, { notFound: '1', notFoundMessage: `No article with ID '${id}' is published. It may be a draft, or have been removed.` })
    return
  }

  const [allTexts, target] = await Promise.all([textsOf(article.id), resolveRefs([article])])
  const texts = allTexts.sort((a, b) => a.language.localeCompare(b.language))
  const lang = pickLanguage(texts.map((t) => t.language), requestedLang)
  const text = texts.find((t) => t.language === lang)

  if (text) document.title = `${text.title} - RaceDB.org`
  fill(document.body, {
    loaded: '1',
    title: text?.title ?? 'Article',
    ...aboutFields(article, target(article)),
    kind: article.kind,
    published: formatDate(article.published_at),
    noTexts: text ? '' : '1',
    hasTexts: text ? '1' : '',
    summary: text?.summary ?? '',
    updated: text ? formatDateTime(text.updated_at) : '',
  })
  if (!text) return

  const languages = section('languages') as HTMLSelectElement
  fillList(languages, texts.map((t) => ({ language: t.language, label: `${t.language} - ${t.title}` })))
  languages.value = text.language
  languages.addEventListener('change', () => {
    location.href = articleHref(article.id, languages.value)
  })

  renderMarkdown(document.querySelector<HTMLElement>('[data-markdown]')!, text.body_md)
}

const params = new URLSearchParams(location.search)
const id = params.get('id')
try {
  if (!id) throw new Error('No article id given. Use article.html?id=<article id>.')
  await load(id, params.get('lang'))
} catch (err) {
  showError(err)
} finally {
  fill(document.body, { loading: '' })
}
