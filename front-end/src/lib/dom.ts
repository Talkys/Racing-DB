// The HTML owns all markup. These helpers only copy fetched values into it:
//   data-field="x"       -> textContent = fields.x
//   data-field-href="x"  -> href = fields.x
//   data-field-value="x" -> value = fields.x (option, input)
//   data-show-if="x"     -> hidden unless fields.x is non-empty
// Elements whose key is absent from `fields` are left untouched, so a page can fill in several passes.
// Repeated rows are cloned from a <template> inside the list element.

export type Fields = Record<string, string | null | undefined>

export function fill(root: ParentNode, fields: Fields): void {
  root.querySelectorAll<HTMLElement>('[data-field]').forEach((el) => {
    const value = fields[el.dataset.field!]
    if (value != null) el.textContent = value
  })
  root.querySelectorAll<HTMLAnchorElement>('[data-field-href]').forEach((el) => {
    const value = fields[el.dataset.fieldHref!]
    if (value != null) el.href = value
  })
  root.querySelectorAll<HTMLOptionElement | HTMLInputElement>('[data-field-value]').forEach((el) => {
    const value = fields[el.dataset.fieldValue!]
    if (value != null) el.value = value
  })
  root.querySelectorAll<HTMLElement>('[data-show-if]').forEach((el) => {
    const key = el.dataset.showIf!
    if (key in fields) el.hidden = !fields[key]
  })
}

export function fillList(list: ParentNode, rows: Fields[]): void {
  const template = list.querySelector('template')!
  const rowsParent = template.parentNode!
  for (const row of rows) {
    const node = template.content.cloneNode(true) as DocumentFragment
    fill(node, row)
    rowsParent.appendChild(node)
  }
}

/** Shows a message in the page's `[data-error]` banner. */
export function showError(err: unknown): void {
  const box = document.querySelector<HTMLElement>('[data-error]')!
  box.textContent = err instanceof Error ? err.message : String(err)
  box.hidden = false
}

export function section(name: string): HTMLElement {
  return document.querySelector<HTMLElement>(`[data-section="${name}"]`)!
}
