import { fillRelated, initDetailPage, listHref, relatedRaces } from '../lib/detail-page.ts'
import { fill } from '../lib/dom.ts'

// driver.html?id=<driver id>

initDetailPage({
  path: 'drivers',
  noun: 'driver',
  articles: 'driver',
  async render(driver) {
    const name = `${driver.first_name} ${driver.last_name}`
    document.title = `${name} - RaceDB.org`
    fill(document.body, { name, code: driver.code, nationality: driver.nationality })
    fillRelated('wins', await relatedRaces(`winner_id,eq,${driver.id}`), listHref('races.html', 'winner_id', driver.id))
  },
})
