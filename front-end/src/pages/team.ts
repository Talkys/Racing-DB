import { fillRelated, initDetailPage, listHref, relatedCars } from '../lib/detail-page.ts'
import { fill } from '../lib/dom.ts'
import { dash } from '../lib/format.ts'

// team.html?id=<team id>

initDetailPage({
  path: 'teams',
  noun: 'team',
  articles: 'team',
  async render(team) {
    document.title = `${team.name} - RaceDB.org`
    fill(document.body, { name: team.name, shortName: team.short_name, country: team.country, founded: dash(team.founded_year) })
    fillRelated('cars', await relatedCars(`team_id,eq,${team.id}`), listHref('cars.html', 'team_id', team.id))
  },
})
