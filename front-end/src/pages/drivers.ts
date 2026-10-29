import { initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'drivers',
  rows: (drivers) => drivers.map((d) => ({ ...d })),
})
