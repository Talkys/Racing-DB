import { afterAll, beforeAll, beforeEach, describe, expect, test } from 'bun:test'
import { mockSearch } from '../src/api/mock.ts'
import type { EntityPath } from '../src/api/types.ts'
import { entityArticles, latestArticles } from '../src/lib/articles.ts'

// Runs the real client against a fetch stub backed by the mock, so the requests themselves can be counted.

const MONACO = 'c39a7593-36c5-4ceb-8d07-6bc82f808381'
const MONACO_2025 = '77000000-0000-4000-8000-000000000004' // has a live blog of more than one page
const RB = '71000000-0000-4000-8000-000000000001'

const realFetch = globalThis.fetch
let requests: URL[] = []

beforeAll(() => {
  Object.defineProperty(navigator, 'languages', { value: ['en'], configurable: true }) // Bun has none
  globalThis.fetch = (async (input: string | URL | Request) => {
    const url = new URL(String(input), 'http://localhost')
    requests.push(url)
    const path = url.pathname.replace('/api/1/', '') as EntityPath
    const rows = mockSearch(path, {
      page: Number(url.searchParams.get('page') ?? 1),
      filter: url.searchParams.getAll('filter'),
      sort: url.searchParams.getAll('sort'),
    })
    return Response.json(rows)
  }) as typeof fetch
})
afterAll(() => {
  globalThis.fetch = realFetch
})
beforeEach(() => {
  requests = []
})

const paths = () => requests.map((r) => r.pathname.replace('/api/1/', ''))

describe('entityArticles', () => {
  test('filters on both entity_type and entity_id in one request', async () => {
    const cards = await entityArticles('race', MONACO)
    expect(cards.map((c) => c.fields.title)).toEqual(['Verstappen masters Monaco'])
    expect(requests[0].searchParams.getAll('filter')).toEqual(['entity_type,eq,race', `entity_id,eq,${MONACO}`])
  })

  test('follows the pages, newest first', async () => {
    const cards = await entityArticles('race', MONACO_2025)
    expect(cards.length).toBe(27)
    expect(cards[0].fields.title).toBe('Monaco 2025 live, update 27')
    expect(cards.at(-1)?.fields.title).toBe('Monaco 2025 live, update 1')
    expect(cards.every((c) => c.text && c.fields.languages === 'en')).toBe(true)
  })

  test('requests grow with pages of rows, not with cards', async () => {
    await entityArticles('race', MONACO)
    expect(paths()).toEqual(['articles', 'article-texts'])

    requests = []
    await entityArticles('race', MONACO_2025)
    // 27 articles and 27 texts: two pages of each.
    expect(paths()).toEqual(['articles', 'articles', 'article-texts', 'article-texts'])
    expect(requests[2].searchParams.get('filter')).toStartWith('article_id,in,')
  })

  test("cards without a text in the browser's language use another one", async () => {
    const [card] = await entityArticles('team', RB)
    expect(card.fields.languages).toBe('pt')
    expect(card.fields.title).toBe('A história da equipe')
  })
})

describe('latestArticles', () => {
  test('one request for the articles and one for all their texts', async () => {
    const cards = await latestArticles(3)
    expect(cards.length).toBe(3)
    expect(paths()).toEqual(['articles', 'article-texts'])
  })
})
