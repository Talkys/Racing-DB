import type { Query } from './client.ts'
import type { Article, ArticleSearchResult, ArticleText, EntityMap, EntityPath, PreviewSegment } from './types.ts'
import { mockData } from './mock-data.ts'

const PAGE_SIZE = 25

type Row = Record<string, unknown>
type Order = [column: string, direction: 'asc' | 'desc'][]

const OPS: Record<string, (a: unknown, b: string) => boolean> = {
  eq: (a, b) => String(a) === b,
  neq: (a, b) => String(a) !== b,
  gt: (a, b) => Number(a) > Number(b),
  gte: (a, b) => Number(a) >= Number(b),
  lt: (a, b) => Number(a) < Number(b),
  lte: (a, b) => Number(a) <= Number(b),
  ilike: (a, b) => String(a).toLowerCase().includes(b.toLowerCase()),
  in: (a, b) => b.split('|').includes(String(a)),
}

/** Each entity's DefaultOrder in back-end/internal/<entity>/repository.go. */
const DEFAULT_ORDER: Record<EntityPath, Order> = {
  categories: [['name', 'asc'], ['id', 'asc']],
  circuits: [['name', 'asc'], ['id', 'asc']],
  drivers: [['last_name', 'asc'], ['first_name', 'asc'], ['id', 'asc']],
  races: [['race_date', 'desc'], ['id', 'asc']],
  teams: [['name', 'asc'], ['id', 'asc']],
  engines: [['manufacturer', 'asc'], ['name', 'asc'], ['id', 'asc']],
  chassis: [['manufacturer', 'asc'], ['name', 'asc'], ['id', 'asc']],
  cars: [['season', 'desc'], ['name', 'asc'], ['id', 'asc']],
  articles: [['published_at', 'desc'], ['id', 'asc']],
  'article-texts': [['language', 'asc'], ['id', 'asc']],
}

/** Like article.VisibleCondition: drafts and future-dated articles (and their texts) do not exist for the API. */
const isPublished = (row: Row) => row.published_at != null && new Date(String(row.published_at)).getTime() <= Date.now()
const visibleArticleIds = () => new Set((mockData.articles as unknown as Row[]).filter(isPublished).map((a) => a.id))

const BASE_CONDITION: Partial<Record<EntityPath, (row: Row) => boolean>> = {
  articles: isPublished,
  'article-texts': (row) => visibleArticleIds().has(row.article_id),
}

function compare(a: unknown, b: unknown): number {
  if (typeof a === 'number' && typeof b === 'number') return a - b
  return String(a).localeCompare(String(b))
}

function rowsOf(path: EntityPath): Row[] {
  const rows = mockData[path] as unknown as Row[]
  const base = BASE_CONDITION[path]
  return base ? rows.filter(base) : [...rows]
}

/**
 * Mimics the real search endpoints (utils/queries/querybuilder.go): optional filter, optional sort with
 * `id` as tiebreaker (else the entity's default order), and pageSize + 1 rows per page.
 */
export function mockSearch<P extends EntityPath>(path: P, query: Query): EntityMap[P][] {
  let rows = rowsOf(path)

  if (query.filter) {
    const [col, op, ...rest] = query.filter.split(',')
    const test = OPS[op]
    if (!test) throw new Error(`mock: invalid operator '${op}'`)
    const value = rest.join(',')
    rows = rows.filter((row) => test(row[col], value))
  }

  const [sortCol, dir] = (query.sort ?? '').split(',')
  const order: Order = sortCol ? [[sortCol, dir?.toLowerCase() === 'desc' ? 'desc' : 'asc'], ['id', 'asc']] : DEFAULT_ORDER[path]
  rows.sort((a, b) => {
    for (const [col, direction] of order) {
      const diff = compare(a[col], b[col])
      if (diff) return direction === 'desc' ? -diff : diff
    }
    return 0
  })

  const offset = ((query.page ?? 1) - 1) * PAGE_SIZE
  return rows.slice(offset, offset + PAGE_SIZE + 1) as unknown as EntityMap[P][]
}

