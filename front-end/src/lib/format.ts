import type { EntityType } from '../api/types.ts'

/** Null, undefined and empty values render as an em dash. */
export const dash = (value: string | number | null | undefined) => (value == null || value === '' ? '-' : String(value))

/** `2026-05-24` or `2026-05-24T12:00:00Z` -> `2026-05-24`. */
export const formatDate = (iso: string | null | undefined) => dash(iso?.slice(0, 10))

export const formatDateTime = (iso: string) => new Date(iso).toLocaleString()

/** `HYBRID_TURBO` -> `hybrid turbo`. */
export const formatEnum = (value: string) => value.replaceAll('_', ' ').toLowerCase()

export const driverName = (d: { first_name: string; last_name: string; code: string }) =>
  `${d.first_name} ${d.last_name} (${d.code})`

export const articleHref = (articleId: string, lang?: string) =>
  `article.html?id=${encodeURIComponent(articleId)}${lang ? `&lang=${encodeURIComponent(lang)}` : ''}`

/** The detail page of each entity type. `chassis.html` is the chassis list, so one chassis has its own page name. */
const DETAIL_PAGES: Record<EntityType, string> = {
  category: 'category.html',
  circuit: 'circuit.html',
  driver: 'driver.html',
  race: 'race.html',
  team: 'team.html',
  engine: 'engine.html',
  chassis: 'chassis-detail.html',
  car: 'car.html',
}

export const entityHref = (entityType: EntityType, entityId: string) => `${DETAIL_PAGES[entityType]}?id=${encodeURIComponent(entityId)}`

export const raceHref = (raceId: string) => entityHref('race', raceId)

export const carHref = (carId: string) => entityHref('car', carId)

/** The href of an optional reference: '' when there is none, which `fill` renders as plain text. */
export const hrefOf = (entityType: EntityType, entityId: string | null | undefined, exists = true) =>
  entityId && exists ? entityHref(entityType, entityId) : ''
