import { ApiError, type Query } from './client.ts'
import { mockData } from './mock-data.ts'
import type { Article, ArticleSearchResult, ArticleText, EntityMap, EntityPath, PreviewSegment } from './types.ts'

// Mock mode (VITE_API_MOCK=true) is for demos and working on the UI without a back end. It follows the API
// closely enough for every page to work, but the back end is the reference: the Go integration tests cover it,
// and this file is only re-checked against it when the API changes.

const PAGE_SIZE = 25
const MAX_FILTERS = 10
const MAX_IN_VALUES = 100

type Row = Record<string, unknown>
type Order = [column: string, direction: 'asc' | 'desc'][]

/** queries.Kind; `slug` is a text id that also takes `in` (categories). */
type Kind = 'text' | 'int' | 'decimal' | 'bool' | 'uuid' | 'slug'

/** The operators each kind accepts, as in querybuilder.go. */
const OPS_BY_KIND: Record<Kind, string[]> = {
  text: ['eq', 'neq', 'ilike'],
  int: ['eq', 'gt', 'gte', 'lt', 'lte'],
  decimal: ['eq', 'gt', 'gte', 'lt', 'lte'],
  bool: ['eq'],
  uuid: ['eq', 'in'],
  slug: ['eq', 'in'],
}

const UUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
const BOOLS: Record<string, boolean> = { 1: true, t: true, T: true, true: true, TRUE: true, True: true, 0: false, f: false, F: false, false: false, FALSE: false, False: false }

/** Parses a filter value like parseFilterValue in querybuilder.go; undefined means it is invalid. */
function parseValue(kind: Kind, raw: string): unknown {
  switch (kind) {
    case 'int':
      return /^[+-]?\d+$/.test(raw) ? Number(raw) : undefined
    case 'decimal':
      return raw.trim() !== '' && Number.isFinite(Number(raw)) ? Number(raw) : undefined
    case 'bool':
      return BOOLS[raw]
    case 'uuid':
      return UUID.test(raw) ? raw.toLowerCase() : undefined
    default:
      return raw
  }
}

/** SQL comparisons: a NULL column never matches. */
const OPS: Record<string, (a: unknown, b: unknown) => boolean> = {
  eq: (a, b) => a === b,
  neq: (a, b) => a !== b,
  gt: (a, b) => (a as number) > (b as number),
  gte: (a, b) => (a as number) >= (b as number),
  lt: (a, b) => (a as number) < (b as number),
  lte: (a, b) => (a as number) <= (b as number),
  ilike: (a, b) => String(a).toLowerCase().includes(String(b).toLowerCase()),
  in: (a, b) => (b as unknown[]).includes(a),
}

/** Each entity's AllowedFilters in back-end/internal/<entity>/<entity>.go. */
const FILTERS: Record<EntityPath, Record<string, Kind>> = {
  categories: { id: 'slug', name: 'text', description: 'text', is_active: 'bool' },
  circuits: { id: 'uuid', name: 'text', location: 'text', country: 'text', length_km: 'decimal' },
  drivers: { id: 'uuid', first_name: 'text', last_name: 'text', code: 'text', nationality: 'text' },
  races: { id: 'uuid', category_id: 'text', circuit_id: 'uuid', winner_id: 'uuid', name: 'text', season: 'int', round: 'int' },
  teams: { id: 'uuid', name: 'text', short_name: 'text', country: 'text', founded_year: 'int' },
  engines: {
    id: 'uuid', manufacturer: 'text', name: 'text', configuration: 'text', displacement_cc: 'int', cylinders: 'int', aspiration: 'text', power_kw: 'decimal',
  },
  chassis: { id: 'uuid', manufacturer: 'text', name: 'text', material: 'text', wheelbase_mm: 'int', weight_kg: 'decimal' },
  cars: { id: 'uuid', category_id: 'text', team_id: 'uuid', engine_id: 'uuid', chassis_id: 'uuid', name: 'text', season: 'int' },
  articles: { id: 'uuid', entity_type: 'text', entity_id: 'text', kind: 'text' },
  'article-texts': { id: 'uuid', article_id: 'uuid', language: 'text', title: 'text' },
}

/** Each entity's AllowedSorts. */
const SORTS: Record<EntityPath, string[]> = {
  categories: ['id', 'name', 'is_active'],
  circuits: ['name', 'location', 'country', 'length_km'],
  drivers: ['last_name', 'first_name', 'code', 'nationality'],
  races: ['race_date', 'name', 'category_id', 'season', 'round'],
  teams: ['name', 'short_name', 'country', 'founded_year'],
  engines: ['manufacturer', 'name', 'configuration', 'displacement_cc', 'cylinders', 'aspiration', 'power_kw'],
  chassis: ['manufacturer', 'name', 'material', 'wheelbase_mm', 'weight_kg'],
  cars: ['name', 'season', 'category_id'],
  articles: ['published_at', 'created_at', 'entity_type', 'kind'],
  'article-texts': ['title', 'language'],
}

