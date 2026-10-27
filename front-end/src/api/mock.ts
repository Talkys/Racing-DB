import type { Query } from './client.ts'
import type { EntityMap, EntityPath } from './types.ts'
import { mockData } from './mock-data.ts'

const PAGE_SIZE = 25

type Row = Record<string, unknown>

const OPS: Record<string, (a: unknown, b: string) => boolean> = {
  eq: (a, b) => String(a) === b,
  neq: (a, b) => String(a) !== b,
  gt: (a, b) => Number(a) > Number(b),
  gte: (a, b) => Number(a) >= Number(b),
  lt: (a, b) => Number(a) < Number(b),
  lte: (a, b) => Number(a) <= Number(b),
  ilike: (a, b) => String(a).toLowerCase().includes(b.toLowerCase()),
}

function compare(a: unknown, b: unknown): number {
  if (typeof a === 'number' && typeof b === 'number') return a - b
  return String(a).localeCompare(String(b))
}

/**
 * Mimics the real search endpoints (utils/queries/querybuilder.go): optional
 * filter, optional sort with `id` as tiebreaker, and pageSize + 1 rows per page.
 */
export function mockSearch<P extends EntityPath>(path: P, query: Query): EntityMap[P][] {
  let rows = [...(mockData[path] as unknown as Row[])]

  if (query.filter) {
    const [col, op, ...rest] = query.filter.split(',')
    const test = OPS[op]
    if (!test) throw new Error(`mock: invalid operator '${op}'`)
    const value = rest.join(',')
    rows = rows.filter((row) => test(row[col], value))
  }

  const [sortCol = 'id', dir = 'asc'] = (query.sort ?? '').split(',')
  const sign = dir.toLowerCase() === 'desc' ? -1 : 1
  rows.sort((a, b) => sign * (compare(a[sortCol], b[sortCol]) || compare(a.id, b.id)))

  const offset = ((query.page ?? 1) - 1) * PAGE_SIZE
  return rows.slice(offset, offset + PAGE_SIZE + 1) as unknown as EntityMap[P][]
}

export function mockGet<P extends EntityPath>(path: P, id: string): EntityMap[P] | null {
  return ((mockData[path] as unknown as Row[]).find((row) => row.id === id) as unknown as EntityMap[P]) ?? null
}
