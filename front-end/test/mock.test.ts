import { describe, expect, test } from 'bun:test'
import { mockGet, mockSearch, mockSearchArticles } from '../src/api/mock.ts'
import { mockData } from '../src/api/mock-data.ts'

const PAGE_SIZE = 25

// Well-known fixture rows (same ids as back-end/scripts/smoke.sh).
const MONACO = 'c39a7593-36c5-4ceb-8d07-6bc82f808381'
const VER = 'b11a43a8-4e36-47a8-8e68-301fc1cae6ef'
const BAG = 'a22b54b9-5f47-48b9-9f79-412fd2dbf7f0'
const ARTICLE_MONACO = '75000000-0000-4000-8000-000000000001'
const ARTICLE_DRAFT = '75000000-0000-4000-8000-000000000003'
const ARTICLE_FUTURE = '75000000-0000-4000-8000-000000000006'

const ids = (rows: { id: string }[]) => rows.map((r) => r.id)

describe('mockSearch filters', () => {
  test('eq', () => {
    const rows = mockSearch('races', { filter: 'category_id,eq,f1' })
    expect(rows.length).toBeGreaterThan(0)
    expect(rows.every((r) => r.category_id === 'f1')).toBe(true)
    expect(ids(rows)).toContain(MONACO)
  })

  test('neq', () => {
    const rows = mockSearch('races', { filter: 'category_id,neq,f1' })
    expect(rows.length).toBeGreaterThan(0)
    expect(rows.some((r) => r.category_id === 'f1')).toBe(false)
  })

  test('numeric comparisons', () => {
    expect(mockSearch('races', { filter: 'season,gte,2026' }).every((r) => r.season >= 2026)).toBe(true)
    expect(mockSearch('races', { filter: 'season,lt,2026' }).every((r) => r.season < 2026)).toBe(true)
    expect(mockSearch('circuits', { filter: 'length_km,gt,10' }).map((c) => c.name)).toEqual(['Circuit de la Sarthe'])
  })

  test('ilike is a case-insensitive substring match', () => {
    const rows = mockSearch('drivers', { filter: 'last_name,ilike,VER' })
    expect(rows.map((d) => d.code)).toContain('VER')
    expect(rows.every((d) => d.last_name.toLowerCase().includes('ver'))).toBe(true)
  })

  test('in takes |-separated values', () => {
    expect(ids(mockSearch('drivers', { filter: `id,in,${VER}|${BAG}` })).sort()).toEqual([BAG, VER].sort())
  })

  test('a value may contain commas', () => {
    expect(mockSearch('drivers', { filter: 'last_name,eq,a,b' })).toEqual([])
  })

  test('an unknown operator throws', () => {
    expect(() => mockSearch('races', { filter: 'season,like,2026' })).toThrow("invalid operator 'like'")
  })
})

describe('mockSearch ordering and paging', () => {
  test('default order follows the back end (races: newest first)', () => {
    const dates = mockSearch('races', {}).map((r) => r.race_date)
    expect(dates).toEqual([...dates].sort().reverse())
  })

  test('default order with several columns (drivers: last name, first name)', () => {
    const rows = mockSearch('drivers', {})
    const keys = rows.map((d) => `${d.last_name}\u0000${d.first_name}`)
    expect(keys).toEqual([...keys].sort((a, b) => a.localeCompare(b)))
  })

  test('explicit sort, both directions, with id as tiebreaker', () => {
    const asc = mockSearch('drivers', { sort: 'last_name,asc' }).map((d) => d.last_name)
    expect(asc).toEqual([...asc].sort((a, b) => a.localeCompare(b)))
    const desc = mockSearch('drivers', { sort: 'last_name,DESC' }).map((d) => d.last_name)
    expect(desc).toEqual([...desc].sort((a, b) => b.localeCompare(a)))
  })

  test('numbers sort numerically', () => {
    const lengths = mockSearch('circuits', { sort: 'length_km,desc' }).map((c) => c.length_km)
    expect(lengths).toEqual([...lengths].sort((a, b) => b - a))
  })

  test('pages hold PAGE_SIZE rows plus one look-ahead row', () => {
    const total = mockData.drivers.length
    expect(total).toBeGreaterThan(PAGE_SIZE)
    const page1 = mockSearch('drivers', { page: 1 })
    const page2 = mockSearch('drivers', { page: 2 })
    expect(page1.length).toBe(PAGE_SIZE + 1)
    expect(page2[0]).toEqual(page1[PAGE_SIZE]) // the look-ahead row opens the next page
    expect(new Set([...ids(page1.slice(0, PAGE_SIZE)), ...ids(page2)]).size).toBe(total)
  })

  test('a page past the end is empty', () => {
    expect(mockSearch('categories', { page: 99 })).toEqual([])
  })
})

