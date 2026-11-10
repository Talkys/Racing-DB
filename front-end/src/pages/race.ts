import { get } from '../api/client.ts'
import { initDetailPage } from '../lib/detail-page.ts'
import { fill, fillList, section } from '../lib/dom.ts'
import { loadDrivers } from '../lib/drivers.ts'
import { dash, driverName, formatDate, hrefOf } from '../lib/format.ts'
import { payloadDriverIds, raceDataFields } from '../lib/payloads.ts'

// race.html?id=<race id>

initDetailPage({
  path: 'races',
  noun: 'race',
  articles: 'race',
  async render(race) {
    const [category, circuit, drivers] = await Promise.all([
      get('categories', race.category_id),
      get('circuits', race.circuit_id),
      loadDrivers([race.winner_id, ...payloadDriverIds(race)]),
    ])
    // Payload driver ids are not foreign keys: show the code when known, else the raw id.
    const code = (driverId: string) => drivers.get(driverId)?.code ?? driverId
    const winner = race.winner_id ? drivers.get(race.winner_id) : undefined

    document.title = `${race.name} ${race.season} - RaceDB.org`
    const raceData = raceDataFields(race, code)
    fill(document.body, {
      name: `${race.name} ${race.season}`,
      date: formatDate(race.race_date),
      season: String(race.season),
      round: String(race.round),
      category: category?.name ?? race.category_id,
      categoryHref: hrefOf('category', race.category_id, !!category),
      winner: winner ? driverName(winner) : race.winner_id ?? '-',
      winnerHref: hrefOf('driver', race.winner_id, !!winner),
      circuitName: dash(circuit?.name),
      circuitHref: hrefOf('circuit', race.circuit_id, !!circuit),
      circuitLocation: dash(circuit?.location),
      circuitCountry: dash(circuit?.country),
      circuitLength: dash(circuit?.length_km),
      ...raceData.fields,
    })
    fillList(section('tyre-strategies'), raceData.tyreStrategies)
  },
})
