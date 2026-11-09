import { get } from '../api/client.ts'
import { entityArticles } from '../lib/articles.ts'
import { type Fields, fill, fillList, section, showError } from '../lib/dom.ts'
import { loadDrivers } from '../lib/drivers.ts'
import { dash, driverName, formatDate } from '../lib/format.ts'
import { payloadDriverIds, raceDataFields } from '../lib/payloads.ts'

// race.html?id=<race id>

async function load(id: string): Promise<void> {
  const race = await get('races', id)
  if (!race) throw new Error(`Race with ID '${id}' does not exist.`)

  const [category, circuit, drivers, articles] = await Promise.all([
    get('categories', race.category_id),
    get('circuits', race.circuit_id),
    loadDrivers([race.winner_id, ...payloadDriverIds(race)]),
    entityArticles('race', race.id),
  ])
  // Payload driver ids are not foreign keys: show the code when known, else the raw id.
  const code = (driverId: string) => drivers.get(driverId)?.code ?? driverId
  const winner = race.winner_id ? drivers.get(race.winner_id) : undefined

  document.title = `${race.name} ${race.season} - RaceDB.org`
  const fields: Fields = {
    loaded: '1',
    name: `${race.name} ${race.season}`,
    date: formatDate(race.race_date),
    season: String(race.season),
    round: String(race.round),
    category: category?.name ?? race.category_id,
    categoryHref: `races.html?col=category_id&op=eq&q=${encodeURIComponent(race.category_id)}`,
    winner: winner ? driverName(winner) : race.winner_id ?? '-',
    circuitName: dash(circuit?.name),
    circuitLocation: dash(circuit?.location),
    circuitCountry: dash(circuit?.country),
    circuitLength: dash(circuit?.length_km),
  }

  const raceData = raceDataFields(race, code)
  fill(document.body, { ...fields, ...raceData.fields })
  fillList(section('tyre-strategies'), raceData.tyreStrategies)
  fillList(section('articles'), articles.map((a) => a.fields))
  section('articles').querySelector<HTMLElement>('[data-empty]')!.hidden = articles.length > 0
}

const id = new URLSearchParams(location.search).get('id')
try {
  if (!id) throw new Error('No race id given. Use race.html?id=<race id>.')
  await load(id)
} catch (err) {
  showError(err)
} finally {
  fill(document.body, { loading: '' })
}
