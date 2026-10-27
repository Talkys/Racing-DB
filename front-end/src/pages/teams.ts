import { dash, initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'team',
  rows: (teams) => teams.map((t) => ({ ...t, founded_year: dash(t.founded_year) })),
})
