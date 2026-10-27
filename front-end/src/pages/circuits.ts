import { initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'circuit',
  rows: (circuits) => circuits.map((c) => ({ ...c, length_km: String(c.length_km) })),
})