describe('mock article visibility', () => {
  test('drafts and future-dated articles are never listed', () => {
    const listed = ids(mockSearch('articles', {}))
    expect(listed).toContain(ARTICLE_MONACO)
    expect(listed).not.toContain(ARTICLE_DRAFT)
    expect(listed).not.toContain(ARTICLE_FUTURE)
  })

  test('nor fetched by id', () => {
    expect(mockGet('articles', ARTICLE_MONACO)?.id).toBe(ARTICLE_MONACO)
    expect(mockGet('articles', ARTICLE_DRAFT)).toBeNull()
    expect(mockGet('articles', ARTICLE_FUTURE)).toBeNull()
  })

  test('their texts are hidden too', () => {
    expect(mockSearch('article-texts', { filter: `article_id,eq,${ARTICLE_DRAFT}` })).toEqual([])
    expect(mockSearch('article-texts', { filter: `article_id,eq,${ARTICLE_FUTURE}` })).toEqual([])
    expect(mockSearch('article-texts', { filter: `article_id,eq,${ARTICLE_MONACO}` }).map((t) => t.language)).toEqual(['en', 'pt-br'])
  })
})

describe('mockGet', () => {
  test('finds by id, null when missing', () => {
    expect(mockGet('races', MONACO)?.name).toBe('Monaco Grand Prix')
    expect(mockGet('categories', 'f1')?.name).toBe('Formula 1')
    expect(mockGet('races', '00000000-0000-4000-8000-000000000099')).toBeNull()
  })
})

describe('mockSearchArticles', () => {
  test('without a query, lists visible articles newest first', () => {
    const results = mockSearchArticles('', 'en', 1)
    expect(ids(results)).not.toContain(ARTICLE_DRAFT)
    expect(ids(results)).not.toContain(ARTICLE_FUTURE)
    const dates = results.map((r) => String(r.published_at))
    expect(dates).toEqual([...dates].sort().reverse())
  })

  test('matches whole words in any field, case-insensitively', () => {
    expect(ids(mockSearchArticles('POLE', 'en', 1))).toEqual([ARTICLE_MONACO])
    expect(mockSearchArticles('pol', 'en', 1)).toEqual([])
  })

  test('every query word must match', () => {
    expect(ids(mockSearchArticles('pole stop', 'en', 1))).toEqual([ARTICLE_MONACO])
    expect(mockSearchArticles('pole nonsense', 'en', 1)).toEqual([])
  })

  test('prefers the requested language, then the same primary subtag', () => {
    expect(mockSearchArticles('verstappen', 'pt-br', 1).find((r) => r.id === ARTICLE_MONACO)?.language).toBe('pt-br')
    expect(mockSearchArticles('verstappen', 'pt', 1).find((r) => r.id === ARTICLE_MONACO)?.language).toBe('pt-br')
    expect(mockSearchArticles('verstappen', 'en', 1).find((r) => r.id === ARTICLE_MONACO)?.language).toBe('en')
  })

  test('marks matched words in the preview and strips Markdown', () => {
    const [result] = mockSearchArticles('controlled', 'en', 1)
    expect(result.preview).toContainEqual({ text: 'controlled', match: true })
    expect(result.preview.map((s) => s.text).join('')).not.toContain('**')
  })

  test('drafts and future-dated articles never match', () => {
    expect(mockSearchArticles('draft', 'en', 1)).toEqual([])
    expect(mockSearchArticles('scheduled', 'en', 1)).toEqual([])
  })
})
