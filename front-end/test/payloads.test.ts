import { afterAll, beforeAll, beforeEach, describe, expect, test } from 'bun:test'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { GlobalRegistrator } from '@happy-dom/global-registrator'
import { hasValidParts, isCarParts, isF1Payload, isF1Race, isMotoGPPayload, isWECPayload } from '../src/api/guards.ts'
import { mockGet, mockSearch } from '../src/api/mock.ts'
import type { Car, CarParts, F1Payload, MotoGPPayload, Race, WECPayload } from '../src/api/types.ts'
import { fill, fillList, section } from '../src/lib/dom.ts'
import { carPartsFields, payloadDriverIds, raceDataFields } from '../src/lib/payloads.ts'

// Well-known fixture rows (same ids as back-end/cmd/api/api_test.go).
const MONACO = 'c39a7593-36c5-4ceb-8d07-6bc82f808381'
const MUGELLO = 'd40b8604-47d6-5dfc-9e18-7cd93f919492'
const LEMANS = 'e51c9715-58e7-6e0d-af29-8de04fa205a3'
const RACE_BROKEN = '77000000-0000-4000-8000-000000000090' // f1 race with payload {}
const CAR_RB22 = '74000000-0000-4000-8000-000000000001'
const CAR_BROKEN = '74000000-0000-4000-8000-000000000090' // parts missing most keys
const VER = 'b11a43a8-4e36-47a8-8e68-301fc1cae6ef'

const race = (id: string) => mockGet('races', id)!
const car = (id: string) => mockGet('cars', id)!
const code = (driverId: string) => `code(${driverId})`
const clone = <T>(x: T): T => structuredClone(x)

describe('payload guards', () => {
  test('every fixture except the broken ones passes its guard', () => {
    const races = Array.from({ length: 3 }, (_, i) => mockSearch('races', { page: i + 1 }).slice(0, 25)).flat()
    expect(races.length).toBeGreaterThan(25)
    for (const r of races) expect({ id: r.id, valid: r.payload_valid }).toEqual({ id: r.id, valid: r.id !== RACE_BROKEN })

    const cars = Array.from({ length: 3 }, (_, i) => mockSearch('cars', { page: i + 1 }).slice(0, 25)).flat()
    for (const c of cars) expect({ id: c.id, valid: c.parts_valid }).toEqual({ id: c.id, valid: c.id !== CAR_BROKEN })
  })

  test('f1: a missing key, a wrong type, a bad enum or a non-UUID driver id fails', () => {
    const ok = race(MONACO).payload as F1Payload
    expect(isF1Payload(ok)).toBe(true)
    // Object.assign lets a fixture take values its type forbids, as a bad ETL row would.
    const broken: ((p: F1Payload) => void)[] = [
      (p) => Object.assign(p, { fastest_lap: undefined }),
      (p) => Object.assign(p.fastest_lap, { lap_time: undefined }),
      (p) => Object.assign(p, { safety_car_deployments: '1' }),
      (p) => Object.assign(p, { track_temperature_celsius: null }),
      (p) => Object.assign(p, { tyre_strategies: {} }),
      (p) => Object.assign(p.tyre_strategies[0].stops[0], { compound: 'SUPERSOFT' }),
      (p) => Object.assign(p.tyre_strategies[0], { stops: [null] }),
      (p) => Object.assign(p.fastest_lap, { driver_id: 'VER' }),
    ]
    for (const breakIt of broken) {
      const p = clone(ok)
      breakIt(p)
      expect(isF1Payload(p)).toBe(false)
    }
  })

  test('motogp and wec', () => {
    const motogp = race(MUGELLO).payload as MotoGPPayload
    expect(isMotoGPPayload(motogp)).toBe(true)
    expect(isMotoGPPayload({ ...motogp, tyre_spec: { front: 'HARD' } })).toBe(false)
    expect(isMotoGPPayload({ ...motogp, weather_conditions: 'MIXED' })).toBe(false)

    const wec = race(LEMANS).payload as WECPayload
    expect(isWECPayload(wec)).toBe(true)
    expect(isWECPayload({ ...wec, total_laps_completed: 368.5 })).toBe(false)
    expect(isWECPayload({ ...wec, class_winners: null })).toBe(false)
  })

  test('car parts', () => {
    const ok = car(CAR_RB22).parts as CarParts
    expect(isCarParts(ok)).toBe(true)
    expect(isCarParts({ ...ok, other: {} })).toBe(true)
    expect(isCarParts({ ...ok, gearbox: { type: 'SEQUENTIAL' } })).toBe(false)
    expect(isCarParts({ ...ok, brakes: { front: 'Carbon' } })).toBe(false)
    expect(isCarParts({ ...ok, other: { kw: 120 } })).toBe(false)
    expect(isCarParts([])).toBe(false)
  })

  test('a row the back end flagged fails its guard even when the shape holds', () => {
    const r: Race = { ...race(MONACO), payload_valid: false }
    expect(isF1Race(r)).toBe(false)
    const c: Car = { ...car(CAR_RB22), parts_valid: false }
    expect(hasValidParts(c)).toBe(false)
  })

  test('a category that is not f1 never passes the f1 guard', () => {
    expect(isF1Race({ ...race(MONACO), category_id: 'f2' })).toBe(false)
  })
})