/** Each entity's DefaultOrder. */
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

const badRequest = (message: string) => new ApiError(400, message)

/** Like article.VisibleCondition: drafts and future-dated articles (and their texts) do not exist for the API. */
const isPublished = (row: Row) => row.published_at != null && new Date(String(row.published_at)).getTime() <= Date.now()
const visibleArticleIds = () => new Set((mockData.articles as unknown as Row[]).filter(isPublished).map((a) => a.id))

const BASE_CONDITION: Partial<Record<EntityPath, (row: Row) => boolean>> = {
  articles: isPublished,
  'article-texts': (row) => visibleArticleIds().has(row.article_id),
}

/** Like Postgres: NULLs sort after every value in ascending order (and so first in descending order). */
function compare(a: unknown, b: unknown): number {
  if (a == null || b == null) return a == null ? (b == null ? 0 : 1) : -1
  if (typeof a === 'number' && typeof b === 'number') return a - b
  if (typeof a === 'boolean' && typeof b === 'boolean') return Number(a) - Number(b)
  return String(a).localeCompare(String(b))
}

function rowsOf(path: EntityPath): Row[] {
  const rows = mockData[path] as unknown as Row[]
  const base = BASE_CONDITION[path]
  return base ? rows.filter(base) : [...rows]
}

/** One row test per filter, rejecting what BuildSearchQuery rejects. */
function parseFilters(path: EntityPath, filters: string[]): ((row: Row) => boolean)[] {
  if (filters.length > MAX_FILTERS) throw badRequest(`too many filters: ${filters.length} (at most ${MAX_FILTERS})`)
  return filters.map((filter) => {
    const [col, op, ...rest] = filter.split(',')
    if (op === undefined || rest.length === 0) throw badRequest(`invalid filter: '${filter}' (expected column,operator,value)`)
    const raw = rest.join(',')
    const kind = FILTERS[path][col]
    if (!kind) throw badRequest(`invalid filter column: ${col}`)
    if (!OPS_BY_KIND[kind].includes(op)) throw badRequest(`invalid operator '${op}' for column '${col}'`)
    const parts = op === 'in' ? raw.split('|') : [raw]
    const values = parts.map((part) => parseValue(kind, part))
    if (parts.length > MAX_IN_VALUES || values.includes(undefined)) throw badRequest(`invalid value '${raw}' for column '${col}'`)
    const value = op === 'in' ? values : values[0]
    const test = OPS[op]
    return (row) => row[col] != null && test(row[col], value)
  })
}

/** Reads `sort` like BuildCriteriaFromURL and buildOrderBy: columns in priority order, each optionally followed by its direction. */
function parseOrder(path: EntityPath, sorts: string[]): Order {
  const order: [string, 'asc' | 'desc' | ''][] = []
  for (const token of sorts.flatMap((s) => s.split(','))) {
    const lower = token.toLowerCase()
    if (!token) continue
    if ((lower === 'asc' || lower === 'desc') && order.length > 0 && order[order.length - 1][1] === '') {
      order[order.length - 1][1] = lower
      continue
    }
    if (!SORTS[path].includes(token)) throw badRequest(`invalid sort column or direction: ${token}`)
    if (order.some(([col]) => col === token)) throw badRequest(`duplicate sort column: ${token}`)
    order.push([token, ''])
  }
  if (order.length === 0) return DEFAULT_ORDER[path]
  const result: Order = order.map(([col, dir]) => [col, dir || 'asc'])
  return order.some(([col]) => col === 'id') ? result : [...result, ['id', 'asc']]
}

/**
 * Mimics the real search endpoints (utils/queries/querybuilder.go): filters AND-ed together, sorts with `id` as
 * tiebreaker (else the entity's default order), pageSize + 1 rows per page, and an ApiError(400) for anything the
 * back end would reject with 400.
 */
export function mockSearch<P extends EntityPath>(path: P, query: Query): EntityMap[P][] {
  const tests = parseFilters(path, [query.filter ?? []].flat())
  const order = parseOrder(path, [query.sort ?? []].flat())

  const rows = rowsOf(path).filter((row) => tests.every((test) => test(row)))
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
 * Approximates GET /articles/search (internal/articlesearch): whole-word, case-insensitive matching of every
 * query word, one text per article, previews around the first match. It is not websearch_to_tsquery: no
 * "phrases", `or` or `-word`, no ranking, and words are split on punctuation more loosely than Postgres does.
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
    text = `${at > 0 ? '… ' : ''}${text.slice(at, at + 240)}`
  } else if (text.length > 240) {
    text = `${text.slice(0, 240).replace(/\s+\S*$/, '')} …`
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
