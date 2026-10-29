import { getMany } from '../api/client.ts'
import { carHref } from '../lib/format.ts'
import { dash, initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'cars',
  async rows(cars) {
    const [categories, teams, engines, chassis] = await Promise.all([
      getMany('categories', cars.map((c) => c.category_id)),
      getMany('teams', cars.map((c) => c.team_id)),
      getMany('engines', cars.map((c) => c.engine_id)),
      getMany('chassis', cars.map((c) => c.chassis_id)),
    ])
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
