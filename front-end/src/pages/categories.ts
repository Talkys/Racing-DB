import { initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'categories',
  rows: (categories) => categories.map((c) => ({ id: c.id, name: c.name, description: c.description, active: c.is_active ? 'Yes' : 'No' })),
})
