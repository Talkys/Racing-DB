import type { Car, CarParts, F1Payload, MotoGPPayload, Race, WECPayload } from './types.ts'

// The payload contract is chosen by the race's category, but the JSONB is whatever the ETL stored. A guard passes
// only when the back end did not flag the row (payload_valid/parts_valid) and the shape holds everything the
// pages read, so a page can trust the narrowed type. These check types and enums, not the back end's ranges.

type Obj = Record<string, unknown>

const isObj = (v: unknown): v is Obj => typeof v === 'object' && v !== null && !Array.isArray(v)
const isStr = (v: unknown): v is string => typeof v === 'string' && v.trim() !== ''
const isNum = (v: unknown): v is number => typeof v === 'number' && Number.isFinite(v)
const isInt = (v: unknown): v is number => Number.isInteger(v)
const oneOf = (values: readonly string[]) => (v: unknown) => typeof v === 'string' && values.includes(v)
// Driver ids are looked up with `id,in,...`, which rejects anything that is not a UUID.
const isUuid = (v: unknown) => typeof v === 'string' && /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(v)
const isPair = (v: unknown, check: (side: unknown) => boolean) => isObj(v) && check(v.front) && check(v.rear)

const isF1Compound = oneOf(['SOFT', 'MEDIUM', 'HARD', 'INTERMEDIATE', 'WET'])
const isMotoGPCompound = oneOf(['SOFT', 'MEDIUM', 'HARD'])

export function isF1Payload(p: unknown): p is F1Payload {
  if (!isObj(p) || !isObj(p.fastest_lap) || !Array.isArray(p.tyre_strategies)) return false
  const lap = p.fastest_lap
  return (
    isUuid(lap.driver_id) &&
    isInt(lap.lap_number) &&
    isStr(lap.lap_time) &&
    p.tyre_strategies.every(
      (s) => isObj(s) && isUuid(s.driver_id) && Array.isArray(s.stops) && s.stops.every((stop) => isObj(stop) && isInt(stop.lap) && isF1Compound(stop.compound)),
    ) &&
    isInt(p.safety_car_deployments) &&
    isNum(p.track_temperature_celsius)
  )
}

export function isMotoGPPayload(p: unknown): p is MotoGPPayload {
  return (
    isObj(p) &&
    isUuid(p.rider_of_the_day_id) &&
    isPair(p.tyre_spec, isMotoGPCompound) &&
    isUuid(p.holeshot_winner_id) &&
    oneOf(['DRY', 'WET'])(p.weather_conditions)
  )
}

export function isWECPayload(p: unknown): p is WECPayload {
  return (
    isObj(p) &&
    isInt(p.total_laps_completed) &&
    isObj(p.class_winners) &&
    isUuid(p.class_winners.hypercar) &&
    isUuid(p.class_winners.lmgt3) &&
    isInt(p.full_course_yellow_count) &&
    isInt(p.red_flag_duration_minutes)
  )
}

export function isCarParts(p: unknown): p is CarParts {
  return (
    isObj(p) &&
    isObj(p.gearbox) &&
    isStr(p.gearbox.type) &&
    isInt(p.gearbox.speeds) &&
    isPair(p.brakes, isStr) &&
    isPair(p.suspension, isStr) &&
    isStr(p.tyre_supplier) &&
    isObj(p.other) &&
    Object.values(p.other).every((v) => typeof v === 'string')
  )
}

/** The categories whose payload has a contract (spec 3.3), with how to check it. */
export const PAYLOAD_CONTRACTS: Record<string, (p: unknown) => boolean> = {
  f1: isF1Payload,
  motogp: isMotoGPPayload,
  wec: isWECPayload,
}

export const isF1Race = (race: Race): race is Race<F1Payload> => race.category_id === 'f1' && race.payload_valid && isF1Payload(race.payload)
export const isMotoGPRace = (race: Race): race is Race<MotoGPPayload> =>
  race.category_id === 'motogp' && race.payload_valid && isMotoGPPayload(race.payload)
export const isWECRace = (race: Race): race is Race<WECPayload> => race.category_id === 'wec' && race.payload_valid && isWECPayload(race.payload)

export const hasValidParts = (car: Car): car is Car & { parts: CarParts } => car.parts_valid && isCarParts(car.parts)
