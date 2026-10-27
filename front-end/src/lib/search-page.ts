import { ApiError, search, type Page, type Query } from '../api/client.ts'
import type { EntityMap, EntityPath } from '../api/types.ts'
import { fillList, type Fields } from './dom.ts'

export { dash } from './format.ts'

// Shared behaviour for the table pages. Each page's HTML owns its form, header and row <template>;
// this only reads the URL, fetches one page of data and copies it into that HTML.
//
// URL params: col, op, q (filter) · sort, dir (ordering) · page.

/** Related names are resolved from the first page of the related list; the API has no get-by-id. */
export async function lookup<P extends EntityPath>(path: P): Promise<Map<string, EntityMap[P]>> {
  const { items } = await search(path)
  return new Map(items.map((item) => [(item as { id: string }).id, item]))
}

export interface SearchPageConfig<P extends EntityPath> {
  path: P
  rows(items: EntityMap[P][]): Fields[] | Promise<Fields[]>
}

const $ = <T extends HTMLElement>(selector: string) => document.querySelector<T>(selector)!

function withParams(params: URLSearchParams, changes: Record<string, string | null>): string {
  const next = new URLSearchParams(params)
  for (const [key, value] of Object.entries(changes)) {
    if (value == null) next.delete(key)
    else next.set(key, value)
  }
  return `?${next}`
}

/** Only offer operators the back end accepts for the selected column's type. */
function syncOperators(form: HTMLFormElement): void {
  const col = form.elements.namedItem('col') as HTMLSelectElement
  const op = form.elements.namedItem('op') as HTMLSelectElement
  const kind = col.selectedOptions[0]?.dataset.kind ?? 'text'
  for (const option of op.options) option.disabled = !option.dataset.kinds?.split(' ').includes(kind)
  if (op.selectedOptions[0]?.disabled) op.value = [...op.options].find((o) => !o.disabled)!.value
}

export async function initSearchPage<P extends EntityPath>(config: SearchPageConfig<P>): Promise<void> {
  const params = new URLSearchParams(location.search)
  const form = $<HTMLFormElement>('[data-search-form]')

  for (const [name, value] of params) {
    const control = form.elements.namedItem(name)
    if (control instanceof HTMLInputElement || control instanceof HTMLSelectElement) control.value = value
  }
  syncOperators(form)
  const colSelect = form.elements.namedItem('col') as HTMLSelectElement
  colSelect.addEventListener('change', () => syncOperators(form))

  const col = params.get('col')
  const q = params.get('q')
  const sort = params.get('sort')
  const page = Math.max(1, Number(params.get('page')) || 1)
  const query: Query = {
    page,
    filter: col && q ? `${col},${params.get('op') ?? 'ilike'},${q}` : undefined,
    sort: sort ? `${sort},${params.get('dir') === 'desc' ? 'desc' : 'asc'}` : undefined,
  }

  document.querySelectorAll<HTMLAnchorElement>('[data-sort]').forEach((link) => {
    const active = sort === link.dataset.sort
    const nextDir = active && params.get('dir') !== 'desc' ? 'desc' : 'asc'
    link.href = withParams(params, { sort: link.dataset.sort!, dir: nextDir, page: null })
    if (active) link.closest('th')!.ariaSort = params.get('dir') === 'desc' ? 'descending' : 'ascending'
  })

  let result: Page<EntityMap[P]>
  try {
    result = await search(config.path, query)
  } catch (err) {
    const box = $('[data-error]')
    box.textContent = err instanceof ApiError ? `The server rejected this search: ${err.message}` : String(err)
    box.hidden = false
    return
  }

  fillList(document.body, await config.rows(result.items))
  $('[data-empty]').hidden = result.items.length > 0
  $('[data-page]').textContent = String(page)

  const prev = $<HTMLAnchorElement>('[data-prev]')
  prev.hidden = page <= 1
  prev.href = withParams(params, { page: String(page - 1) })
  const next = $<HTMLAnchorElement>('[data-next]')
  next.hidden = !result.hasNext
  next.href = withParams(params, { page: String(page + 1) })
}
