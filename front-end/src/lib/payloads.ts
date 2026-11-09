import { hasValidParts, isF1Race, isMotoGPRace, isWECRace, PAYLOAD_CONTRACTS } from '../api/guards.ts'
import type { Car, Race } from '../api/types.ts'
import type { Fields } from './dom.ts'
import { formatEnum } from './format.ts'

// What race.html and car.html show for the JSONB columns. Rows that pass their guard get the typed boxes; any
// other row falls back to the raw-JSON panel, with a note when the row was meant to follow a contract.

const json = (value: unknown) => JSON.stringify(value, null, 2)

/** "… so it is shown as stored", plus the back end's list of problems when it sent one. */
function brokenContractNote(what: string, errors: string[] | undefined): string {
  const note = `${what} does not match its expected format, so it is shown as stored.`
  return errors?.length ? `${note} Problems: ${errors.join('; ')}.` : note
}

/** Driver ids the race's payload refers to; empty when the payload cannot be trusted. */
export function payloadDriverIds(race: Race): string[] {
  if (isF1Race(race)) return [race.payload.fastest_lap.driver_id, ...race.payload.tyre_strategies.map((s) => s.driver_id)]
  if (isMotoGPRace(race)) return [race.payload.rider_of_the_day_id, race.payload.holeshot_winner_id]
  if (isWECRace(race)) return [race.payload.class_winners.hypercar, race.payload.class_winners.lmgt3]
  return []
}

/**
 * Fields for race.html's "Race Data" box and the rows of its tyre-strategy table. `code` turns a payload driver
 * id into what to show (payload driver ids are not foreign keys).
 */
export function raceDataFields(race: Race, code: (driverId: string) => string): { fields: Fields; tyreStrategies: Fields[] } {
  const fields: Fields = { f1: '', motogp: '', wec: '', rawPayload: '', payloadNote: '' }

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
    const tyreStrategies = p.tyre_strategies.map((s) => ({
      driver: code(s.driver_id),
      stops: s.stops.map((stop) => `lap ${stop.lap}: ${stop.compound}`).join(' → ') || 'no stops',
    }))
    return { fields, tyreStrategies }
  }

  if (isMotoGPRace(race)) {
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
  } else if (race.category_id in PAYLOAD_CONTRACTS || !race.payload_valid) {
    // A contracted category whose payload failed its guard, or a row the back end flagged.
    fields.payloadNote = brokenContractNote('This race data', race.payload_errors)
    fields.rawPayload = json(race.payload)
  } else {
    // Categories without a payload contract: show what the back end stored.
    fields.rawPayload = Object.keys(race.payload ?? {}).length ? json(race.payload) : 'No race data.'
  }
  return { fields, tyreStrategies: [] }
}

/** Fields for car.html's "Parts" box and the rows of its "other" list. */
export function carPartsFields(car: Car): { fields: Fields; other: Fields[] } {
  if (!hasValidParts(car)) {
    return {
      fields: { partsValid: '', partsNote: brokenContractNote('This parts list', car.parts_errors), rawParts: json(car.parts) },
      other: [],
    }
  }
  const { parts } = car
  return {
    fields: {
      partsValid: '1',
      partsNote: '',
      rawParts: '',
      gearboxType: formatEnum(parts.gearbox.type),
      gearboxSpeeds: String(parts.gearbox.speeds),
      brakesFront: parts.brakes.front,
      brakesRear: parts.brakes.rear,
      suspensionFront: parts.suspension.front,
      suspensionRear: parts.suspension.rear,
      tyreSupplier: parts.tyre_supplier,
    },
    other: Object.entries(parts.other).map(([key, value]) => ({ key: formatEnum(key), value })),
  }
}
