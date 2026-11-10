import { entityHref } from '../lib/format.ts'
import { dash, initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'teams',
  rows: (teams) => teams.map((t) => ({ ...t, founded_year: dash(t.founded_year), href: entityHref('team', t.id) })),
})
