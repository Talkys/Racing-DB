import { search } from '../api/client.ts'
import type { Article, ArticleText } from '../api/types.ts'
import type { Fields } from './dom.ts'
import { articleHref, formatDate } from './format.ts'
import { pickLanguage } from './language.ts'

/** Drafts (no published_at) and future-dated articles behave as non-existent. The API already hides them; this is a safety net. */
export const isVisible = (article: Article) =>
  article.published_at != null && new Date(article.published_at).getTime() <= Date.now()

export const textsOf = async (articleId: string) =>
  (await search('article-texts', { filter: `article_id,eq,${articleId}` })).items

export interface ArticleCard {
  article: Article
  /** The text in the best available language, if the article has any. */
  text: ArticleText | undefined
  /** title, summary, href, kind, published, languages */
  fields: Fields
}

async function toCard(article: Article): Promise<ArticleCard> {
  const texts = await textsOf(article.id)
  const lang = pickLanguage(texts.map((t) => t.language))
  const text = texts.find((t) => t.language === lang)
  return {
    article,
    text,
    fields: {
      title: text?.title ?? '(untitled)',
      summary: text?.summary ?? '',
      href: articleHref(article.id, lang ?? undefined),
      kind: article.kind,
      published: formatDate(article.published_at),
      languages: texts.map((t) => t.language).join(', ') || 'none',
    },
  }
}

const byNewest = (a: Article, b: Article) => (b.published_at ?? '').localeCompare(a.published_at ?? '')

/** Visible articles about one entity. Category slugs are not unique across types, hence the entity_type check. */
export async function entityArticles(entityType: Article['entity_type'], entityId: string): Promise<ArticleCard[]> {
  const { items } = await search('articles', { filter: `entity_id,eq,${entityId}` })
  return Promise.all(items.filter((a) => a.entity_type === entityType && isVisible(a)).sort(byNewest).map(toCard))
}

/** Most recently published articles. */
export async function latestArticles(limit: number): Promise<ArticleCard[]> {
  const { items } = await search('articles', { sort: 'published_at,desc' })
  return Promise.all(items.filter(isVisible).slice(0, limit).map(toCard))
}
