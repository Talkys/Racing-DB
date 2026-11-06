import type { ArticleSearchResult, EntityMap, EntityPath } from './types.ts'

const BASE_URL = '/api/1'
const TIMEOUT_MS = 3000

/** The back end returns pageSize + 1 rows; the extra row only signals that another page exists. */
export const PAGE_SIZE = 25

/** Set VITE_API_MOCK=true to skip the network entirely and use the in-memory mock. There is no automatic fallback. */
export const USE_MOCK = import.meta.env.VITE_API_MOCK === 'true'

/** Every failed request: the back end's own errors (400, 404, 500) and network failures (status 0). */
export class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

export interface Query {
  page?: number
  /**
   * `column,operator,value`, e.g. `language,eq,en` (operators: eq neq gt gte lt lte ilike in; `in` values are `|`
   * separated). Several filters are AND-ed.
   */
  filter?: string | string[]
  /** Columns in priority order, each optionally followed by its direction: `season,desc,name` or `['season,desc', 'name']`. */
  sort?: string | string[]
}

/** The back end caps `in` lists at this many values. */
const MAX_IN_VALUES = 100

export interface Page<T> {
  items: T[]
  hasNext: boolean
}

function toPage<T>(rows: T[]): Page<T> {
  return { items: rows.slice(0, PAGE_SIZE), hasNext: rows.length > PAGE_SIZE }
}

function buildUrl(path: EntityPath, query: Query): string {
  const params = new URLSearchParams()
  if (query.page) params.set('page', String(query.page))
  for (const filter of [query.filter ?? []].flat()) params.append('filter', filter)
  for (const sort of [query.sort ?? []].flat()) params.append('sort', sort)
  const qs = params.toString()
  return `${BASE_URL}/${path}${qs ? `?${qs}` : ''}`
}

/** Error bodies are `{"status": 400, "error": "..."}`; anything else (e.g. a proxy page) falls back to the status text. */
async function errorFrom(res: Response): Promise<ApiError> {
  const body = await res.json().catch(() => null)
  const message = typeof body?.error === 'string' ? body.error : `${res.status} ${res.statusText || 'request failed'}`
  return new ApiError(res.status, message)
}

async function request(url: string): Promise<Response> {
  try {
    return await fetch(url, { signal: AbortSignal.timeout(TIMEOUT_MS) })
  } catch (err) {
    throw new ApiError(0, `The server could not be reached (${err instanceof Error ? err.message : String(err)}).`)
  }
}

/** GET /api/1/<path>. */
export async function search<P extends EntityPath>(path: P, query: Query = {}): Promise<Page<EntityMap[P]>> {
  if (USE_MOCK) {
    const { mockSearch } = await import('./mock.ts')
    return toPage(mockSearch(path, query))
  }
  const res = await request(buildUrl(path, query))
  if (!res.ok) throw await errorFrom(res)
  return toPage<EntityMap[P]>(await res.json())
}

/** GET /api/1/<path>/<id>. Resolves to null when the entity does not exist. */
export async function get<P extends EntityPath>(path: P, id: string): Promise<EntityMap[P] | null> {
  if (USE_MOCK) {
    const { mockGet } = await import('./mock.ts')
    return mockGet(path, id)
  }
  const res = await request(`${BASE_URL}/${path}/${encodeURIComponent(id)}`)
  if (res.status === 404) return null
  if (!res.ok) throw await errorFrom(res)
  return (await res.json()) as EntityMap[P]
}

/** Every row matching the query, following the pages. Meant for lists known to stay small. */
export async function searchAll<P extends EntityPath>(path: P, query: Omit<Query, 'page'> = {}): Promise<EntityMap[P][]> {
  const rows: EntityMap[P][] = []
  for (let page = 1; ; page++) {
    const result = await search(path, { ...query, page })
    rows.push(...result.items)
    if (!result.hasNext) return rows
  }
}

/**
 * Every row whose `column` is one of `values` (`column,in,a|b|c`), in as few requests as the back end allows:
 * one per 100 values, plus one per extra page of results.
 */
export async function searchIn<P extends EntityPath>(path: P, column: string, values: (string | null | undefined)[]): Promise<EntityMap[P][]> {
  const unique = [...new Set(values.filter((v): v is string => !!v))]
  const chunks: string[][] = []
  for (let i = 0; i < unique.length; i += MAX_IN_VALUES) chunks.push(unique.slice(i, i + MAX_IN_VALUES))
  const results = await Promise.all(chunks.map((chunk) => searchAll(path, { filter: `${column},in,${chunk.join('|')}` })))
  return results.flat()
}

/** Fetches the given ids of one entity, keyed by id. Missing ids are simply absent from the map. */
export async function getMany<P extends EntityPath>(path: P, ids: (string | null | undefined)[]): Promise<Map<string, EntityMap[P]>> {
  const rows = await searchIn(path, 'id', ids)
  return new Map(rows.map((row) => [(row as { id: string }).id, row]))
}

/**
 * GET /api/1/articles/search. Whole-word search over article texts; an empty `q` lists every article
 * newest first. `lang` is the preferred text language.
 */
export async function searchArticles(q: string, lang: string, page = 1): Promise<Page<ArticleSearchResult>> {
  if (USE_MOCK) {
    const { mockSearchArticles } = await import('./mock.ts')
    return toPage(mockSearchArticles(q, lang, page))
  }
  const params = new URLSearchParams({ lang })
  if (q) params.set('q', q)
  if (page > 1) params.set('page', String(page))
  const res = await request(`${BASE_URL}/articles/search?${params}`)
  if (!res.ok) throw await errorFrom(res)
  return toPage<ArticleSearchResult>(await res.json())
}
