import { afterAll, beforeAll, beforeEach, describe, expect, test } from 'bun:test'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { GlobalRegistrator } from '@happy-dom/global-registrator'
import { mockGet, mockSearch } from '../src/api/mock.ts'
import type { EntityPath } from '../src/api/types.ts'
import { type DetailPageConfig, fillRelated, initDetailPage, relatedRaces } from '../src/lib/detail-page.ts'
import { fill } from '../src/lib/dom.ts'
import { aboutFields, resolveRefs } from '../src/lib/entities.ts'

// Runs the detail-page skeleton against the real HTML and a fetch stub backed by the mock, like the back end would
// answer: 404 for an unknown id, 400 for a malformed one.

const VER = 'b11a43a8-4e36-47a8-8e68-301fc1cae6ef'
const UNKNOWN = '00000000-0000-4000-8000-000000000099'
const DRIVER_DELETED = '70000000-0000-4000-8000-000000000999' // article 7 is about it (mock-data.ts)
const UUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

const body = (file: string) => {
  const html = readFileSync(resolve(import.meta.dirname, '..', file), 'utf8')
  return html.slice(html.indexOf('<body'), html.indexOf('</body>')).replace(/<script[^>]*><\/script>/g, '')
}

beforeAll(() => {
  GlobalRegistrator.register({ url: 'http://localhost/' })
  Object.defineProperty(navigator, 'languages', { value: ['en'], configurable: true })
  globalThis.fetch = (async (input: string | URL | Request) => {
    const url = new URL(String(input), 'http://localhost')
    const [path, id] = url.pathname.replace('/api/1/', '').split('/') as [EntityPath, string?]
    if (id !== undefined) {
      if (path !== 'categories' && !UUID.test(id)) return Response.json({ status: 400, error: 'invalid id' }, { status: 400 })
      const row = mockGet(path, decodeURIComponent(id))
      return row ? Response.json(row) : Response.json({ status: 404, error: 'not found' }, { status: 404 })
    }
    return Response.json(mockSearch(path, { page: Number(url.searchParams.get('page') ?? 1), filter: url.searchParams.getAll('filter'), sort: url.searchParams.getAll('sort') }))
  }) as typeof fetch
})
afterAll(() => GlobalRegistrator.unregister())

const shown = (key: string) => !document.querySelector<HTMLElement>(`[data-show-if="${key}"]`)!.hidden
const text = (selector: string) => document.querySelector(selector)!.textContent

async function open<P extends EntityPath>(file: string, id: string, config: DetailPageConfig<P>): Promise<void> {
  document.body.innerHTML = body(file)
  history.replaceState({}, '', `/${file}?id=${encodeURIComponent(id)}`)
  await initDetailPage(config)
}

const driverPage: DetailPageConfig<'drivers'> = {
  path: 'drivers',
  noun: 'driver',
  articles: 'driver',
  async render(driver) {
    fill(document.body, { name: `${driver.first_name} ${driver.last_name}` })
    fillRelated('wins', await relatedRaces(`winner_id,eq,${driver.id}`), 'races.html')
  },
}

describe('initDetailPage', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
  })

  test('renders the entity and its articles', async () => {
    await open('driver.html', VER, driverPage)
    expect(shown('loaded')).toBe(true)
    expect(shown('loading')).toBe(false)
    expect(shown('notFound')).toBe(false)
    expect(text('[data-field="name"]')).toBe('Max Verstappen')
    expect(document.querySelectorAll('[data-section="articles"] li').length).toBeGreaterThan(0)
  })

  // Checked on the rows, not the DOM: happy-dom loses the content of a <template> inside <tbody>.
  test('related races link to their race and category pages', async () => {
    const { rows } = await relatedRaces(`winner_id,eq,${VER}`)
    expect(rows.length).toBeGreaterThan(0)
    expect(rows.every((r) => r.href?.startsWith('race.html?id=') && r.categoryHref?.startsWith('category.html?id='))).toBe(true)
  })

  test.each([
    ['unknown', UNKNOWN],
    ['malformed', 'not-a-uuid'],
  ])('an %s id shows the not-found state, not an error', async (_, id) => {
    await open('driver.html', id, driverPage)
    expect(shown('notFound')).toBe(true)
    expect(shown('loaded')).toBe(false)
    expect(shown('loading')).toBe(false)
    expect(text('[data-field="notFoundMessage"]')).toContain(id)
    expect(document.querySelector<HTMLElement>('[data-error]')!.hidden).toBe(true)
  })

  test('a missing id is an error', async () => {
    document.body.innerHTML = body('driver.html')
    history.replaceState({}, '', '/driver.html')
    await initDetailPage(driverPage)
    expect(document.querySelector<HTMLElement>('[data-error]')!.hidden).toBe(false)
    expect(shown('notFound')).toBe(false)
  })

  test('an inactive category is shown with its badge', async () => {
    await open('category.html', 'a1gp', {
      path: 'categories',
      noun: 'category',
      render: (c) => fill(document.body, { name: c.name, inactive: c.is_active ? '' : '1' }),
    })
    expect(shown('loaded')).toBe(true)
    expect(shown('inactive')).toBe(true)
  })
})

describe('article "About" references', () => {
  test('link to their target, or degrade to text when it is gone', async () => {
    const live = { entity_type: 'driver' as const, entity_id: VER }
    const gone = { entity_type: 'driver' as const, entity_id: DRIVER_DELETED }
    const target = await resolveRefs([live, gone])

    expect(aboutFields(live, target(live))).toEqual({ about: 'driver: Max Verstappen', aboutHref: `driver.html?id=${VER}` })
    const missing = aboutFields(gone, target(gone))
    expect(missing.aboutHref).toBe('')

    document.body.innerHTML = '<a href="#" data-field="about" data-field-href="aboutHref"></a>'
    fill(document.body, missing)
    expect(document.querySelector('a')!.hasAttribute('href')).toBe(false)
  })
})
