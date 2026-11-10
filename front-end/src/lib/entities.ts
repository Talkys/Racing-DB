import { getMany } from '../api/client.ts'
import type { Car, Category, Chassis, Circuit, Driver, Engine, EntityMap, EntityType, Race, Team } from '../api/types.ts'
import { entityHref } from './format.ts'

// What an article is about. articles.entity_id has no foreign key (K12), so the target may be gone: callers get
// `undefined` for it and show the reference as plain text.

export const ENTITY_PATHS = {
  category: 'categories',
  circuit: 'circuits',
  driver: 'drivers',
  race: 'races',
  team: 'teams',
  engine: 'engines',
  chassis: 'chassis',
  car: 'cars',
} as const satisfies Record<EntityType, keyof EntityMap>

type Row = EntityMap[(typeof ENTITY_PATHS)[EntityType]]

/** How each entity type is named in links and titles. */
export function entityName(type: EntityType, row: Row): string {
  switch (type) {
    case 'driver': {
      const d = row as Driver
      return `${d.first_name} ${d.last_name}`
    }
    case 'race': {
      const r = row as Race
      return `${r.name} ${r.season}`
    }
    case 'car': {
      const c = row as Car
      return `${c.name} (${c.season})`
    }
    case 'engine':
    case 'chassis': {
      const m = row as Engine | Chassis
      return `${m.manufacturer} ${m.name}`
    }
    default:
      return (row as Category | Circuit | Team).name
  }
}

export interface EntityRef {
  entity_type: EntityType
  entity_id: string
}

export interface EntityLink {
  name: string
  href: string
}

const refKey = (type: EntityType, id: string) => `${type}:${id.toLowerCase()}`

/** Looks up the targets of the given references, one request per entity type. Missing targets are absent. */
export async function resolveRefs(refs: EntityRef[]): Promise<(ref: EntityRef) => EntityLink | undefined> {
  const types = [...new Set(refs.map((r) => r.entity_type))]
  const found = await Promise.all(
    types.map(async (type) => {
      const rows = await getMany(ENTITY_PATHS[type], refs.filter((r) => r.entity_type === type).map((r) => r.entity_id))
      return [...rows].map(([id, row]) => [refKey(type, id), { name: entityName(type, row), href: entityHref(type, id) }] as const)
    }),
  )
  const links = new Map<string, EntityLink>(found.flat())
  return (ref) => links.get(refKey(ref.entity_type, ref.entity_id))
}

/** "About" fields for an article: the target's name and link, or the bare type with no link when it is gone. */
export function aboutFields(ref: EntityRef, link: EntityLink | undefined): { about: string; aboutHref: string } {
  return link
    ? { about: `${ref.entity_type}: ${link.name}`, aboutHref: link.href }
    : { about: `${ref.entity_type} (no longer in the database)`, aboutHref: '' }
}
