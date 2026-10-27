import { initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'driver',
  rows: (drivers) => drivers.map((d) => ({ ...d })),
})
