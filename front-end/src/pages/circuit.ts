import { fillRelated, initDetailPage, listHref, relatedRaces } from '../lib/detail-page.ts'
import { fill } from '../lib/dom.ts'

// circuit.html?id=<circuit id>

initDetailPage({
  path: 'circuits',
  noun: 'circuit',
  articles: 'circuit',
  async render(circuit) {
    document.title = `${circuit.name} - RaceDB.org`
    fill(document.body, { name: circuit.name, location: circuit.location, country: circuit.country, length: String(circuit.length_km) })
    fillRelated('races', await relatedRaces(`circuit_id,eq,${circuit.id}`), listHref('races.html', 'circuit_id', circuit.id))
  },
})
