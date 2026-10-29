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
  /** `column,operator,value`, e.g. `language,eq,en` (operators: eq neq gt gte lt lte ilike in; `in` values are `|` separated). */
  filter?: string
  /** `column,direction`, e.g. `race_date,desc`. */
  sort?: string
}

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
  if (query.filter) params.set('filter', query.filter)
  if (query.sort) params.set('sort', query.sort)
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

/**
 * Fetches the given ids of one entity in a single request (`id,in,a|b|c`), keyed by id.
 * Missing ids are simply absent from the map. Meant for the related rows of one page, so at most 100 ids.
 */
export async function getMany<P extends EntityPath>(path: P, ids: (string | null | undefined)[]): Promise<Map<string, EntityMap[P]>> {
  const unique = [...new Set(ids.filter((id): id is string => !!id))]
  if (unique.length === 0) return new Map()
  const rows: EntityMap[P][] = []
  // Pages hold PAGE_SIZE rows, so a long id list may need more than one.
  for (let page = 1; ; page++) {
    const result = await search(path, { filter: `id,in,${unique.join('|')}`, page })
    rows.push(...result.items)
    if (!result.hasNext) break
  }
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