describe('raceDataFields', () => {
  test('a valid payload fills its category box', () => {
    const { fields, tyreStrategies } = raceDataFields(race(MONACO), code)
    expect(fields).toMatchObject({ f1: '1', motogp: '', wec: '', rawPayload: '', payloadNote: '', fastestLapDriver: code(VER) })
    expect(tyreStrategies).toEqual([{ driver: code(VER), stops: 'lap 30: HARD' }])
    expect(payloadDriverIds(race(MONACO))).toEqual([VER, VER])
  })

  test('a broken payload falls back to the raw JSON with a note', () => {
    const { fields, tyreStrategies } = raceDataFields(race(RACE_BROKEN), code)
    expect(fields).toMatchObject({ f1: '', motogp: '', wec: '', rawPayload: '{}' })
    expect(fields.payloadNote).toContain('does not match its expected format')
    expect(tyreStrategies).toEqual([])
    expect(payloadDriverIds(race(RACE_BROKEN))).toEqual([])
  })

  test("the note lists the back end's problems", () => {
    const r: Race = { ...race(RACE_BROKEN), payload_errors: ['fastest_lap: missing', 'tyre_strategies: missing'] }
    expect(raceDataFields(r, code).fields.payloadNote).toContain('Problems: fastest_lap: missing; tyre_strategies: missing.')
  })

  test('a category without a contract shows its raw JSON without a note', () => {
    const r: Race = { ...race(MONACO), category_id: 'indycar', payload: { pole: 'x' } }
    expect(raceDataFields(r, code).fields).toMatchObject({ rawPayload: '{\n  "pole": "x"\n}', payloadNote: '' })
    expect(raceDataFields({ ...r, payload: {} }, code).fields.rawPayload).toBe('No race data.')
  })
})

describe('carPartsFields', () => {
  test('valid parts fill the parts box', () => {
    const { fields, other } = carPartsFields(car(CAR_RB22))
    expect(fields).toMatchObject({ partsValid: '1', rawParts: '', partsNote: '', gearboxSpeeds: '8', tyreSupplier: 'Pirelli' })
    expect(other).toEqual([{ key: 'fuel supplier', value: 'Petronas' }])
  })

  test('broken parts fall back to the raw JSON with a note', () => {
    const { fields, other } = carPartsFields(car(CAR_BROKEN))
    expect(fields.partsValid).toBe('')
    expect(JSON.parse(fields.rawParts!)).toEqual(car(CAR_BROKEN).parts)
    expect(fields.partsNote).toContain('does not match its expected format')
    expect(other).toEqual([])
  })
})

// The real templates, filled the way race.ts/car.ts fill them.
describe('pages render the fallback', () => {
  // The page's <body>, without its <script> tags (happy-dom would try to load them).
  const body = (file: string) =>
    readFileSync(resolve(import.meta.dirname, '..', file), 'utf8')
      .match(/<body[^>]*>([\s\S]*)<\/body>/)![1]
      .replace(/<script[\s\S]*?<\/script>/g, '')
  const visible = (selector: string) => {
    const el = document.querySelector<HTMLElement>(selector)!
    for (let node: HTMLElement | null = el; node; node = node.parentElement) if (node.hidden) return false
    return true
  }

  beforeAll(() => GlobalRegistrator.register())
  afterAll(() => GlobalRegistrator.unregister())
  beforeEach(() => {
    document.body.innerHTML = ''
  })

  test('race.html with the broken f1 race', () => {
    document.body.innerHTML = body('race.html')
    fill(document.body, { loaded: '1', ...raceDataFields(race(RACE_BROKEN), code).fields })

    expect(visible('[data-field="rawPayload"]')).toBe(true)
    expect(document.querySelector('[data-field="rawPayload"]')!.textContent).toBe('{}')
    expect(visible('[data-field="payloadNote"]')).toBe(true)
    expect(visible('[data-show-if="f1"]')).toBe(false)
  })

  test('race.html with a valid f1 race', () => {
    document.body.innerHTML = body('race.html')
    fill(document.body, { loaded: '1', ...raceDataFields(race(MONACO), code).fields })
    expect(visible('[data-show-if="f1"]')).toBe(true)
    expect(visible('[data-field="rawPayload"]')).toBe(false)
    expect(visible('[data-field="payloadNote"]')).toBe(false)
  })

  test('car.html with the broken car', () => {
    document.body.innerHTML = body('car.html')
    const { fields, other } = carPartsFields(car(CAR_BROKEN))
    fill(document.body, { loaded: '1', ...fields })
    fillList(section('parts'), other)

    expect(visible('[data-field="rawParts"]')).toBe(true)
    expect(visible('[data-field="partsNote"]')).toBe(true)
    expect(visible('[data-section="parts"]')).toBe(false)
  })

  test('car.html with a valid car', () => {
    document.body.innerHTML = body('car.html')
    fill(document.body, { loaded: '1', ...carPartsFields(car(CAR_RB22)).fields })
    expect(visible('[data-section="parts"]')).toBe(true)
    expect(visible('[data-field="rawParts"]')).toBe(false)
  })
})
