import type { EntityMap, EntityPath } from './types.ts'

const BASE_URL = '/api/1'
const TIMEOUT_MS = 3000

/** The back end returns pageSize + 1 rows; the extra row only signals that another page exists. */
export const PAGE_SIZE = 25

/** Set VITE_API_MOCK=true to skip the network entirely and always use mock data. */
const FORCE_MOCK = import.meta.env.VITE_API_MOCK === 'true'

/** The back end answered, but with an error (e.g. 500 for an invalid filter). Never masked by mock data. */
export class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

export interface Query {
  page?: number
  /** `column,operator,value`, e.g. `language,eq,en` (operators: eq neq gt gte lt lte ilike). */
  filter?: string
  /** `column,direction`, e.g. `race_date,desc`. */
  sort?: string
}

export interface Page<T> {
  items: T[]
  hasNext: boolean
}

function toPage<T>(rows: T[] | null): Page<T> {
  const all = rows ?? [] // Go encodes an empty result as `null`
  return { items: all.slice(0, PAGE_SIZE), hasNext: all.length > PAGE_SIZE }
}

function buildUrl(path: EntityPath, query: Query): string {
  const params = new URLSearchParams()
  if (query.page) params.set('page', String(query.page))
  if (query.filter) params.set('filter', query.filter)
  if (query.sort) params.set('sort', query.sort)
  const qs = params.toString()
  return `${BASE_URL}/${path}${qs ? `?${qs}` : ''}`
}

/**
 * GET /api/1/<path>. If the back end is unreachable, falls back to the in-memory
 * mock so the front end can be developed in isolation. Errors the back end
 * reports itself are thrown as ApiError instead.
 */
export async function search<P extends EntityPath>(
  path: P,
  query: Query = {},
): Promise<Page<EntityMap[P]>> {
  if (!FORCE_MOCK) {
    try {
      const res = await fetch(buildUrl(path, query), { signal: AbortSignal.timeout(TIMEOUT_MS) })
      if (res.ok) return toPage<EntityMap[P]>(await res.json())
      const body = (await res.text()).trim()
      // The Vite dev proxy answers an unreachable target with a bodiless 5xx; Go's http.Error always has a body.
      if (body) throw new ApiError(res.status, body)
      console.warn(`[api] ${path} unreachable (${res.status}), using mock data`)
    } catch (err) {
      if (err instanceof ApiError) throw err
      console.warn(`[api] ${path} unreachable, using mock data`, err)
    }
  }
  const { mockSearch } = await import('./mock.ts')
  return toPage(mockSearch(path, query))
}

/**
 * GET /api/1/<path>/<id>. Resolves to null when the entity does not exist.
 * (Not implemented by the back end yet: until it is, the mock answers when the route is missing.)
 */
export async function get<P extends EntityPath>(path: P, id: string): Promise<EntityMap[P] | null> {
  if (!FORCE_MOCK) {
    try {
      const res = await fetch(`${BASE_URL}/${path}/${encodeURIComponent(id)}`, { signal: AbortSignal.timeout(TIMEOUT_MS) })
      if (res.ok) return (await res.json()) as EntityMap[P]
      const body = (await res.text()).trim()
      // Go's mux answers an unregistered route with "404 page not found"; a real "no such entity" 404 is the back end's own.
      if (res.status === 404 && body !== '404 page not found') return null
      if (body) throw new ApiError(res.status, body)
      console.warn(`[api] ${path}/${id} unreachable (${res.status}), using mock data`)
    } catch (err) {
      if (err instanceof ApiError) throw err
      console.warn(`[api] ${path}/${id} unreachable, using mock data`, err)
    }
  }
  const { mockGet } = await import('./mock.ts')
  return mockGet(path, id)
}
