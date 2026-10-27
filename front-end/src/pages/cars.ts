import { carHref } from '../lib/format.ts'
import { dash, initSearchPage, lookup } from '../lib/search-page.ts'

initSearchPage({
  path: 'car',
  async rows(cars) {
    const [categories, teams, engines, chassis] = await Promise.all([lookup('categories'), lookup('team'), lookup('engine'), lookup('chassis')])
    return cars.map((c) => ({
      name: c.name,
      href: carHref(c.id),
      season: String(c.season),
      category: categories.get(c.category_id)?.name ?? c.category_id,
      team: dash(teams.get(c.team_id)?.short_name),
      engine: dash(engines.get(c.engine_id)?.name),
      chassis: dash(chassis.get(c.chassis_id)?.name),
    }))
  },
})
