import { search, searchAll, searchIn } from '../api/client.ts'
import type { Article, ArticleText } from '../api/types.ts'
import type { Fields } from './dom.ts'
import { articleHref, formatDate } from './format.ts'
import { pickLanguage } from './language.ts'

export const textsOf = (articleId: string) => searchAll('article-texts', { filter: `article_id,eq,${articleId}` })

export interface ArticleCard {
  article: Article
  /** The text in the best available language, if the article has any. */
  text: ArticleText | undefined
  /** title, summary, href, kind, published, languages */
  fields: Fields
}

/** Cards for the given articles, in their order. All their texts come in one batch, not one request per card. */
async function toCards(articles: Article[]): Promise<ArticleCard[]> {
  const texts = new Map<string, ArticleText[]>()
  for (const t of await searchIn('article-texts', 'article_id', articles.map((a) => a.id))) {
    texts.set(t.article_id, [...(texts.get(t.article_id) ?? []), t])
  }
  return articles.map((article) => {
    const own = texts.get(article.id) ?? []
    const lang = pickLanguage(own.map((t) => t.language))
    const text = own.find((t) => t.language === lang)
    return {
      article,
      text,
      fields: {
        title: text?.title ?? '(untitled)',
        summary: text?.summary ?? '',
        href: articleHref(article.id, lang ?? undefined),
        kind: article.kind,
        published: formatDate(article.published_at),
        languages: own.map((t) => t.language).join(', ') || 'none',
      },
    }
  })
}

/** Every published article about one entity, newest first (the API hides drafts and future-dated ones). Category slugs are not unique across types, hence the entity_type filter. */
export async function entityArticles(entityType: Article['entity_type'], entityId: string): Promise<ArticleCard[]> {
  return toCards(await searchAll('articles', { filter: [`entity_type,eq,${entityType}`, `entity_id,eq,${entityId}`], sort: 'published_at,desc' }))
}

/** Most recently published articles. */
export async function latestArticles(limit: number): Promise<ArticleCard[]> {
  const { items } = await search('articles', { sort: 'published_at,desc' })
  return toCards(items.slice(0, limit))
}
