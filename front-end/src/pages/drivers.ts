import { entityHref } from '../lib/format.ts'
import { initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'drivers',
  rows: (drivers) => drivers.map((d) => ({ ...d, href: entityHref('driver', d.id) })),
})
