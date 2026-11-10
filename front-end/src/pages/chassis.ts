import { entityHref } from '../lib/format.ts'
import { dash, initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'chassis',
  rows: (chassis) =>
    chassis.map((c) => ({ ...c, wheelbase_mm: dash(c.wheelbase_mm), weight_kg: dash(c.weight_kg), href: entityHref('chassis', c.id) })),
})
