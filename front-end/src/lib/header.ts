// The header search box jumps to the matching table page with a "contains" filter.
const form = document.querySelector<HTMLFormElement>('[data-header-search]')

form?.addEventListener('submit', (event) => {
  event.preventDefault()
  const type = form.elements.namedItem('type') as HTMLSelectElement
  const q = (form.elements.namedItem('q') as HTMLInputElement).value.trim()
  const params = new URLSearchParams()
  if (q) { params.set('col', type.selectedOptions[0].dataset.col!); params.set('op', 'ilike'); params.set('q', q) }
  location.href = `${type.value}${params.size ? `?${params}` : ''}`
})
