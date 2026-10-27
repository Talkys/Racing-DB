import { dash, initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'engine',
  rows: (engines) => engines.map((e) => ({
    ...e,
    displacement_cc: String(e.displacement_cc),
    cylinders: String(e.cylinders),
    aspiration: e.aspiration.replace('_', ' ').toLowerCase(),
    power_kw: dash(e.power_kw),
  })),
})
