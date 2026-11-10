import { entityHref } from '../lib/format.ts'
import { initSearchPage } from '../lib/search-page.ts'

// Inactive categories are listed like the others, with an "inactive" badge.

initSearchPage({
  path: 'categories',
  rows: (categories) =>
    categories.map((c) => ({
      id: c.id,
      name: c.name,
      href: entityHref('category', c.id),
      description: c.description,
      active: c.is_active ? 'Yes' : 'No',
      inactive: c.is_active ? '' : '1',
    })),
})
