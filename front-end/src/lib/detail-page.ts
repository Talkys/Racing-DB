import { ApiError, get, getMany, search } from '../api/client.ts'
import type { Car, EntityMap, EntityPath, EntityType, Race } from '../api/types.ts'
import { entityArticles } from './articles.ts'
import { type Fields, fill, fillList, section, showError } from './dom.ts'
import { carHref, formatDate, hrefOf, raceHref } from './format.ts'

// Shared behaviour for the detail pages (<page>.html?id=<id>). Each page's HTML owns its markup:
//   [data-show-if="loading"]   shown until the load settles
//   [data-show-if="loaded"]    the entity itself, filled by `render`
//   [data-show-if="notFound"]  shown with [data-field="notFoundMessage"] when the id matches no row
//   [data-section="articles"]  optional: the published articles about the entity
//   [data-error]               any other failure
// Related lists (races at a circuit, cars of a team, ...) are sections holding a row <template>, a [data-empty]
// line and a [data-show-if="more"] link to the full list ([data-field-href="allHref"]).

export interface DetailPageConfig<P extends EntityPath> {
  path: P
  /** What the page shows, as in `No driver with ID ...`. */
  noun: string
  /** Set when articles can be about this entity: lists them in [data-section="articles"]. */
  articles?: EntityType
  /** Fills the page from the row. The page is shown once it resolves. */
  render(row: EntityMap[P]): void | Promise<void>
}

function fillArticles(cards: { fields: Fields }[]): void {
  const box = section('articles')
  fillList(box, cards.map((a) => a.fields))
  box.querySelector<HTMLElement>('[data-empty]')!.hidden = cards.length > 0
}

export async function initDetailPage<P extends EntityPath>(config: DetailPageConfig<P>): Promise<void> {
  const file = location.pathname.split('/').pop() || `${config.noun}.html`
  const id = new URLSearchParams(location.search).get('id')?.trim()
  try {
    if (!id) throw new Error(`No ${config.noun} id given. Use ${file}?id=<${config.noun} id>.`)
    // A malformed id is rejected with 400, but for the reader it is the same as an unknown one.
    const row = await get(config.path, id).catch((err) => (err instanceof ApiError && err.status === 400 ? null : Promise.reject(err)))
    if (!row) {
      fill(document.body, { notFound: '1', notFoundMessage: `No ${config.noun} with ID '${id}' exists. It may have been removed.` })
      return
    }
    const entityId = (row as { id: string }).id
    const [articles] = await Promise.all([config.articles ? entityArticles(config.articles, entityId) : undefined, config.render(row)])
    if (articles) fillArticles(articles)
    fill(document.body, { loaded: '1' })
  } catch (err) {
    showError(err)
  } finally {
    fill(document.body, { loading: '' })
  }
}

/** First page of a related list, plus whether there is more than that page. */
export interface Related {
  rows: Fields[]
  more: boolean
}

/** Fills a related-list section; `allHref` is the list page holding every row. */
export function fillRelated(name: string, related: Related, allHref: string): void {
  const box = section(name)
  fillList(box, related.rows)
  box.querySelector<HTMLElement>('[data-empty]')!.hidden = related.rows.length > 0
  fill(box, { more: related.more ? '1' : '', allHref })
}

/** Races matching `filter`, newest first: date, name, href, category, categoryHref. */
export async function relatedRaces(filter: string): Promise<Related> {
  const { items, hasNext } = await search('races', { filter })
  const categories = await getMany('categories', items.map((r) => r.category_id))
  return {
    more: hasNext,
    rows: items.map((r: Race) => ({
      date: formatDate(r.race_date),
      name: `${r.name} ${r.season}`,
      href: raceHref(r.id),
      category: categories.get(r.category_id)?.name ?? r.category_id,
      categoryHref: hrefOf('category', r.category_id),
    })),
  }
}

/** Cars matching `filter`, newest season first: season, name, href, category, categoryHref. */
export async function relatedCars(filter: string): Promise<Related> {
  const { items, hasNext } = await search('cars', { filter })
  const categories = await getMany('categories', items.map((c) => c.category_id))
  return {
    more: hasNext,
    rows: items.map((c: Car) => ({
      season: String(c.season),
      name: c.name,
      href: carHref(c.id),
      category: categories.get(c.category_id)?.name ?? c.category_id,
      categoryHref: hrefOf('category', c.category_id),
    })),
  }
}

/** A list page filtered on one column, e.g. every race at a circuit. */
export const listHref = (page: string, column: string, value: string) => `${page}?col=${column}&op=eq&q=${encodeURIComponent(value)}`