export function mockGet<P extends EntityPath>(path: P, id: string): EntityMap[P] | null {
  return (rowsOf(path).find((row) => row.id === id) as unknown as EntityMap[P]) ?? null
}

/** Rough port of markdown_to_text in init.sql. */
const markdownToText = (md: string) =>
  md
    .replace(/<(https?:\/\/[^>]+)>/g, '$1')
    .replace(/<[^>]+>/g, ' ')
    .replace(/!?\[([^\]]*)\](\([^)]*\)|\[[^\]]*\])/g, '$1')
    .replace(/^\s*```.*$|^\s{0,3}\[[^\]]+\]:\s+\S+.*$/gm, '')
    .replace(/^\s{0,3}(#{1,6}\s+|>\s?)|^\s*([*+-]|\d+[.)])\s+/gm, '')
    .replace(/[*~`|]+/g, ' ')
    .replace(/(^|[^\p{L}\p{N}])_+|_+([^\p{L}\p{N}]|$)/gu, '$1$2')
    .replace(/\s+/g, ' ')
    .trim()

const languageRank = (language: string, wanted: string) =>
  language === wanted ? 0 : language.split('-')[0] === wanted.split('-')[0] ? 1 : 2

/**
 * Mimics GET /articles/search (internal/articlesearch): whole-word, case-insensitive matching of every
 * query word (no phrase/or/- syntax), one text per article, previews around the first match.
 */
export function mockSearchArticles(q: string, lang: string, page: number): ArticleSearchResult[] {
  const words = q.toLowerCase().split(/[^\p{L}\p{N}]+/u).filter(Boolean)
  const texts = rowsOf('article-texts') as unknown as ArticleText[]
  const tokens = (s: string | null) => new Set((s ?? '').toLowerCase().split(/[^\p{L}\p{N}]+/u))

  // Newest first (the back end orders search results by rank first; the mock has no ranking).
  const articles = (rowsOf('articles') as unknown as Article[]).sort((a, b) => String(b.published_at).localeCompare(String(a.published_at)) || a.id.localeCompare(b.id))
  const results = articles.flatMap((article) => {
    const candidates = texts
      .filter((t) => t.article_id === article.id)
      .map((t) => ({ t, plain: markdownToText(t.body_md) }))
      .filter(({ t, plain }) => words.every((w) => tokens(t.title).has(w) || tokens(t.summary).has(w) || tokens(plain).has(w)))
      .sort((a, b) => languageRank(a.t.language, lang) - languageRank(b.t.language, lang) || a.t.language.localeCompare(b.t.language))
    if (words.length && !candidates.length) return []
    const best = candidates[0]
    return [{
      id: article.id, entity_type: article.entity_type, entity_id: article.entity_id, kind: article.kind, published_at: article.published_at,
      language: best?.t.language ?? null, title: best?.t.title ?? null, summary: best?.t.summary ?? null,
      preview: best ? mockPreview(best.plain, best.t.title, words) : [],
    }]
  })

  const offset = (page - 1) * PAGE_SIZE
  return results.slice(offset, offset + PAGE_SIZE + 1)
}

function mockPreview(plain: string, title: string, words: string[]): PreviewSegment[] {
  let text = plain.toLowerCase().startsWith(title.toLowerCase()) ? plain.slice(title.length).replace(/^[\s\p{P}]+/u, '') || plain : plain
  if (words.length) {
    const at = Math.max(0, text.toLowerCase().indexOf(words[0]) - 80)
    text = (at > 0 ? '… ' : '') + text.slice(at, at + 240)
  } else if (text.length > 240) {
    text = text.slice(0, 240).replace(/\s+\S*$/, '') + ' …'
  }
  const segments: PreviewSegment[] = []
  for (const part of text.split(/([\p{L}\p{N}]+)/u)) {
    if (!part) continue
    const match = words.includes(part.toLowerCase())
    const last = segments.at(-1)
    if (last && !last.match && !match) last.text += part
    else segments.push({ text: part, match })
  }
  return segments
}
