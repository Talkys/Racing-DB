import { get } from '../api/client.ts'
import { isVisible, textsOf } from '../lib/articles.ts'
import { fill, fillList, section, showError } from '../lib/dom.ts'
import { articleHref, entityHref, formatDate, formatDateTime } from '../lib/format.ts'
import { pickLanguage } from '../lib/language.ts'
import { renderMarkdown } from '../lib/markdown.ts'

// article.html?id=<article id>[&lang=<tag>]
// Articles have no index page. Each language is an independent text, shown exactly as stored.

async function load(id: string, requestedLang: string | null): Promise<void> {
  const article = await get('articles', id)
  if (!article || !isVisible(article)) throw new Error(`Article with ID '${id}' does not exist.`)

  const texts = (await textsOf(article.id)).sort((a, b) => a.language.localeCompare(b.language))
  const lang = pickLanguage(texts.map((t) => t.language), requestedLang)
  const text = texts.find((t) => t.language === lang)
  const entityPage = entityHref(article.entity_type, article.entity_id)

  if (text) document.title = `${text.title} - RaceDB.org`
  fill(document.body, {
    loaded: '1',
    title: text?.title ?? 'Article',
    entityType: article.entity_type,
    entityHref: entityPage ?? '',
    entityPlain: entityPage ? '' : article.entity_type,
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
