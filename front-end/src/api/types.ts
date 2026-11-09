// Mirrors the JSON models in back-end/internal/*/models.go.
// uuid.UUID and time.Time are serialised as strings; pointer fields are nullable.

export interface Category {
  id: string
  name: string
  description: string
  is_active: boolean
}

export interface Circuit {
  id: string
  name: string
  location: string
  country: string
  length_km: number
  created_at: string
}

export interface Driver {
  id: string
  first_name: string
  last_name: string
  code: string
  nationality: string
  created_at: string
}

export interface Race<TPayload = RacePayload> {
  id: string
  category_id: string
  circuit_id: string
  winner_id: string | null
  name: string
  season: number
  round: number
  race_date: string
  /** As stored; check it with the guards in guards.ts before reading nested fields. */
  payload: TPayload
  created_at: string
  /** False when the back end found that payload breaks its category's contract; payload_errors then says how. */
  payload_valid: boolean
  payload_errors?: string[]
}

// Category payload contracts (races.payload). Driver ids inside payloads are not foreign keys.

export type F1Compound = 'SOFT' | 'MEDIUM' | 'HARD' | 'INTERMEDIATE' | 'WET'

export interface F1Payload {
  fastest_lap: { driver_id: string; lap_number: number; lap_time: string }
  tyre_strategies: { driver_id: string; stops: { lap: number; compound: F1Compound }[] }[]
  safety_car_deployments: number
  track_temperature_celsius: number
}

export type MotoGPCompound = 'SOFT' | 'MEDIUM' | 'HARD'

export interface MotoGPPayload {
  rider_of_the_day_id: string
  tyre_spec: { front: MotoGPCompound; rear: MotoGPCompound }
  holeshot_winner_id: string
  weather_conditions: 'DRY' | 'WET'
}

export interface WECPayload {
  total_laps_completed: number
  class_winners: { hypercar: string; lmgt3: string }
  full_course_yellow_count: number
  red_flag_duration_minutes: number
}

/** Categories without a payload contract keep an untyped object. */
export type RacePayload = F1Payload | MotoGPPayload | WECPayload | Record<string, unknown>

export interface Team {
  id: string
  name: string
  short_name: string
  country: string
  founded_year: number | null
  created_at: string
}

export interface Engine {
  id: string
  manufacturer: string
  name: string
  configuration: string
  displacement_cc: number
  cylinders: number
  aspiration: 'NATURALLY_ASPIRATED' | 'TURBOCHARGED' | 'HYBRID_TURBO'
  power_kw: number | null
  created_at: string
}

export interface Chassis {
  id: string
  manufacturer: string
  name: string
  material: string
  wheelbase_mm: number | null
  weight_kg: number | null
  created_at: string
}

export interface Car {
  id: string
  category_id: string
  team_id: string
  engine_id: string
  chassis_id: string
  name: string
  season: number
  /** As stored; check it with hasValidParts (guards.ts) before reading nested fields. */
  parts: CarParts | Record<string, unknown>
  created_at: string
  /** False when the back end found that parts breaks the contract; parts_errors then says how. */
  parts_valid: boolean
  parts_errors?: string[]
}

/** cars.parts: same shape for every category. */
export interface CarParts {
  gearbox: { type: string; speeds: number }
  brakes: { front: string; rear: string }
  suspension: { front: string; rear: string }
  tyre_supplier: string
  other: Record<string, string>
}

export interface Article {
  id: string
  entity_type: 'category' | 'circuit' | 'driver' | 'race' | 'team' | 'engine' | 'chassis' | 'car'
  entity_id: string
  kind: string
  published_at: string | null
  created_at: string
}

export interface ArticleText {
  id: string
  article_id: string
  language: string
  title: string
  summary: string | null
  body_md: string
  created_at: string
  updated_at: string
}

/** Fields the back end computes on read (JSONB validation); they are not database columns. */
export type ComputedFields = 'payload_valid' | 'payload_errors' | 'parts_valid' | 'parts_errors'

/** Entity name -> model, keyed by the route segment under /api/1. */
export interface EntityMap {
  categories: Category
  circuits: Circuit
  drivers: Driver
  races: Race
  teams: Team
  engines: Engine
  chassis: Chassis
  cars: Car
  articles: Article
  'article-texts': ArticleText
}

export type EntityPath = keyof EntityMap

/** A piece of a search preview; `match` marks words that matched the query. */
export interface PreviewSegment {
  text: string
  match: boolean
}

/**
 * GET /api/1/articles/search: one article with the text chosen for it (the requested language when it
 * qualifies, else the best match). language/title/summary are null only when browsing an article with no texts.
 */
export interface ArticleSearchResult {
  id: string
  entity_type: Article['entity_type']
  entity_id: string
  kind: string
  published_at: string | null
  language: string | null
  title: string | null
  summary: string | null
  preview: PreviewSegment[]
}
