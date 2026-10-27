import { get } from '../api/client.ts'
import { isF1Race, isMotoGPRace, isWECRace } from '../api/guards.ts'
import { entityArticles } from '../lib/articles.ts'
import { fill, fillList, section, showError, type Fields } from '../lib/dom.ts'
import { loadDrivers } from '../lib/drivers.ts'
import { dash, driverName, formatDate } from '../lib/format.ts'

// race.html?id=<race id>

async function load(id: string): Promise<void> {
  const race = await get('race', id)
  if (!race) throw new Error(`Race with ID '${id}' does not exist.`)

  const payloadDriverIds = isF1Race(race)
    ? [race.payload.fastest_lap.driver_id, ...race.payload.tyre_strategies.map((s) => s.driver_id)]
    : isMotoGPRace(race)
      ? [race.payload.rider_of_the_day_id, race.payload.holeshot_winner_id]
      : isWECRace(race)
        ? [race.payload.class_winners.hypercar, race.payload.class_winners.lmgt3]
        : []

  const [category, circuit, drivers, articles] = await Promise.all([
    get('categories', race.category_id),
    get('circuit', race.circuit_id),
    loadDrivers([race.winner_id, ...payloadDriverIds]),
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
    winner: winner ? driverName(winner) : race.winner_id ?? '—',
    circuitName: dash(circuit?.name),
    circuitLocation: dash(circuit?.location),
    circuitCountry: dash(circuit?.country),
    circuitLength: dash(circuit?.length_km),
    f1: '',
    motogp: '',
    wec: '',
    rawPayload: '',
  }

  if (isF1Race(race)) {
    const p = race.payload
    Object.assign(fields, {
      f1: '1',
      fastestLapDriver: code(p.fastest_lap.driver_id),
      fastestLapNumber: String(p.fastest_lap.lap_number),
      fastestLapTime: p.fastest_lap.lap_time,
      safetyCars: String(p.safety_car_deployments),
      trackTemperature: String(p.track_temperature_celsius),
    })
    fillList(
      section('tyre-strategies'),
      p.tyre_strategies.map((s) => ({
        driver: code(s.driver_id),
        stops: s.stops.map((stop) => `lap ${stop.lap}: ${stop.compound}`).join(' → ') || 'no stops',
      })),
    )
  } else if (isMotoGPRace(race)) {
    const p = race.payload
    Object.assign(fields, {
      motogp: '1',
      riderOfTheDay: code(p.rider_of_the_day_id),
      holeshotWinner: code(p.holeshot_winner_id),
      tyreFront: p.tyre_spec.front,
      tyreRear: p.tyre_spec.rear,
      weather: p.weather_conditions,
    })
  } else if (isWECRace(race)) {
    const p = race.payload
    Object.assign(fields, {
      wec: '1',
      totalLaps: String(p.total_laps_completed),
      hypercarWinner: code(p.class_winners.hypercar),
      lmgt3Winner: code(p.class_winners.lmgt3),
      fcyCount: String(p.full_course_yellow_count),
      redFlagMinutes: String(p.red_flag_duration_minutes),
    })
  } else {
    // Categories without a payload contract: show what the back end stored.
    fields.rawPayload = Object.keys(race.payload).length ? JSON.stringify(race.payload, null, 2) : 'No race data.'
  }

  fill(document.body, fields)
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
